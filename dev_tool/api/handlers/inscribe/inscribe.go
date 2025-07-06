package inscribe

import (
	"bytes"
	"dev_tool/api/response"
	"dev_tool/chain/metaid-inscribe-sdk/common"
	"dev_tool/chain/metaid-inscribe-sdk/inscribe"
	"dev_tool/models"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/gin-gonic/gin"
)

// Utxo 结构
type Utxo struct {
	OutTxId  string `json:"outTxId"`
	OutIndex uint32 `json:"outIndex"`
	OutValue int64  `json:"outValue"`
}

// CreateInscribe 创建铭文
func CreateInscribe(c *gin.Context) {
	var req CreateInscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "无效的请求参数: "+err.Error())
		return
	}

	// 查找链信息
	var chain models.Chain
	if err := models.DB.First(&chain, req.ChainID).Error; err != nil {
		response.Error(c, "指定的链不存在")
		return
	}

	// 创建铭文请求
	metaIdOpRequest := &inscribe.MetaIdInscribeRequest{
		Net:           common.GetNetParams(string(chain.ChainType)),
		MetaIdFlag:    req.MetaIDFlag,
		Path:          req.Path,
		Operation:     req.Operation,
		Payload:       req.Payload,
		PinOutValue:   req.PinOutValue,
		PinOutAddress: req.Address,
		ChangeAddress: req.Address,
		OtherOuts:     []*inscribe.OtherOut{},
	}

	// 构建铭文交易
	metaIdInscribeBuilder, minerFee, err := inscribe.MetaIdInscribeBuilder(metaIdOpRequest, req.FeeRate)
	if err != nil {
		response.Error(c, "构建铭文交易失败: "+err.Error())
		return
	}

	// 获取 Reveal 交易相关信息
	revealPrivateKey := metaIdInscribeBuilder.RevealPrivateKeyHex
	revealAddress := metaIdInscribeBuilder.RevealAddress
	revealPkScript := hex.EncodeToString(metaIdInscribeBuilder.TxCtxData.CommitTxAddressPkScript)
	revealInputIndex := int(metaIdInscribeBuilder.RevealTaprootDataInputIndex)
	redeemScript := hex.EncodeToString(metaIdInscribeBuilder.TxCtxData.InscriptionScript)
	controlBlockWitness := hex.EncodeToString(metaIdInscribeBuilder.TxCtxData.ControlBlockWitness)

	// 创建 Commit 交易
	commitTx, commitTxRaw, commitTxOutIndex, err := makeCommitTx(&CommitTxParams{
		ChainType:      string(chain.ChainType),
		RevealAddress:  revealAddress,
		FromAddress:    req.Address,
		RevealFee:      minerFee,
		NetworkFeeRate: req.FeeRate,
		UtxoList:       req.UtxoList, // 使用请求中的 UTXO 列表
	})
	if err != nil {
		response.Error(c, "创建Commit交易失败: "+err.Error())
		return
	}

	// 完成 Reveal 交易
	txHash := commitTx.TxHash()
	commitPreOutPoint := wire.NewOutPoint(&txHash, uint32(commitTxOutIndex))
	metaIdInscribeBuilder.RevealPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxIn[revealInputIndex].PreviousOutPoint = *commitPreOutPoint

	// 签名 Reveal 交易
	taprootInSigner := &common.InputSign{
		UtxoType:            common.Taproot,
		Index:               revealInputIndex,
		PkScript:            revealPkScript,
		RedeemScript:        redeemScript,
		ControlBlockWitness: controlBlockWitness,
		Amount:              uint64(minerFee),
		SighashType:         txscript.SigHashAll,
		PriHex:              revealPrivateKey,
	}

	err = metaIdInscribeBuilder.RevealPsbtBuilder.UpdateAndSignTaprootInput([]*common.InputSign{taprootInSigner})
	if err != nil {
		response.Error(c, "签名Reveal交易失败: "+err.Error())
		return
	}

	// 提取 Reveal 交易
	revealTxRaw, err := metaIdInscribeBuilder.RevealPsbtBuilder.ExtractPsbtTransaction()
	if err != nil {
		response.Error(c, "提取Reveal交易失败: "+err.Error())
		return
	}

	// 直接返回交易内容
	response.Success(c, InscribeResponse{
		CommitTxRaw: commitTxRaw,
		RevealTxRaw: revealTxRaw,
		CommitTxID:  commitTx.TxHash().String(),
		MinerFee:    minerFee,
	})
}

// GetInscribe 获取铭文信息
func GetInscribe(c *gin.Context) {
	id := c.Param("id")

	var inscribe models.Inscribe
	if err := models.DB.Preload("Chain").First(&inscribe, id).Error; err != nil {
		response.Error(c, "铭文不存在")
		return
	}

	response.Success(c, inscribe)
}

