package controllers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
)

type TxController struct {
}

// GetTx 获取交易信息
func (c *TxController) GetTx(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.TxDataInfo{}
	var result models.TxDataInfoRes

	txHashTemp := ctx.Param("txHash")
	if txHashTemp == "" {
		res.Response(ctx, statecode.TxHashEmptyErr, nil)
		return
	}
	txHash := common.HexToHash(txHashTemp)

	errCode := validate.NewTxDataInfo().TxDataInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewTx().TxDataInfo(req.ChainId, txHash, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

// GetTxReceipt 获取交易收据信息
func (c *TxController) GetTxReceipt(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.TxReceiptDataInfo{}
	var result models.TxReceiptDataInfoRes

	txHashTemp := ctx.Param("txHash")
	if txHashTemp == "" {
		res.Response(ctx, statecode.TxHashEmptyErr, nil)
		return
	}
	txHash := common.HexToHash(txHashTemp)

	errCode := validate.NewTxReceiptDataInfo().TxReceiptDataInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewTxReceipt().TxReceiptDataInfo(req.ChainId, txHash, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}
