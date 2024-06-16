package task_table_summary

import (
	"fmt"
	"strconv"
	"strings"

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
	result.Table.Columns = append(result.Table.Columns, &common.ResultTableColumn{Name: ""})

	result.Table.Columns[1].Align = "right"

	fAdd := func(name string, value string, comment string) {
		var tableItem common.ResultTableItem
		tableItem.Values = append(tableItem.Values, name)
		tableItem.Values = append(tableItem.Values, fmt.Sprint(value))
		tableItem.Values = append(tableItem.Values, comment)
		result.Table.Items = append(result.Table.Items, &tableItem)
	}

	validTransactionPercents := 0.0
	invalidTransactionPercents := 0.0
	if itemTotalCount > 0 {
		validTransactionPercents = 100 * itemValidCount / itemTotalCount
		invalidTransactionPercents = 100 - validTransactionPercents
	}

	fAdd("Total Transactions", formatNumberWithSpaces(strconv.FormatFloat(itemTotalCount, 'f', 0, 64)), "")
	fAdd("Valid Transactions", formatNumberWithSpaces(strconv.FormatFloat(itemValidCount, 'f', 0, 64)), " ("+strconv.FormatFloat(validTransactionPercents, 'f', 0, 64)+"%)")
	fAdd("Invalid Transactions", formatNumberWithSpaces(strconv.FormatFloat(itemInvalidCount, 'f', 0, 64)), " ("+strconv.FormatFloat(invalidTransactionPercents, 'f', 0, 64)+"%)")
	fAdd("New Contracts", formatNumberWithSpaces(strconv.FormatFloat(itemNewContractCount, 'f', 0, 64)), "")

	fAdd("Regular Transfers Count", formatNumberWithSpaces(strconv.FormatFloat(itemNativeTransfers, 'f', 0, 64)), "")
	fAdd("Contract Calls Count", formatNumberWithSpaces(strconv.FormatFloat(itemContractCalls, 'f', 0, 64)), "")

	fAdd("Regular Transfers Amount", formatNumberWithSpaces(strconv.FormatFloat(itemTotalTransferAmount, 'f', 0, 64)), "")
	fAdd("Contract Calls Amount", formatNumberWithSpaces(strconv.FormatFloat(itemTotalContractCallAmount, 'f', 0, 64)), "")
}

func formatNumberWithSpaces(s string) string {
	var result strings.Builder
	length := len(s)

	// Пройдемся по строке в обратном порядке и добавим пробелы каждые три цифры
	for i := 0; i < length; i++ {
		if i > 0 && i%3 == 0 {
			result.WriteString(" ")
		}
		result.WriteByte(s[length-i-1])
	}

	// Перевернем строку, чтобы получить правильный порядок
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
