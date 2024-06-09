package task_timechart_price

import (
	"strconv"
	"strings"
	"time"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/market"
	"github.com/ipoluianov/gomisc/logger"
)

func New(symbol string, ticker string) *common.Task {
	var c common.Task

	if symbol != "ETH" && symbol != "BTC" {
		c.Code = "token-" + strings.ToLower(symbol) + "-price"
		c.Name = "Token " + symbol + " Price (USDT)"
	} else {
		c.Code = strings.ToLower(symbol) + "-price"
		c.Name = symbol + " Price (USDT)"
	}

	c.Type = "timechart"
	c.Fn = Run
	c.Description = "Price of " + symbol
	c.Text = ""
	c.Ticker = ticker
	c.Symbol = symbol
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::taskPrice begin")

	/*if strings.ToLower(tokenSymbol) == strings.ToLower("TONCOIN") {
		tokenSymbol = "TON"
	}*/

	price := market.Instance.GetData(strings.ToUpper(task.Ticker))

	for i := 0; i < len(price); i++ {
		src := price[i]
		var item common.ResultTimeChartItem
		item.Index = i
		item.DT = uint64(src.StartTime.Unix())
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")
		openPrice, _ := strconv.ParseFloat(src.LowPrice, 64)
		closePrice, _ := strconv.ParseFloat(src.HighPrice, 64)
		item.Value = (openPrice + closePrice) / 2
		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}

	result.Count = len(result.TimeChart.Items)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::taskPrice end")
}
