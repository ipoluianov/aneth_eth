package main

import (
	"github.com/ipoluianov/aneth_eth/app"
	"github.com/ipoluianov/aneth_eth/application"
	"github.com/ipoluianov/gomisc/logger"
)

func main() {
	name := "aneth_eth"
	application.Name = name
	application.ServiceName = name
	application.ServiceDisplayName = name
	application.ServiceDescription = name
	application.ServiceRunFunc = app.RunAsService
	application.ServiceStopFunc = app.StopService
	logger.Init(logger.CurrentExePath() + "/logs")

	if !application.TryService() {
		app.RunDesktop()
	}
}
