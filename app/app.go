package app

import (
	"fmt"

	"github.com/ipoluianov/aneth_eth/an"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/httpserver"
	"github.com/ipoluianov/gomisc/logger"
)

func Start() {
	logger.Println("Start begin")
	TuneFDs()

	/*router := gin.Default()
	router.GET("/state", api.State)
	router.GET("/analytic/:code", api.Analytic)
	router.GET("/blocks", api.Blocks)
	router.GET("/latest_block_number", api.LatestBlockNumber)
	router.GET("/block/:id", api.Block)
	go router.Run(":8201")*/

	httpserver.Instance.Start()
	db.Instance.Start()
	an.Instance.Start()

	logger.Println("Start end")
}

func Stop() {
}

func RunDesktop() {
	logger.Println("Running as console application")
	Start()
	fmt.Scanln()
	logger.Println("Console application exit")
}

func RunAsService() error {
	Start()
	return nil
}

func StopService() {
	Stop()
}
