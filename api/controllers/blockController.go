package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
	"strconv"
)

type BlockController struct {
}

func (c *BlockController) GetBlockInfo(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.BlockDataInfo{}
	var result models.BlockDataInfoRes

	blockNum := ctx.Param("block_num")
	if blockNum == "" {
		res.Response(ctx, statecode.ParameterEmptyErr, nil)
		return
	}

	full, _ := strconv.ParseBool(ctx.Query("full"))

	errCode := validate.NewBlockDataInfo().BlockDataInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewBlock().BlockDataInfo(req.ChainId, blockNum, full, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}
