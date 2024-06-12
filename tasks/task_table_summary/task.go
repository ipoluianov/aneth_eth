package task_table_summary

import (
	"fmt"
	"strconv"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
)

func New() *common.Task {
	var c common.Task
	c.Code = "summary"
	c.Name = "Summary"
	c.Type = "table"
	c.Fn = Run
	c.Description = "Ethereum - Summary"
	c.Text = ""
	c.Ticker = ""
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	type Item struct {
		Name  string
		Value float64
	}

	itemTotalCount := 0.0
	itemValidCount := 0.0
	itemInvalidCount := 0.0
	itemNewContractCount := 0.0
	itemNativeTransfers := 0.0
	itemContractCalls := 0.0
	itemTotalTransferAmount := 0.0
	itemTotalContractCallAmount := 0.0

	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		itemTotalCount += 1
		if !tx.TxValid {
			itemInvalidCount += 1
			continue
		}
		itemValidCount += 1

		value, _ := tx.TxValue.Float64()

		if tx.TxNewContract != "" {
			itemNewContractCount++
		}

		if len(tx.TxDataMethod) == 0 {
			itemNativeTransfers++
			itemTotalTransferAmount += value
		}

		if len(tx.TxDataMethod) != 0 && tx.TxNewContract == "" {
			itemContractCalls++
			itemTotalContractCallAmount += value
		}
	}

	itemTotalTransferAmount = itemTotalTransferAmount / 1000000000000000000
	itemTotalContractCallAmount = itemTotalContractCallAmount / 1000000000000000000

	result.Table.Columns = append(result.Table.Columns, &common.ResultTableColumn{Name: "Name"})
	result.Table.Columns = append(result.Table.Columns, &common.ResultTableColumn{Name: "Value"})

	fAdd := func(name string, value string) {
		var tableItem common.ResultTableItem
		tableItem.Values = append(tableItem.Values, name)
		tableItem.Values = append(tableItem.Values, fmt.Sprintln(value))
		result.Table.Items = append(result.Table.Items, &tableItem)
	}

	fAdd("Total Transactions", strconv.FormatFloat(itemTotalCount, 'f', 0, 64))
	fAdd("Valid Transactions", strconv.FormatFloat(itemValidCount, 'f', 0, 64))
	fAdd("Invalid Transactions", strconv.FormatFloat(itemInvalidCount, 'f', 0, 64))
	fAdd("New Contracts", strconv.FormatFloat(itemNewContractCount, 'f', 0, 64))
	fAdd("Regular Transfers Count", strconv.FormatFloat(itemNativeTransfers, 'f', 0, 64))
	fAdd("Regular Transfers Amount", strconv.FormatFloat(itemTotalTransferAmount, 'f', 3, 64))
	fAdd("Contract Calls Count", strconv.FormatFloat(itemContractCalls, 'f', 0, 64))
	fAdd("Contract Calls Amount", strconv.FormatFloat(itemTotalContractCallAmount, 'f', 3, 64))
}
