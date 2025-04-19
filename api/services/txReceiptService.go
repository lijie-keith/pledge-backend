package services

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/config"
	"pledge-backend/log"
)

type txReceiptService struct{}

func NewTxReceipt() *txReceiptService {
	return &txReceiptService{}
}

func (s *txReceiptService) TxReceiptDataInfo(chainId int64, txHash common.Hash, result *models.TxReceiptDataInfoRes) int {
	err := models.NewTxReceiptData().TxReceiptDataInfo(chainId, txHash, result)
	if err != nil {
		return statecode.CommonErrServerErr
	}
	if result.TxReceiptData.TxHash == "" {
		//1.查询链上数据
		txReceiptData, err := getTxReceiptDataFromChain(txHash, chainId)
		if err != nil {
			return statecode.CommonErrServerErr
		}
		//2.链上数据入库
		err = models.NewTxReceiptData().CreateTxReceiptDataInfo(txReceiptData, result)
		if err != nil {
			return statecode.CommonErrServerErr
		}
	}
	return statecode.CommonSuccess
}

func getTxReceiptDataFromChain(txHash common.Hash, chainId int64) (*models.TxReceiptData, error) {
	ethereumConn, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error("Failed to establish the connection:" + err.Error())
		return nil, err
	}
	receipt, err := ethereumConn.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Logger.Error("Failed to query the transaction receipt information from blockChain:" + err.Error())
		return nil, err
	}

	return &models.TxReceiptData{
		ChainId:           chainId,
		TxHash:            receipt.TxHash.Hex(),
		TxIndex:           receipt.TransactionIndex,
		BlockHash:         receipt.BlockHash.Hex(),
		BlockNumber:       receipt.BlockNumber.Int64(),
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		GasUsed:           receipt.GasUsed,
		ContractAddress:   receipt.ContractAddress.Hex(),
		Type:              receipt.Type,
		PostState:         string(receipt.PostState),
	}, nil
}
