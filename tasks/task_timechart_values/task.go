package task_timechart_values

import (
	"fmt"
	"time"

	"github.com/ipoluianov/aneth_eth/cache"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/gomisc/logger"
)

func New() *common.Task {
	var c common.Task
	c.Code = "eth-transfer-volume-per-minute"
	c.Name = "ETH transfer volume per minute"
	c.Type = "timechart"
	c.Fn = Run
	c.Description = "The graph shows the total volume of ETH transfers. These can be either regular transfers between accounts or transfers to smart merchant addresses."
	c.Text = ""
	c.Ticker = ""
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::anTrValue begin")
	for i := 0; i < len(txsByMin.Items); i++ {
		src := txsByMin.Items[i]
		var item common.ResultTimeChartItem
		item.Index = i
		item.DT = src.DT
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")

		cacheId := result.Code + "_" + item.DTStr + "_" + fmt.Sprint(len(src.TXS))

		v := float64(0)
		cacheItem := cache.Instance.Get(cacheId)
		if cacheItem == nil {
			for _, t := range src.TXS {
				if !t.TxValid {
					continue
				}
				tv, _ := t.TxValue.Float64()
				v += tv
			}
			cache.Instance.Set(cacheId, v)
		} else {
			v = cacheItem.Value
		}
		item.Value = v / 1000000000000000000

		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}
	result.Count = len(result.TimeChart.Items)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::anTrValue end")
}
