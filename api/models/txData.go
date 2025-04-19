package models

import (
	"github.com/ethereum/go-ethereum/common"
	"pledge-backend/db"
	"pledge-backend/log"
)

type TxData struct {
	Id                   int    `json:"-" gorm:"column:id;primaryKey"`
	ChainId              int64  `json:"chainId" gorm:"column:chain_id"`
	Nonce                uint64 `json:"nonce" gorm:"column:nonce"`
	TxHash               string `json:"txHash" gorm:"column:tx_hash;"`
	From                 string `json:"from" gorm:"column:from;"`
	To                   string `json:"to" gorm:"column:to;"`
	Value                int64  `json:"value" gorm:"column:value;"`
	GasPrice             int64  `json:"gasPrice" gorm:"column:gas_price;"`
	GasLimit             uint64 `json:"gasLimit" gorm:"column:gas_limit;"`
	Data                 string `json:"data" gorm:"column:data"`
	MaxFeePerGas         int64  `json:"maxFeePerGas" gorm:"column:max_fee_per_gas;"`
	MaxPriorityFeePerGas int64  `json:"maxPriorityFeePerGas" gorm:"column:max_priority_fee_per_gas;"`
}

type TxDataInfoRes struct {
	TxData TxData `json:"txData"`
}

func NewTxData() *TxData {
	return &TxData{}
}

func (tx *TxData) TableName() string {
	return "tx_data"
}

func (tx *TxData) TxDataInfo(chainId int64, txHash common.Hash, res *TxDataInfoRes) error {
	var txData TxData

	err := db.Mysql.Table("tx_data").Where("chain_id=?", chainId).Where("tx_hash", txHash.Hex()).Find(&txData).Debug().Error
	if err != nil {
		log.Logger.Error("Failed to query the transaction information from db:" + err.Error())
		return err
	}

	if txData.TxHash != "" {
		*res = TxDataInfoRes{
			TxData: txData,
		}
	}

	return nil
}

func (tx *TxData) CreateTxDataInfo(txData *TxData, res *TxDataInfoRes) error {
	err := db.Mysql.Table("tx_data").Create(txData).Error
	if err != nil {
		log.Logger.Error("Failed to add the transaction information: " + err.Error())
		return err
	}
	*res = TxDataInfoRes{
		TxData: *txData,
	}
	return nil
}
