package task_timechart_new_contracts

import (
	"fmt"
	"time"

	"github.com/ipoluianov/aneth_eth/cache"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
)

func New() *common.Task {
	var c common.Task
	c.Code = "number-of-new-contracts-per-minute"
	c.Name = "Number of new contracts per minute"
	c.Type = "timechart"
	c.Fn = Run
	c.Description = "The graph displays the number of transactions that create new smart contracts on the network."
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

		cacheId := result.Code + "_" + item.DTStr + "_" + fmt.Sprint(len(src.TXS))

		v := float64(0)

		cacheItem := cache.Instance.Get(cacheId)
		if cacheItem == nil {
			for _, t := range src.TXS {
				if !t.TxValid {
					continue
				}
				if t.TxTo == "" && len(t.TxDataMethod) > 0 {
					v += 1
				}
			}
			cache.Instance.Set(cacheId, v)
		} else {
			v = cacheItem.Value
		}

		item.Value = v

		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}
}
