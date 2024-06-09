package task_table_new_contracts

import (
	"fmt"
	"sort"
	"time"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/gomisc/logger"
)

func New() *common.Task {
	var c common.Task
	c.Code = "new-eth-contracts-list"
	c.Name = "New ETH Contracts - Last 24 hours"
	c.Type = "table"
	c.Fn = Run
	c.Description = "List of new smart contracts"
	c.Text = ""
	c.Ticker = ""
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::taskNewContracts begin")

	type Item struct {
		ContractAddress string
		DT              uint64
		DTStr           string
		GasUsed         uint64
	}

	m := make(map[string]*Item)
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		if tx.TxNewContract == "" {
			continue
		}
		if !tx.TxValid {
			continue
		}
		var item Item
		item.ContractAddress = tx.TxNewContract
		item.DT = tx.BlDT
		item.DTStr = time.Unix(int64(tx.BlDT), 0).Format("2006-01-02 15:04:05")
		item.GasUsed = tx.TxGasUsed
		m[tx.TxNewContract] = &item
	}

	items := make([]*Item, 0)
	for _, value := range m {
		items = append(items, value)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].DT < items[j].DT
	})

	result.Table.Columns = append(result.Table.Columns, &common.ResultTableColumn{Name: "Contract Address"})
	result.Table.Columns = append(result.Table.Columns, &common.ResultTableColumn{Name: "Date/Time"})
	result.Table.Columns = append(result.Table.Columns, &common.ResultTableColumn{Name: "Gas Used"})

	for i := 0; i < len(items); i++ {
		item := items[i]
		var tableItem common.ResultTableItem
		tableItem.Values = append(tableItem.Values, item.ContractAddress)
		tableItem.Values = append(tableItem.Values, item.DTStr)
		tableItem.Values = append(tableItem.Values, fmt.Sprint(item.GasUsed))
		result.Table.Items = append(result.Table.Items, &tableItem)
	}

	result.Count = len(result.Table.Items)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::taskNewContracts end")
}
