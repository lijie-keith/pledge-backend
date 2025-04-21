package services

import (
	"github.com/ethereum/go-ethereum/rpc"
	"pledge-backend/api/common/redisKey"
	"pledge-backend/api/models"
	"pledge-backend/api/services"
	"pledge-backend/db"
	"pledge-backend/log"
)

type BlockService struct{}

func NewBlockService() *BlockService {
	return &BlockService{}
}

func (blockService *BlockService) UpdateHeadBlock() {
	updateBlockInfoFromBlockChain("head")
}

func (blockService *BlockService) UpdateFinalizedBlock() {
	updateBlockInfoFromBlockChain("finalized")
}

func (blockService *BlockService) UpdateSafeBlock() {
	updateBlockInfoFromBlockChain("safe")
}

func updateBlockInfoFromBlockChain(blockType string) {
	var result = &models.BlockDataInfoRes{}
	switch blockType {
	case "head":
		updateInfo(nil, result)
	case "finalized":
		blockNum := rpc.FinalizedBlockNumber.Int64()
		updateInfo(&blockNum, result)
	case "safe":
		blockNum := rpc.SafeBlockNumber.Int64()
		updateInfo(&blockNum, result)
	}

	err := db.RedisSet(redisKey.BlockKey+blockType, result, 0)
	if err != nil {
		log.Logger.Warn("update block info into redis failed :" + err.Error())
		return
	}
}

func updateInfo(blockNum *int64, result *models.BlockDataInfoRes) {
	err := services.NewBlock().GetHeaderInfoFromBlockChain(blockNum, result)
	if err != nil {
		log.Logger.Error("update block info failed :" + err.Error())
		return
	}
	err = services.NewBlock().GetTxInfoFromBlockChain(result)
	if err != nil {
		log.Logger.Error("update tx info failed :" + err.Error())
		return
	}
}
