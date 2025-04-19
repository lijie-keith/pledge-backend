package request

type TxDataInfo struct {
	ChainId int64 `form:"chainId" binding:"required"`
}