// ListInscribes 获取铭文列表
func ListInscribes(c *gin.Context) {
	var inscribes []models.Inscribe

	query := models.DB.Preload("Chain")

	// 支持按状态和链ID筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if chainID := c.Query("chain_id"); chainID != "" {
		query = query.Where("chain_id = ?", chainID)
	}

	if err := query.Find(&inscribes).Error; err != nil {
		response.Error(c, "获取铭文列表失败")
		return
	}

	response.Success(c, inscribes)
}

// CommitTxParams commit 交易参数
type CommitTxParams struct {
	ChainType      string  // 链类型
	RevealAddress  string  // reveal 地址
	FromAddress    string  // 支付地址
	RevealFee      int64   // reveal 费用
	NetworkFeeRate int64   // 网络费率
	UtxoList       []*Utxo // UTXO 列表
}

// makeCommitTx 创建 commit 交易
func makeCommitTx(params *CommitTxParams) (*wire.MsgTx, string, int64, error) {
	var (
		utxoAmount        int64 = 0
		txSize            int64 = 0
		minerFee          int64 = 0
		commitPsbtBuilder *common.PsbtBuilder
		inputs            []common.Input      = make([]common.Input, 0)
		outputs           []common.Output     = make([]common.Output, 0)
		inSigns           []*common.InputSign = make([]*common.InputSign, 0)
		err               error
		commitTxRaw       string
		commitTxOutIndex  int64 = 0
	)

	// 从数据库获取地址私钥
	var address models.Address
	if err := models.DB.Where("address = ?", params.FromAddress).First(&address).Error; err != nil {
		return nil, "", 0, fmt.Errorf("获取地址私钥失败: %v", err)
	}

	// 获取地址的公钥脚本
	pkScript, err := common.AddressToPkScript(common.GetNetParams(params.ChainType), params.FromAddress)
	if err != nil {
		return nil, "", 0, fmt.Errorf("生成公钥脚本失败: %v", err)
	}
	fmt.Println("--------------------------------")
	for _, utxo := range params.UtxoList {
		fmt.Printf("utxo: %+v", utxo)
	}
	// 构建输入
	for _, utxo := range params.UtxoList {
		inputs = append(inputs, common.Input{
			OutTxId:  utxo.OutTxId,
			OutIndex: utxo.OutIndex,
		})
	}

	// 构建输出
	outputs = []common.Output{
		{
			Address: params.RevealAddress,
			Amount:  uint64(params.RevealFee),
		},
		{
			Address: params.FromAddress, // 找零地址
			Amount:  0,
		},
	}

	// 创建 PSBT 构建器
	commitPsbtBuilder, err = common.CreatePsbtBuilder(common.GetNetParams(params.ChainType), inputs, outputs)
	if err != nil {
		return nil, "", 0, fmt.Errorf("创建PSBT构建器失败: %v", err)
	}

	// 添加签名信息
	for i, utxo := range params.UtxoList {
		inSigns = append(inSigns, &common.InputSign{
			UtxoType:    common.Witness,
			Index:       int(i),
			PkScript:    pkScript,
			Amount:      uint64(utxo.OutValue),
			SighashType: txscript.SigHashAll,
			PriHex:      address.PrivateKey,
		})
		utxoAmount += utxo.OutValue
	}

	// 更新输入见证数据
	if err = commitPsbtBuilder.UpdateAndAddInputWitness(inSigns); err != nil {
		return nil, "", 0, fmt.Errorf("更新输入见证数据失败: %v", err)
	}

	// 计算交易大小和矿工费
	txSize, err = commitPsbtBuilder.CalTxSize()
	if err != nil {
		return nil, "", 0, fmt.Errorf("计算交易大小失败: %v", err)
	}
	minerFee = txSize * params.NetworkFeeRate

	// 设置找零金额
	commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut[1].Value = int64(utxoAmount - params.RevealFee - minerFee)
	if commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut[1].Value < 546 {
		commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut = commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut[:1]
		commitPsbtBuilder.PsbtUpdater.Upsbt.Outputs = commitPsbtBuilder.PsbtUpdater.Upsbt.Outputs[:1]
	}

	// 签名交易
	if err = commitPsbtBuilder.UpdateAndSignInput(inSigns); err != nil {
		return nil, "", 0, fmt.Errorf("签名交易失败: %v", err)
	}

	// 提取交易
	commitTxRaw, err = commitPsbtBuilder.ExtractPsbtTransaction()
	if err != nil {
		return nil, "", 0, fmt.Errorf("提取交易失败: %v", err)
	}

	// 反序列化交易
	commitTx := wire.NewMsgTx(wire.TxVersion)
	b, _ := hex.DecodeString(commitTxRaw)
	buf := bytes.NewBuffer(b)
	if err = commitTx.Deserialize(buf); err != nil {
		return nil, "", 0, fmt.Errorf("反序列化交易失败: %v", err)
	}

	return commitTx, commitTxRaw, commitTxOutIndex, nil
}
