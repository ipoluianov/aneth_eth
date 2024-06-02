package an

import (
	"sort"
	"time"

	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/gomisc/logger"
)

func (c *An) taskAccountsBySendCount(result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::taskAccountsBySendCount begin")

	m := make(map[string]float64)
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		if !tx.TxValid {
			continue
		}
		m[tx.TxFrom] += 1
	}

	type Item struct {
		addr  string
		count float64
	}

	items := make([]*Item, 0)
	for key, value := range m {
		var item Item
		item.addr = key
		item.count = value
		items = append(items, &item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].count > items[j].count
	})

	for i := 0; i < 10; i++ {
		if i >= len(items) {
			break
		}
		item := items[i]
		var tableItem ResultItemTable
		tableItem.Text = item.addr
		tableItem.Values = append(tableItem.Values, item.count)
		result.ItemsTable = append(result.ItemsTable, &tableItem)
	}

	result.Count = len(result.ItemsTable)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::taskAccountsBySendCount end")
}
