package an

import (
	"fmt"
	"time"

	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/utils"
	"github.com/ipoluianov/gomisc/logger"
)

func (c *An) taskMinutesCountOfUsdt(result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx) {
	logger.Println("An::taskMinutesCountOfUsdt begin")
	for i := 0; i < len(txsByMin.Items); i++ {
		src := txsByMin.Items[i]
		var item ResultItemByMinutes
		item.Index = i
		item.DT = src.DT
		item.DTStr = time.Unix(int64(item.DT), 0).UTC().Format("2006-01-02 15:04:05")

		transferMethodId := []byte{0xA9, 0x05, 0x9C, 0xBB}

		cacheId := result.Code + "_" + item.DTStr + "_" + fmt.Sprint(len(src.TXS))

		v := float64(0)

		cacheItem := c.cache.Get(cacheId)
		if cacheItem == nil {
			for _, t := range src.TXS {
				if !t.TxValid {
					continue
				}
				if len(t.TxTo) > 0 && t.TxTo == "0x514910771AF9Ca656af840dff83E8264EcF986CA" {
					isTransfer := utils.CompareMethodId(transferMethodId, t.TxDataMethod)
					if isTransfer {
						v += 1
					}
				}

			}
			c.cache.Set(cacheId, v)
		} else {
			v = cacheItem.Value
		}

		item.Value = v

		result.ItemsByMinutes = append(result.ItemsByMinutes, &item)
	}
	result.Count = len(result.ItemsByMinutes)
	result.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	logger.Println("An::taskMinutesCountOfUsdt end")
}
