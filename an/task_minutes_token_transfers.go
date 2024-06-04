package an

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/tokens"
	"github.com/ipoluianov/aneth_eth/utils"
	"github.com/ipoluianov/gomisc/logger"
)

func (c *An) taskMinutesTokenTransfers(result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::taskMinutesTokenTransfers begin")

	parts := strings.FieldsFunc(result.Code, func(r rune) bool {
		return r == '-'
	})

	if len(parts) == 0 {
		return
	}

	tokenSymbol := parts[0]

	var token *tokens.Token

	tokens := tokens.Instance.GetTokens()
	for _, t := range tokens {
		if t.Symbol == tokenSymbol {
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
		var item ResultTimeChartItem
		item.Index = i
		item.DT = src.DT
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")

		transferMethodId := []byte{0xA9, 0x05, 0x9C, 0xBB}

		cacheId := result.Code + "_" + item.DTStr + "_" + fmt.Sprint(len(src.TXS))

		vFloat := 0.0
		v := big.NewInt(0)

		cacheItem := c.cache.Get(cacheId)
		if cacheItem == nil {
			for _, t := range src.TXS {
				if !t.TxValid {
					continue
				}
				if len(t.TxTo) > 0 && t.TxTo == token.Address {
					isTransfer := utils.CompareMethodId(transferMethodId, t.TxDataMethod)
					if isTransfer {
						valueBigInt := utils.ParseDatabigInt(t.TxDataP2)
						//valueBigInt = valueBigInt.Div(valueBigInt, div)
						v = v.Add(v, valueBigInt)
					}
				}

			}
			v = v.Div(v, div)
			vFloat, _ := v.Float64()
			c.cache.Set(cacheId, vFloat)
		} else {
			vFloat = cacheItem.Value
		}

		item.Value = vFloat

		result.TimeChart.Items = append(result.TimeChart.Items, &item)
	}
	result.Count = len(result.TimeChart.Items)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::taskMinutesTokenTransfers end")
}
