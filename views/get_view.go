package views

import (
	"fmt"
	"strings"

	"github.com/ipoluianov/aneth_eth/an"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/static"
)

func GetView(code string, defaultTitle string, defaultDescription string, instance string, chartHeight int, showTitle, showDesc bool, showText bool, showHorScale bool) (result string, title string, description string) {
	title = defaultTitle
	description = defaultDescription
	task := an.Instance.GetTask(code)
	if task == nil {
		return
	}

	if task.Type == "timechart" {
		result = static.FileViewChart
	}

	if task.Type == "table" {
		result = static.FileViewTable
	}

	title = task.Name + " - " + common.GlobalSiteName
	description = task.Description + " " + defaultDescription

	displayDescription := task.Description
	displayText := task.Text
	displayName := task.Name

	displayStyleName := "none"
	if showTitle {
		displayStyleName = "block"
	}

	displayStyleDesc := "none"
	if showTitle {
		displayStyleDesc = "block"
	}

	displayStyleText := "none"
	if showTitle {
		displayStyleText = "block"
	}

	result = strings.ReplaceAll(result, "%VIEW_CODE%", task.Code)
	result = strings.ReplaceAll(result, "%VIEW_NAME%", displayName)
	result = strings.ReplaceAll(result, "%VIEW_DESC%", displayDescription)
	result = strings.ReplaceAll(result, "%VIEW_TEXT%", displayText)
	result = strings.ReplaceAll(result, "VIEW_INSTANCE", instance)
	result = strings.ReplaceAll(result, "VIEW_DISPLAY_NAME", displayStyleName)
	result = strings.ReplaceAll(result, "VIEW_DISPLAY_DESC", displayStyleDesc)
	result = strings.ReplaceAll(result, "VIEW_DISPLAY_TEXT", displayStyleText)
	result = strings.ReplaceAll(result, "VIEW_CHART_HEIGHT", fmt.Sprint(chartHeight))
	result = strings.ReplaceAll(result, "VIEW_DRAW_HOR_SCALE", fmt.Sprint(showHorScale))

	return
}
