package task_timechart_volatility

import (
	"math"
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
		c.Code = "token-" + strings.ToLower(symbol) + "-volatility"
		c.Name = "Token " + name + " - volatility"
	} else {
		c.Code = strings.ToLower(symbol) + "-volatility"
		c.Name = name + " - Volatility"
	}

	c.Type = "timechart"
	c.Fn = Run
	c.Description = "Volatility of " + name
	c.Text = ""
	c.Ticker = ticker
	c.Symbol = symbol
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	price := market.Instance.GetData(strings.ToUpper(task.Ticker))
	for i := 0; i < len(price); i++ {
		src := price[i]
		var item common.ResultTimeChartItem
		item.Index = i
		item.DT = uint64(src.StartTime.Unix())
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")
		minPrice, _ := strconv.ParseFloat(src.LowPrice, 64)
		maxPrice, _ := strconv.ParseFloat(src.HighPrice, 64)
		item.Value = math.Abs(maxPrice - minPrice)
		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}
	items := make([]float64, len(result.TimeChart.Items))
	for i := 0; i < len(result.TimeChart.Items); i++ {
		items[i] = result.TimeChart.Items[i].Value
	}
	itemsSmoothed := movingAverage(items, 10)

	m := 0.0
	for i := 0; i < len(itemsSmoothed); i++ {
		m += itemsSmoothed[i]
	}
	m = m / float64(len(itemsSmoothed))

	baseline := 0.0
	baselineCount := 0
	for i := 0; i < len(itemsSmoothed); i++ {
		if itemsSmoothed[i] > 0 && itemsSmoothed[i] < m {
			baseline += itemsSmoothed[i]
			baselineCount++
		}
	}
	baseline = baseline / float64(baselineCount)
	for i := 0; i < len(itemsSmoothed); i++ {
		itemsSmoothed[i] = itemsSmoothed[i] - baseline
		if itemsSmoothed[i] < 0 {
			itemsSmoothed[i] = 0
		}
	}

	for i := 0; i < len(result.TimeChart.Items); i++ {
		result.TimeChart.Items[i].Value = itemsSmoothed[i]
	}
}

func movingAverage(arr []float64, windowSize int) []float64 {
	if windowSize <= 0 {
		return nil
	}

	smoothed := make([]float64, len(arr))

	for i := range arr {
		start := max(0, i-windowSize/2)
		end := min(len(arr), i+windowSize/2+1)
		sum := 0.0
		count := 0

		for j := start; j < end; j++ {
			sum += arr[j]
			count++
		}

		smoothed[i] = sum / float64(count)
	}

	return smoothed
}
