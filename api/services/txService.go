package services

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/config"
	"pledge-backend/log"
)

type txService struct{}

func NewTx() *txService {
	return &txService{}
}

func (s *txService) TxDataInfo(chainId int64, txHash common.Hash, result *models.TxDataInfoRes) int {
	err := models.NewTxData().TxDataInfo(chainId, txHash, result)
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr
	}
	if result.TxData.TxHash == "" {
		//1.查询链上数据
		txData, err := getTxDataFromChain(txHash, chainId)
		if err != nil {
			return statecode.CommonErrServerErr
		}
		//2.链上数据入库
		err = models.NewTxData().CreateTxDataInfo(txData, result)
		if err != nil {
			return statecode.CommonErrServerErr
		}
	}
	return statecode.CommonSuccess
}

func getTxDataFromChain(txHash common.Hash, chainId int64) (*models.TxData, error) {
	ethereumConn, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error("Failed to establish the connection:" + err.Error())
		return nil, err
	}
	tx, _, err := ethereumConn.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Logger.Error("Failed to query the transaction information from blockchain:" + err.Error())
		return nil, err
	}

	var from common.Address
	sender, err := types.LatestSignerForChainID(tx.ChainId()).Sender(tx)

	if err != nil {
		log.Logger.Error("Failed to query the from address:" + err.Error())
		return nil, err
	}
	from = sender

	return ProcessTransaction(tx, from, chainId)
}

func ProcessTransaction(tx *types.Transaction, from common.Address, chainId int64) (*models.TxData, error) {
	var to common.Address
	if tx.To() != nil {
		to = *tx.To()
	}
	txData := &models.TxData{
		TxHash:   tx.Hash().Hex(),
		ChainId:  chainId,
		Nonce:    tx.Nonce(),
		From:     from.Hex(),
		To:       to.Hex(),
		Value:    tx.Value().Int64(),
		GasLimit: tx.Gas(),
		GasPrice: tx.GasPrice().Int64(),
		//Data:     string(tx.Data()),
	}
	switch tx.Type() {
	case types.LegacyTxType:
		// Legacy 交易字段
		return txData, nil
	case types.AccessListTxType:
		// EIP-2930 交易字段
		return txData, nil
	case types.DynamicFeeTxType:
		// EIP-1559 交易字段
		txData.MaxPriorityFeePerGas = tx.GasTipCap().Int64()
		txData.MaxFeePerGas = tx.GasFeeCap().Int64()
		return txData, nil
	default:
		return nil, errors.New("unknown tx type")
	}
}

func (s *txService) GetTxDataByTx(tx *types.Transaction, chainId int64) (*models.TxData, error) {
	var from common.Address
	sender, err := types.LatestSignerForChainID(tx.ChainId()).Sender(tx)

	if err != nil {
		log.Logger.Error("Failed to query the from address:" + err.Error())
		return nil, err
	}
	from = sender

	return ProcessTransaction(tx, from, chainId)
}
