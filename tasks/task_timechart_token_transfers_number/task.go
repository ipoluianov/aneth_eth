package task_timechart_token_transfers_number

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ipoluianov/aneth_eth/cache"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/tokens"
	"github.com/ipoluianov/aneth_eth/utils"
)

func New(symbol string, name string) *common.Task {
	var c common.Task

	if symbol != "ETH" {
		c.Code = "token-" + strings.ToLower(symbol) + "-number-of-transactions"
		c.Name = "Token " + name + " - Number of transactions per minute"
	} else {
		c.Code = strings.ToLower(symbol) + "-number-of-transfers"
		c.Name = name + " - Number of transactions per minute"
	}

	c.Type = "timechart"
	c.Fn = Run
	c.Description = name + " - number of transactions"
	c.Text = ""
	c.Ticker = ""
	c.Symbol = symbol
	return &c
}

func Run(task *common.Task, result *common.Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	var token *tokens.Token

	tokens := tokens.Instance.GetTokens()
	for _, t := range tokens {
		if t.Symbol == task.Symbol {
			token = t
			break
		}
	}

	if token == nil {
		return
	}

	var div, e = big.NewInt(10), big.NewInt(int64(token.Decimals))
	div.Exp(div, e, nil)

	for i := 0; i < len(txsByMin.Items); i++ {
		src := txsByMin.Items[i]
		var item common.ResultTimeChartItem
		item.Index = i
		item.DT = src.DT
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")

		transferMethodId := []byte{0xA9, 0x05, 0x9C, 0xBB}

		cacheId := result.Code + "_" + item.DTStr + "_" + fmt.Sprint(len(src.TXS))

		v := float64(0)

		cacheItem := cache.Instance.Get(cacheId)
		if cacheItem == nil {
			for _, t := range src.TXS {
				if !t.TxValid {
					continue
				}
				if len(t.TxTo) > 0 && t.TxTo == token.Address {
					isTransfer := utils.CompareMethodId(transferMethodId, t.TxDataMethod)
					if isTransfer {
						v += 1
					}
				}

			}
			cache.Instance.Set(cacheId, v)
		} else {
			v = cacheItem.Value
		}

		item.Value = v

		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}
}
