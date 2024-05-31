package db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ipoluianov/aneth_eth/utils"
	"github.com/ipoluianov/gomisc/logger"
)

type Block struct {
	Number uint64
	Time   uint64
	Txs    []*Tx
}

func NewBlock() *Block {
	var c Block
	return &c
}

func (c *Block) String() string {
	return fmt.Sprint("Block# ", c.Number, " Time:", c.Time)
}

func (c *Block) Write(fileName string) error {
	bs, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		logger.Println("write block error:", err)
		return err
	}
	bs = utils.PackBytes(bs)
	err = os.WriteFile(fileName, bs, 0666)
	if err != nil {
		logger.Println("write block error:", err)
		return err
	}
	return nil
}

func (c *Block) From(tx *types.Transaction) string {
	from, _ := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	return from.Hex()
}

func (c *Block) Read(fileName string) error {
	bs, err := os.ReadFile(fileName)
	if err != nil {
		logger.Println("Block::Read", "read block error:", err)
		return err
	}
	bs, err = utils.UnpackBytes(bs)
	if err != nil {
		logger.Println("Block::Read", "read block error:", err)
		return err
	}
	err = json.Unmarshal(bs, c)

	return err
}
