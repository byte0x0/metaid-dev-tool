package inscribe

import (
	"dev_tool/chain/metaid-inscribe-sdk/common"
	"encoding/hex"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type MetaIdBuilder struct {
	Net        *chaincfg.Params
	MetaIdData *MetaIdData

	FeeRate   int64
	OtherOuts []*OtherOut

	RevealPrivateKeyHex string
	RevealAddress       string

	RevealTaprootDataInputIndex uint32
	TxCtxData                   *inscriptionTxCtxData
	RevealPsbtBuilder           *common.PsbtBuilder
	revealTx                    *wire.MsgTx

	metaIdOutValue   int64
	metaIdOutAddress string
}

func (m *MetaIdBuilder) buildEmptyRevealPsbt() error {
	var (
		revealPsbtBuilder     *common.PsbtBuilder
		inputs                []common.Input      = make([]common.Input, 0)
		inSigners             []*common.InputSign = make([]*common.InputSign, 0)
		outputs               []common.Output     = make([]common.Output, 0)
		taprootDataInputIndex uint32              = 0
		err                   error
	)

	emptyTxId := "0000000000000000000000000000000000000000000000000000000000000000"
	taprootDataIn := common.Input{
		OutTxId:  emptyTxId,
		OutIndex: 0,
	}

	inputs = append(inputs, taprootDataIn)

	outPin := common.Output{
		Address: m.metaIdOutAddress,
		Amount:  uint64(m.metaIdOutValue),
	}
	outputs = append(outputs, outPin)

	if m.OtherOuts != nil && len(m.OtherOuts) != 0 {
		for _, v := range m.OtherOuts {
			out := common.Output{
				Address: v.Address,
				Amount:  uint64(v.Amount),
				Script:  v.Script,
			}
			outputs = append(outputs, out)
		}
	}
	revealPsbtBuilder, err = common.CreatePsbtBuilder(m.Net, inputs, outputs)
	if err != nil {
		return err
	}
	m.RevealPsbtBuilder = revealPsbtBuilder

	taprootDataInSigner := &common.InputSign{
		UtxoType:            common.Taproot,
		Index:               int(taprootDataInputIndex),
		PkScript:            hex.EncodeToString(m.TxCtxData.CommitTxAddressPkScript),
		RedeemScript:        hex.EncodeToString(m.TxCtxData.InscriptionScript),
		ControlBlockWitness: hex.EncodeToString(m.TxCtxData.ControlBlockWitness),
		Amount:              uint64(m.CalRevealPsbtFee(m.FeeRate)),
		SighashType:         txscript.SigHashAll,
		PriHex:              "",
	}
	inSigners = append(inSigners, taprootDataInSigner)

	err = revealPsbtBuilder.UpdateAndAddInputWitness(inSigners)
	if err != nil {
		return err
	}

	m.RevealAddress, err = common.PkScriptToAddress(m.Net, hex.EncodeToString(m.TxCtxData.CommitTxAddressPkScript))
	if err != nil {
		return err
	}
	m.RevealPrivateKeyHex = m.TxCtxData.RecoveryPrivateKeyHex

	m.RevealPsbtBuilder = revealPsbtBuilder
	m.RevealTaprootDataInputIndex = taprootDataInputIndex
	m.TxCtxData.revealTxPrevOutput = &wire.TxOut{
		PkScript: m.TxCtxData.CommitTxAddressPkScript,
		Value:    m.CalRevealPsbtFee(m.FeeRate),
	}
	return nil
}

func (m *MetaIdBuilder) CalRevealPsbtFee(feeRate int64) int64 {
	var (
		tx          *wire.MsgTx = m.RevealPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx
		txTotalSize int         = tx.SerializeSize()
		txBaseSize  int         = tx.SerializeSizeStripped()
		txFee       int64       = 0
		weight      int64       = 0
		vSize       int64       = 0

		revealOutValues = int64(0)
	)
	revealOutValues += m.metaIdOutValue
	if m.OtherOuts != nil && len(m.OtherOuts) > 0 {
		for _, v := range m.OtherOuts {
			revealOutValues += v.Amount
		}
	}

	emptySignature := make([]byte, 64)
	emptyControlBlockWitness := make([]byte, 33)
	txTotalSize += wire.TxWitness{emptySignature, m.TxCtxData.InscriptionScript, emptyControlBlockWitness}.SerializeSize()

	weight = int64(txBaseSize*3 + txTotalSize)
	vSize = (weight + (blockchain.WitnessScaleFactor - 1)) / blockchain.WitnessScaleFactor
	vSize = vSize + 1
	txFee = vSize * feeRate
	return txFee + revealOutValues
}

func (m *MetaIdBuilder) completeRevealPsbt(commitTxId string, commitTxOutIndex uint32) error {
	var (
		commitPreOutPoint *wire.OutPoint
		txHash            *chainhash.Hash
		err               error
	)
	txHash, err = chainhash.NewHashFromStr(commitTxId)
	if err != nil {
		return err
	}
	commitPreOutPoint = wire.NewOutPoint(txHash, commitTxOutIndex)
	m.RevealPsbtBuilder.PsbtUpdater.Upsbt.UnsignedTx.TxIn[m.RevealTaprootDataInputIndex].PreviousOutPoint = *commitPreOutPoint
	return nil
}

func (m *MetaIdBuilder) signRevealPsbt(taprootInSigner *common.InputSign) error {
	var (
		revealSigners        []*common.InputSign = make([]*common.InputSign, 0)
		revealTaprootSigners []*common.InputSign = make([]*common.InputSign, 0)
		err                  error
	)

	err = m.RevealPsbtBuilder.UpdateAndSignInput(revealSigners)
	if err != nil {
		return err
	}

	if taprootInSigner == nil {
		taprootInSigner = &common.InputSign{
			UtxoType:            common.Taproot,
			Index:               int(m.RevealTaprootDataInputIndex),
			PkScript:            hex.EncodeToString(m.TxCtxData.CommitTxAddressPkScript),
			RedeemScript:        hex.EncodeToString(m.TxCtxData.InscriptionScript),
			ControlBlockWitness: hex.EncodeToString(m.TxCtxData.ControlBlockWitness),
			Amount:              uint64(m.CalRevealPsbtFee(m.FeeRate)),
			SighashType:         txscript.SigHashAll,
			PriHex:              m.TxCtxData.RecoveryPrivateKeyHex,
		}
		revealTaprootSigners = append(revealTaprootSigners, taprootInSigner)
	}

	err = m.RevealPsbtBuilder.UpdateAndSignTaprootInput(revealTaprootSigners)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetaIdBuilder) ExtractRevealTransaction() (string, string, error) {
	var (
		commitTxHex string
		revealTxHex string
		err         error
	)

	revealTxHex, err = m.RevealPsbtBuilder.ExtractPsbtTransaction()
	if err != nil {
		return "", "", err
	}
	return commitTxHex, revealTxHex, nil
}

func createMetaIdTxCtxData(net *chaincfg.Params, metaIdData *MetaIdData) (*inscriptionTxCtxData, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	inscriptionBuilder := txscript.NewScriptBuilder().
		AddData(schnorr.SerializePubKey(privateKey.PubKey())).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_FALSE).
		AddOp(txscript.OP_IF).
		AddData([]byte(metaIdData.MetaIDFlag)). //<metaid_flag>
		AddData([]byte(metaIdData.Operation))   //<operation>

	inscriptionBuilder.AddData([]byte(metaIdData.Path)) //<path>
	if metaIdData.Encryption == "" {
		inscriptionBuilder.AddOp(txscript.OP_0)
	} else {
		inscriptionBuilder.AddData([]byte(metaIdData.Encryption)) //<Encryption>
	}

	if metaIdData.Version == "" {
		inscriptionBuilder.AddOp(txscript.OP_0)
	} else {
		inscriptionBuilder.AddData([]byte(metaIdData.Version)) //<version>
	}

	if metaIdData.ContentType == "" {
		inscriptionBuilder.AddOp(txscript.OP_0)
	} else {
		inscriptionBuilder.AddData([]byte(metaIdData.ContentType)) //<content-type>
	}
	maxChunkSize := 520
	bodySize := len(metaIdData.Content)
	for i := 0; i < bodySize; i += maxChunkSize {
		end := i + maxChunkSize
		if end > bodySize {
			end = bodySize
		}
		inscriptionBuilder.AddFullData(metaIdData.Content[i:end]) //<payload>
	}

	inscriptionScript, err := inscriptionBuilder.Script()
	if err != nil {
		return nil, err
	}
	inscriptionScript = append(inscriptionScript, txscript.OP_ENDIF)

	proof := &txscript.TapscriptProof{
		TapLeaf:  txscript.NewBaseTapLeaf(schnorr.SerializePubKey(privateKey.PubKey())),
		RootNode: txscript.NewBaseTapLeaf(inscriptionScript),
	}

	controlBlock := proof.ToControlBlock(privateKey.PubKey())
	controlBlockWitness, err := controlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	tapHash := proof.RootNode.TapHash()
	commitTxAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(privateKey.PubKey(), tapHash[:])), net)
	if err != nil {
		return nil, err
	}
	commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
	if err != nil {
		return nil, err
	}

	recoveryPrivateKeyWIF, err := btcutil.NewWIF(txscript.TweakTaprootPrivKey(*privateKey, tapHash[:]), net, true)
	if err != nil {
		return nil, err
	}

	recoveryPrivateKeyHex := hex.EncodeToString(privateKey.Serialize())

	return &inscriptionTxCtxData{
		privateKey:              privateKey,
		InscriptionScript:       inscriptionScript,
		CommitTxAddressPkScript: commitTxAddressPkScript,
		ControlBlockWitness:     controlBlockWitness,
		recoveryPrivateKeyWIF:   recoveryPrivateKeyWIF.String(),
		RecoveryPrivateKeyHex:   recoveryPrivateKeyHex,
	}, nil
}
