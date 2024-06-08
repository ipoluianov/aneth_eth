package an

type ResultTimeChartItem struct {
	Index int
	DT    uint64
	DTStr string
	Value float64
}

type ResultTimeChart struct {
	VAxisName string
	HAxisName string
	Items     []*ResultTimeChartItem
	Tables    []*ResultTable
}
