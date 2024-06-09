package task_timechart_rejected

import (
	"fmt"
	"time"

	"github.com/ipoluianov/aneth_eth/cache"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
)

func New() *common.Task {
	var c common.Task
	c.Code = "number-of-rejected-transactions-per-minute"
	c.Name = "Number of rejected transactions per minute"
	c.Type = "timechart"
	c.Fn = Run
	c.Description = "Displays the number of unsuccessful transactions recently. An increase in the number of such transactions indicates possible unsuccessful attacks on the network."
	c.Text = "The graph shows the number of transactions that, after being included in a block, were rejected as a result of executing a smart contract, per minute. Possible reasons for rejection include incorrect call parameters, insufficient funds to complete the operation, errors in the smart contract logic, and failure to meet contract conditions."
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
