package models

import (
	"github.com/ethereum/go-ethereum/common"
	"pledge-backend/db"
	"pledge-backend/log"
)

type TxReceiptData struct {
	Id                int    `json:"-" gorm:"column:id;primaryKey"`
	ChainId           int64  `json:"chainId" gorm:"column:chain_id"`
	TxHash            string `json:"txHash" gorm:"column:tx_hash"`
	TxIndex           uint   `json:"TxIndex" gorm:"column:tx_index"`
	BlockHash         string `json:"blockHash" gorm:"column:block_hash"`
	BlockNumber       int64  `json:"blockNumber" gorm:"column:block_number"`
	Status            uint64 `json:"status" gorm:"column:status"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gorm:"column:cumulative_gas_used"`
	GasUsed           uint64 `json:"gasUsed" gorm:"column:gas_used"`
	ContractAddress   string `json:"contractAddress" gorm:"column:contract_address;"`
	Logs              string `json:"logs" gorm:"column:logs"`
	LogsBloom         string `json:"logsBloom" gorm:"column:logs_bloom"`
	Type              uint8  `json:"type" gorm:"column:type"`
	PostState         string `json:"postState" gorm:"column:post_state"`
}

type TxReceiptDataInfoRes struct {
	TxReceiptData TxReceiptData `json:"txReceiptData"`
}

func NewTxReceiptData() *TxReceiptData {
	return &TxReceiptData{}
}

func (tx *TxReceiptData) TableName() string {
	return "tx_receipt_data"
}

func (tx *TxReceiptData) TxReceiptDataInfo(chainId int64, txHash common.Hash, res *TxReceiptDataInfoRes) error {
	var txReceiptData TxReceiptData

	err := db.Mysql.Table("tx_receipt_data").Where("chain_id=?", chainId).Where("tx_hash", txHash.Hex()).Find(&txReceiptData).Debug().Error
	if err != nil {
		log.Logger.Error("Failed to query the transaction receipt information from db:" + err.Error())
		return err
	}

	if txReceiptData.TxHash != "" {
		*res = TxReceiptDataInfoRes{
			TxReceiptData: txReceiptData,
		}
	}

	return nil
}

func (tx *TxReceiptData) CreateTxReceiptDataInfo(txReceiptData *TxReceiptData, res *TxReceiptDataInfoRes) error {
	err := db.Mysql.Table("tx_receipt_data").Create(txReceiptData).Error
	if err != nil {
		log.Logger.Error("Failed to add the transaction receipt information :" + err.Error())
		return err
	}
	*res = TxReceiptDataInfoRes{
		TxReceiptData: *txReceiptData,
	}
	return nil
}
