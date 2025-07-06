package inscribe

import "github.com/btcsuite/btcd/chaincfg"

type MetaIdInscribeRequest struct {
	Net           *chaincfg.Params
	MetaIdFlag    string
	Operation     string
	Path          string
	Payload       string
	PinOutValue   int64
	PinOutAddress string
	ChangeAddress string

	OtherOuts []*OtherOut
}

func MetaIdInscribeBuilder(opRep *MetaIdInscribeRequest, feeRate int64) (*MetaIdBuilder, int64, error) {
	var (
		err           error
		metaIdBuilder *MetaIdBuilder
		fee           int64 = 0

		content                = opRep.Payload
		path                   = opRep.Path
		metaIdData *MetaIdData = &MetaIdData{
			MetaIDFlag:  opRep.MetaIdFlag,
			Operation:   opRep.Operation,
			Path:        path,
			Content:     []byte(content),
			Encryption:  "",
			Version:     "",
			ContentType: "application/json",
		}
	)
	metaIdBuilder = &MetaIdBuilder{
		Net:              opRep.Net,
		MetaIdData:       metaIdData,
		FeeRate:          feeRate,
		OtherOuts:        opRep.OtherOuts,
		metaIdOutValue:   opRep.PinOutValue,
		metaIdOutAddress: opRep.PinOutAddress,
	}

	txCtxData, err := createMetaIdTxCtxData(opRep.Net, metaIdBuilder.MetaIdData)
	if err != nil {
		return nil, 0, err
	}
	metaIdBuilder.TxCtxData = txCtxData

	err = metaIdBuilder.buildEmptyRevealPsbt()
	if err != nil {
		return nil, 0, err
	}
	fee = metaIdBuilder.CalRevealPsbtFee(feeRate)
	return metaIdBuilder, fee, nil
}
