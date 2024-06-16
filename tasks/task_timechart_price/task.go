package task_timechart_price

import (
	"strconv"
	"strings"
	"time"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/market"
)

func New(symbol string, name string, ticker string) *common.Task {
	var c common.Task

	if symbol != "ETH" && symbol != "BTC" {
		c.Code = "token-" + strings.ToLower(symbol) + "-price"
		c.Name = "Token " + name + " - Price (USDT)"
	} else {
		c.Code = strings.ToLower(symbol) + "-price"
		c.Name = name + " - Price (USDT)"
	}

	c.Type = "timechart"
	c.Fn = Run
	c.Description = "Price of " + name
	c.Text = ""
	c.Ticker = ticker
	c.Symbol = symbol
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	price := market.Instance.GetData(strings.ToUpper(task.Ticker))
	values := make([]float64, 0)
	for i := 0; i < len(price); i++ {
		src := price[i]
		var item common.ResultTimeChartItem
		item.Index = i
		item.DT = uint64(src.StartTime.Unix())
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")
		openPrice, _ := strconv.ParseFloat(src.LowPrice, 64)
		closePrice, _ := strconv.ParseFloat(src.HighPrice, 64)
		item.Value = (openPrice + closePrice) / 2
		values = append(values, item.Value)
		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}

}
