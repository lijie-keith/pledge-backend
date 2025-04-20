package models

import (
	"pledge-backend/db"
	"pledge-backend/log"
)

type BlockData struct {
	Id          int    `json:"-" gorm:"column:id;primaryKey"`
	ChainId     int64  `json:"chainId" gorm:"column:chain_id"`
	BlockNum    int    `json:"blockNum" gorm:"column:block_num"`
	ParentHash  string `json:"parentHash" gorm:"column:parent_hash"`
	UncleHash   string `json:"uncleHash" gorm:"column:uncle_hash"`
	Coinbase    string `json:"coinbase" gorm:"column:coinbase"`
	Root        string `json:"root" gorm:"column:root"`
	TxHash      string `json:"txHash" gorm:"column:tx_hash"`
	ReceiptHash string `json:"receiptHash" gorm:"column:receipt_hash"`
	Bloom       string `json:"bloom" gorm:"column:bloom"`
	Difficulty  int    `json:"difficulty" gorm:"column:difficulty"`
	GasLimit    uint64 `json:"gasLimit" gorm:"column:gas_limit"`
	GasUsed     uint64 `json:"gasUsed" gorm:"column:gas_used"`
	Time        uint64 `json:"time" gorm:"column:time"`
	Extra       string `json:"extra" gorm:"column:extra"`
	MixDigest   string `json:"mixDigest" gorm:"column:mix_digest"`
	Nonce       uint64 `json:"nonce" gorm:"column:nonce"`
	BaseFee     int    `json:"baseFee" gorm:"column:base_fee"` // BaseFee was added by EIP-1559 and is ignored in legacy headers.
}

type BlockDataInfoRes struct {
	BlockData     BlockData       `json:"blockData"`
	TxDataInfoRes []TxDataInfoRes `json:"txDataInfoRes"`
}

func NewBlockData() *BlockData {
	return &BlockData{}
}

func (tx *BlockData) TableName() string {
	return "block_data"
}

func (tx *BlockData) BlockDataInfo(chainId int64, blockNum int64, res *BlockDataInfoRes) error {
	var blockData BlockData

	err := db.Mysql.Table("block_data").Where("chain_id=?", chainId).Where("block_num", blockNum).Find(&blockData).Debug().Error
	if err != nil {
		log.Logger.Error("Failed to query the block information from db:" + err.Error())
		return err
	}

	if blockData.BlockNum != 0 {
		*res = BlockDataInfoRes{
			BlockData: blockData,
		}
	}

	return nil
}

func (tx *BlockData) CreateBlockDataInfo(blockData *BlockData, res *BlockDataInfoRes) error {
	err := db.Mysql.Table("block_data").Create(blockData).Error
	if err != nil {
		log.Logger.Error("Failed to add the block information :" + err.Error())
		return err
	}
	*res = BlockDataInfoRes{
		BlockData: *blockData,
	}
	return nil
}
