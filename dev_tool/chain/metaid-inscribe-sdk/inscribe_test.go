package metaid_inscribe_sdk

import (
	"bytes"
	"dev_tool/chain/metaid-inscribe-sdk/common"
	"dev_tool/chain/metaid-inscribe-sdk/inscribe"
	"dev_tool/chain/metaid-inscribe-sdk/tool"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func TestMetaIdInscribe(t *testing.T) {
	var (
		//以下是可以改的
		net         string                 = "regtest"
		address     string                 = "bcrt1qkn6ca856lrptq3j5g0caq5dhugycz5208d4mrc"
		metaIdFlag  string                 = "testid"
		path        string                 = "/protocols/rollup"
		metaDataMap map[string]interface{} = map[string]interface{}{
			"test": "1", //内容
		}
		pinOutValue int64                = 546
		feeRate     int64                = 6
		otherOuts   []*inscribe.OtherOut = []*inscribe.OtherOut{}

		//不用动的
		metaIdOpRequest *inscribe.MetaIdInscribeRequest

		metaIdInscribeBuilder *inscribe.MetaIdBuilder
		minerFee              int64 = 0
		err                   error

		revealPrePsbtRaw     string
		revealPrivateKey     string = ""
		revealAddress        string
		revealInputIndex     int    = 0
		revealPkScript       string = ""
		redeemScript         string = ""
		controlBlockWitness  string = ""
		taprootInSigner      *common.InputSign
		revealTaprootSigners []*common.InputSign = make([]*common.InputSign, 0)
		revealTxRaw          string
	)

	metaIdOpRequest = &inscribe.MetaIdInscribeRequest{
		Net:           common.GetNetParams(net),
		MetaIdFlag:    metaIdFlag,
		Path:          path,
		Payload:       tool.AnyToStr(metaDataMap),
		PinOutValue:   pinOutValue,
		PinOutAddress: address,
		ChangeAddress: address,
		OtherOuts:     otherOuts,
	}

	metaIdInscribeBuilder, minerFee, err = inscribe.MetaIdInscribeBuilder(metaIdOpRequest, feeRate)
	if err != nil {
		t.Fatalf("MetaIdInscribeBuilder error: %v", err)
	}
	revealPrePsbtRaw, err = metaIdInscribeBuilder.RevealPsbtBuilder.ToString()
	if err != nil {
		t.Fatalf("RevealPsbtBuilder ToString error: %v", err)
	}
	_ = revealPrePsbtRaw
	revealPrivateKey = metaIdInscribeBuilder.RevealPrivateKeyHex
	revealAddress = metaIdInscribeBuilder.RevealAddress
	revealPkScript = hex.EncodeToString(metaIdInscribeBuilder.TxCtxData.CommitTxAddressPkScript)
	revealInputIndex = int(metaIdInscribeBuilder.RevealTaprootDataInputIndex)
	redeemScript = hex.EncodeToString(metaIdInscribeBuilder.TxCtxData.InscriptionScript)
	controlBlockWitness = hex.EncodeToString(metaIdInscribeBuilder.TxCtxData.ControlBlockWitness)
	//make commit tx
	commitTx, commitTxRaw, commitTxOutIndex, err := makeCommitTx(revealAddress, minerFee, feeRate)
	fmt.Println("commitTxHash:", commitTx.TxHash().String())
	//complete reveal psbt and sign
	txHash := commitTx.TxHash()
	commitPreOutPoint := wire.NewOutPoint(&txHash, uint32(commitTxOutIndex))
	metaIdInscribeBuilder.RevealPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxIn[revealInputIndex].PreviousOutPoint = *commitPreOutPoint

	taprootInSigner = &common.InputSign{
		UtxoType:            common.Taproot,
		Index:               revealInputIndex,
		PkScript:            revealPkScript,
		RedeemScript:        redeemScript,
		ControlBlockWitness: controlBlockWitness,
		Amount:              uint64(minerFee),
		SighashType:         txscript.SigHashAll,
		PriHex:              revealPrivateKey,
	}
	fmt.Printf("%+v", taprootInSigner)
	revealTaprootSigners = append(revealTaprootSigners, taprootInSigner)
	err = metaIdInscribeBuilder.RevealPsbtBuilder.UpdateAndSignTaprootInput(revealTaprootSigners)
	if err != nil {
		t.Fatalf("UpdateAndSignTaprootInput error: %v", err)
	}

	revealTxRaw, err = metaIdInscribeBuilder.RevealPsbtBuilder.ExtractPsbtTransaction()
	if err != nil {
		t.Fatalf("ExtractPsbtTransaction error: %v", err)
	}
	fmt.Printf("CommitTxRaw: %s\n", commitTxRaw)
	fmt.Printf("RevealTxRaw: %s\n", revealTxRaw)
}

type Utxo struct {
	OutTxId  string
	OutIndex uint32
	OutValue int64
	OutRaw   string
}

func makeCommitTx(revealAddress string, revealFee, networkFeeRate int64) (*wire.MsgTx, string, int64, error) {
	var (
		//以下是可以改的
		net           string  = "regtest"
		privateKeyHex string  = "214c743547d1d7a21f331bd8a06f06dc8d9e4d24b0dce3194d6a94b5b6b03a22"
		address       string  = "bcrt1qkn6ca856lrptq3j5g0caq5dhugycz5208d4mrc"
		changeAddress         = address
		pkScript, _           = common.AddressToPkScript(common.GetNetParams(net), address)
		utxoList      []*Utxo = []*Utxo{
			{
				OutTxId:  "9b586521a95619762731fafdfedb9031dc469fac62d6dc6fe046ccd67c7c1ffe",
				OutIndex: 1,
				OutValue: 49784644,
			},
		}

		//不用动的
		utxoAmount        int64 = 0
		txSize            int64 = 0
		minerFee          int64 = 0
		commitPsbtBuilder *common.PsbtBuilder
		inputs            []common.Input      = make([]common.Input, 0)
		outputs           []common.Output     = make([]common.Output, 0)
		inSigns           []*common.InputSign = make([]*common.InputSign, 0)
		err               error
		commitTxRaw       string
		commitTxOutIndex  int64
	)
	//make commit tx
	for _, utxo := range utxoList {
		inputs = append(inputs, common.Input{
			OutTxId:  utxo.OutTxId,
			OutIndex: utxo.OutIndex,
		})
	}
	outputs = []common.Output{
		{
			Address: revealAddress,
			Amount:  uint64(revealFee),
			//Script:  "",
		},
		{
			Address: changeAddress,
			Amount:  0,
		},
	}

	commitPsbtBuilder, err = common.CreatePsbtBuilder(common.GetNetParams(net), inputs, outputs)
	if err != nil {
		return nil, "", 0, err
	}

	for i, utxo := range utxoList {
		inSigns = append(inSigns, &common.InputSign{
			UtxoType: common.Witness,
			Index:    int(i),
			//OutRaw:         utxo.OutRaw,
			PkScript:    pkScript,
			Amount:      uint64(utxo.OutValue),
			SighashType: txscript.SigHashAll,
			PriHex:      privateKeyHex,
		})
		utxoAmount += utxo.OutValue
	}

	if err = commitPsbtBuilder.UpdateAndAddInputWitness(inSigns); err != nil {
		return nil, "", 0, err
	}

	txSize, err = commitPsbtBuilder.CalTxSize()
	if err != nil {
		return nil, "", 0, err
	}
	minerFee = txSize * networkFeeRate

	commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut[1].Value = int64(utxoAmount - revealFee - minerFee)
	if commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut[1].Value < 546 {
		commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut = commitPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxOut[:1]
		commitPsbtBuilder.PsbtUpdater.Upsbt.Outputs = commitPsbtBuilder.PsbtUpdater.Upsbt.Outputs[:1]
	}
	if err = commitPsbtBuilder.UpdateAndSignInput(inSigns); err != nil {
		return nil, "", 0, err
	}
	commitTxRaw, err = commitPsbtBuilder.ExtractPsbtTransaction()
	if err != nil {
		return nil, "", 0, err
	}
	commitTxOutIndex = 0
	commitTx := wire.NewMsgTx(wire.TxVersion)
	//buf := bytes.NewBufferString(commitTxRaw)
	b, _ := hex.DecodeString(commitTxRaw)
	buf := bytes.NewBufferString(string(b))
	err = commitTx.Deserialize(buf)
	if err != nil {
		return nil, "", 0, err
	}
	return commitTx, commitTxRaw, commitTxOutIndex, nil
}
