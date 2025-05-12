package common

//商品
type Product struct {
	Id            int
	Name          string
	Size          int     //产品尺寸
	Sale          int     //销量
	ShipAddress   string  //发货地址
	Price         float64 //单价
	PositiveRatio float64 //好评率
	RatioCount    int     //评论量
}
