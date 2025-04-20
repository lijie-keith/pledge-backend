package services

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"pledge-backend/api/common/redisKey"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/config"
	"pledge-backend/db"
	"pledge-backend/log"
	"strconv"
)

type blockService struct{}

func NewBlock() *blockService {
	return &blockService{}
}

func (s *blockService) BlockDataInfo(chainId int64, blockNumStr string, full bool, result *models.BlockDataInfoRes) int {
	var err error
	switch blockNumStr {
	case "head":
		blockType := blockNumStr
		err = getBlockInfoByBlockType(nil, blockType, full, result)
	case "finalized":
		blockNum := rpc.FinalizedBlockNumber.Int64()
		blockType := blockNumStr
		err = getBlockInfoByBlockType(&blockNum, blockType, full, result)
	case "safe":
		blockNum := rpc.SafeBlockNumber.Int64()
		blockType := blockNumStr
		err = getBlockInfoByBlockType(&blockNum, blockType, full, result)
	default:
		err = getBlockInfoByBlockNum(chainId, blockNumStr, full, result)
	}

	if err != nil {
		return statecode.CommonErrServerErr
	}

	return statecode.CommonSuccess
}

func getBlockInfoByBlockNum(chainId int64, blockNumStr string, full bool, result *models.BlockDataInfoRes) error {
	blockNum, err := strconv.ParseInt(blockNumStr, 10, 64)
	if err != nil {
		log.Logger.Error("The wrong block number: " + err.Error())
		return err
	}

	//1.查询数据库
	err = models.NewBlockData().BlockDataInfo(chainId, blockNum, result)
	if err != nil {
		return err
	}
	//2.从链上查询数据
	if result.BlockData.BlockNum == 0 {
		err := NewBlock().GetHeaderInfoFromBlockChain(&blockNum, result)
		if err != nil {
			return err
		}
	}

	// 3.判断是否查询交易信息
	if full {
		//4.这里应该也可以从数据库获取交易信息,但这里直接从链上查询交易数据
		err := NewBlock().GetTxInfoFromBlockChain(result)
		if err != nil {
			return err
		}
	}

	return nil
}

func getBlockInfoByBlockType(blockNum *int64, blockType string, full bool, result *models.BlockDataInfoRes) error {
	blockHeadByte, err := db.RedisGet(redisKey.BlockKey + blockType)
	if err != nil {
		log.Logger.Warn("Failed to query block head information from redis :" + err.Error())
	}
	if blockHeadByte == nil || len(blockHeadByte) == 0 {
		//1.从链上查询head信息
		err := NewBlock().GetHeaderInfoFromBlockChain(blockNum, result)
		if err != nil {
			return err
		}

		//2.总会从链上获取该区块交易数据
		err = NewBlock().GetTxInfoFromBlockChain(result)
		if err != nil {
			return err
		}

		// 3.这里放入redis失败其实不影响业务数据返回，所以不返回err，只记录
		err = db.RedisSet(redisKey.BlockKey+blockType, result, 0)
		if err != nil {
			log.Logger.Error("Failed to set block head information into redis :" + err.Error())
		}
	} else {
		err := json.Unmarshal(blockHeadByte, &result)
		if err != nil {
			return err
		}
	}

	if !full {
		result.TxDataInfoRes = nil
	}

	return nil
}

func (s *blockService) GetHeaderInfoFromBlockChain(blockNum *int64, result *models.BlockDataInfoRes) error {
	ethereumConn, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error("Failed to establish the connection :" + err.Error())
		return err
	}
	defer ethereumConn.Close()
	var blockHeaderData *types.Header
	if blockNum == nil {
		blockHeaderData, err = ethereumConn.HeaderByNumber(context.Background(), nil)
	} else {
		blockHeaderData, err = ethereumConn.HeaderByNumber(context.Background(), big.NewInt(*blockNum))
	}

	if err != nil {
		log.Logger.Error("Failed to query block head information from blockChain :" + err.Error())
		return err
	}
	chainId, err := ethereumConn.ChainID(context.Background())
	if err != nil {
		log.Logger.Error("Failed to get chainId :" + err.Error())
	}
	err = models.NewBlockData().CreateBlockDataInfo(&models.BlockData{
		BlockNum:    int(blockHeaderData.Number.Int64()),
		ChainId:     chainId.Int64(),
		ParentHash:  blockHeaderData.ParentHash.String(),
		UncleHash:   blockHeaderData.UncleHash.String(),
		Coinbase:    blockHeaderData.Coinbase.String(),
		Root:        blockHeaderData.Root.String(),
		TxHash:      blockHeaderData.TxHash.String(),
		ReceiptHash: blockHeaderData.ReceiptHash.String(),
		Difficulty:  blockHeaderData.Difficulty.BitLen(),
		GasLimit:    blockHeaderData.GasLimit,
		GasUsed:     blockHeaderData.GasUsed,
		Time:        blockHeaderData.Time,
		MixDigest:   blockHeaderData.MixDigest.String(),
		Nonce:       blockHeaderData.Nonce.Uint64(),
		BaseFee:     blockHeaderData.BaseFee.BitLen(),
	}, result)
	if err != nil {
		return err
	}
	return nil
}

func (s *blockService) GetTxInfoFromBlockChain(result *models.BlockDataInfoRes) error {
	ethereumConn, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error("Failed to establish the connection :" + err.Error())
		return err
	}
	defer ethereumConn.Close()
	blockInfo, err := ethereumConn.BlockByNumber(context.Background(), big.NewInt(int64(result.BlockData.BlockNum)))
	if err != nil {
		log.Logger.Error("Failed to query block data information from blockChain :" + err.Error())
		return err
	}
	chainId, err := ethereumConn.ChainID(context.Background())
	if err != nil {
		log.Logger.Error("Failed to get chainId :" + err.Error())
	}

	for _, tx := range blockInfo.Transactions() {
		var txDataInfoRes models.TxDataInfoRes
		err := models.NewTxData().TxDataInfo(chainId.Int64(), tx.Hash(), &txDataInfoRes)
		if err != nil {
			return err
		}
		if txDataInfoRes.TxData.TxHash != "" {
			result.TxDataInfoRes = append(result.TxDataInfoRes, txDataInfoRes)
		} else {
			byTx, err := NewTx().GetTxDataByTx(tx, chainId.Int64())
			if err != nil {
				return err
			}
			err = models.NewTxData().CreateTxDataInfo(byTx, &txDataInfoRes)
			if err != nil {
				return err
			}
			result.TxDataInfoRes = append(result.TxDataInfoRes, txDataInfoRes)
		}
	}
	return nil
}
