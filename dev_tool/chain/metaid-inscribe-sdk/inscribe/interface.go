package inscribe

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
)

type MetaIdData struct {
	MetaIDFlag  string
	Operation   string
	Path        string
	Content     []byte
	Encryption  string
	Version     string
	ContentType string
}

type inscriptionTxCtxData struct {
	privateKey              *btcec.PrivateKey
	InscriptionScript       []byte
	CommitTxAddressPkScript []byte
	ControlBlockWitness     []byte
	recoveryPrivateKeyWIF   string
	RecoveryPrivateKeyHex   string
	revealTxPrevOutput      *wire.TxOut
}

type OtherOut struct {
	Address string
	Amount  int64
	Script  string
}
