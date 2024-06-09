package task_timechart_count

import (
	"time"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
)

func New() *common.Task {
	var c common.Task
	c.Code = "number-of-transactions-per-minute"
	c.Name = "Number of transactions per minute"
	c.Type = "timechart"
	c.Fn = Run
	c.Description = "Number of transactions by minute on the chart. Only successful transactions are displayed. This data indicates the overall activity of the network."
	c.Text = ""
	c.Ticker = ""
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	for i := 0; i < len(txsByMin.Items); i++ {
		src := txsByMin.Items[i]
		var item common.ResultTimeChartItem
		item.Index = i
		item.DT = src.DT
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")
		item.Value = float64(len(src.TXS))
		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}
	result.Count = len(result.TimeChart.Items)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
}
