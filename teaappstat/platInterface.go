package teaappstat

type plateInter interface {
	getAdList() (error, [][]string)
	getPayList() map[string]*userOrderPay
	getPlatName() string
}
