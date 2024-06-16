package task_timechart_price_minus_btc

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/market"
	"github.com/ipoluianov/gomisc/logger"
)

func New(symbol string, name string, ticker string) *common.Task {
	var c common.Task

	if symbol != "ETH" && symbol != "BTC" {
		c.Code = "token-" + strings.ToLower(symbol) + "-price-minus-btc-price"
		c.Name = "Token " + name + " - Price (USDT) minus BTC price"
	} else {
		c.Code = strings.ToLower(symbol) + "-price-minus-btc-price"
		c.Name = name + " - Price (USDT) minus BTC price"
	}

	c.Type = "timechart"
	c.Fn = Run
	c.Description = "Price of " + name + " minus BTC price"
	c.Text = ""
	c.Ticker = ticker
	c.Symbol = symbol
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("task_timechart_price_minus_btc begin: ", task.Ticker)
	priceToken := market.Instance.GetData(strings.ToUpper(task.Ticker))
	priceBTC := market.Instance.GetData("BTCUSDT")

	count := len(priceToken)

	if len(priceBTC) < len(priceToken) {
		count = len(priceBTC)
	}

	logger.Println("task_timechart_price_minus_btc 1", count)

	pricesToken := make([]float64, count)
	pricesBTC := make([]float64, count)
	// TODO: min max

	btcMin := math.MaxFloat64
	btcMax := float64(0)

	tokenMin := math.MaxFloat64
	tokenMax := float64(0)

	for i := 0; i < count; i++ {
		srcToken := priceToken[i]
		srcBTC := priceBTC[i]

		minPriceToken, _ := strconv.ParseFloat(srcToken.LowPrice, 64)
		maxPriceToken, _ := strconv.ParseFloat(srcToken.HighPrice, 64)
		pricesToken[i] = (minPriceToken + maxPriceToken) / 2
		if pricesToken[i] < tokenMin {
			tokenMin = pricesToken[i]
		}
		if pricesToken[i] > tokenMax {
			tokenMax = pricesToken[i]
		}

		minPriceBTC, _ := strconv.ParseFloat(srcBTC.LowPrice, 64)
		maxPriceBTC, _ := strconv.ParseFloat(srcBTC.HighPrice, 64)
		pricesBTC[i] = (minPriceBTC + maxPriceBTC) / 2
		if pricesBTC[i] < btcMin {
			btcMin = pricesBTC[i]
		}
		if pricesBTC[i] > btcMax {
			btcMax = pricesBTC[i]
		}
	}

	tokenDiff := tokenMax - tokenMin
	for i := 0; i < count; i++ {
		pricesBTC[i] = pricesBTC[i] - btcMin
		btcDiff := (btcMax - btcMin)
		if btcDiff > math.SmallestNonzeroFloat64 {
			pricesBTC[i] = pricesBTC[i] / btcDiff
		}
		pricesBTC[i] = pricesBTC[i]*tokenDiff + tokenMin
		pricesToken[i] = pricesToken[i] - pricesBTC[i]
	}

	for i := 0; i < count; i++ {
		var item common.ResultTimeChartItem
		item.Index = i
		item.DT = uint64(priceToken[i].StartTime.Unix())
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")
		item.Value = pricesToken[i]
		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}

	logger.Println("task_timechart_price_minus_btc end: ", task.Ticker, len(result.TimeChart.Items))
}
