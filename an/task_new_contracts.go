package an

import (
	"sort"
	"time"

	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/gomisc/logger"
)

func (c *An) taskNewContracts(result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::taskNewContracts begin")

	m := make(map[string]float64)
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		if tx.TxNewContract == "" {
			continue
		}
		if !tx.TxValid {
			continue
		}
		m[tx.TxNewContract] = float64(tx.TxGasUsed)
	}

	type Item struct {
		ContractAddress string
		GasUsed         float64
	}

	items := make([]*Item, 0)
	for key, value := range m {
		var item Item
		item.ContractAddress = key
		item.GasUsed = value
		items = append(items, &item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].GasUsed > items[j].GasUsed
	})

	for i := 0; i < len(items); i++ {
		item := items[i]
		var tableItem ResultItemTable
		tableItem.Text = item.ContractAddress
		tableItem.Values = append(tableItem.Values, item.GasUsed)
		result.ItemsTable = append(result.ItemsTable, &tableItem)
	}

	result.Count = len(result.ItemsTable)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::taskNewContracts end")
}
