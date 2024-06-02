package db

import (
	"math/big"
)

type Tx struct {
	BlNumber      uint64   `json:"bl_n"`
	BlDT          uint64   `json:"bl_dt"`
	TxFrom        string   `json:"tx_from"`
	TxTo          string   `json:"tx_to"`
	TxValue       *big.Int `json:"tx_value"`
	TxValueFloat  float64  `json:"tx_value_float"`
	TxValid       bool     `json:"tx_valid"`
	TxStatus      uint64   `json:"tx_status"`
	TxGasUsed     uint64   `json:"tx_gas_used"`
	TxNewContract string   `json:"tx_new_contract"`
	TxDataMethod  []byte   `json:"tx_data_method"`
	TxDataP1      []byte   `json:"tx_data_p1"`
	TxDataP2      []byte   `json:"tx_data_p2"`
	TxDataP3      []byte   `json:"tx_data_p3"`
	TxDataP4      []byte   `json:"tx_data_p4"`
}

type TxsByMinutes struct {
	Items []*TxsByMinute
}

type TxsByMinute struct {
	DT  uint64
	TXS []*Tx
}
