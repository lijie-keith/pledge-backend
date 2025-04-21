package request

type BlockDataInfo struct {
	ChainId int64 `form:"chainId" binding:"required"`
}
