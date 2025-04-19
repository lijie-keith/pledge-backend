package request

type TxReceiptDataInfo struct {
	ChainId int64 `form:"chainId" binding:"required"`
}
