package an

import (
	"fmt"
	"time"

	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/gomisc/logger"
)

func (c *An) taskMinutesNewContracts(result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::taskMinutesNewContracts begin")
	for i := 0; i < len(txsByMin.Items); i++ {
		src := txsByMin.Items[i]
		var item ResultTimeChartItem
		item.Index = i
		item.DT = src.DT
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")

		cacheId := result.Code + "_" + item.DTStr + "_" + fmt.Sprint(len(src.TXS))

		v := float64(0)

		cacheItem := c.cache.Get(cacheId)
		if cacheItem == nil {
			for _, t := range src.TXS {
				if !t.TxValid {
					continue
				}
				if t.TxTo == "" && len(t.TxDataMethod) > 0 {
					v += 1
				}
			}
			c.cache.Set(cacheId, v)
		} else {
			v = cacheItem.Value
		}

		item.Value = v

		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}
	result.Count = len(result.TimeChart.Items)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::taskMinutesNewContracts end")
}
