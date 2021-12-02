// Package main package main
package main

import (
	"os"

	"github.com/manifoldco/promptui"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/simulateprocess"
)

func simulatationEntry() {
	var err error
	tmpChan := make(chan string)
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		prompt := promptui.Prompt{
			Label: "Simulate?(y/n)",
		}
		var result string
		result, err = prompt.Run()
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		if result == "y" {
			ansArr := simulationPrompt()
			simulateprocess.Simulate(ansArr[0], ansArr[1], ansArr[2], ansArr[3])
			<-tmpChan
		}
	}
	getConds()
}

func getConds() {
	var err error
	tmpChan := make(chan string)
	global.ForwardCond, err = simulate.GetBestForwardCondByTradeDay(global.TradeDay, database.GetAgent())
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	if global.ForwardCond.Model.ID == 0 {
		simulateprocess.Simulate("a", "n", "n", "1")
		global.ForwardCond, err = simulate.GetBestForwardCondByTradeDay(global.TradeDay, database.GetAgent())
		if err != nil {
			logger.GetLogger().Panic(err)
		}
	}
	global.ReverseCond, err = simulate.GetBestReverseCondByTradeDay(global.TradeDay, database.GetAgent())
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	if global.ReverseCond.Model.ID == 0 {
		simulateprocess.Simulate("b", "n", "n", "1")
		global.ReverseCond, err = simulate.GetBestReverseCondByTradeDay(global.TradeDay, database.GetAgent())
		if err != nil {
			logger.GetLogger().Panic(err)
		}
	}
	if global.ForwardCond.Model.ID == 0 || global.ReverseCond.Model.ID == 0 {
		logger.GetLogger().Warn("no cond to trade")
		<-tmpChan
	}
	// forwardResult, err := simulate.GetResultByCond(int(global.ForwardCond.ID), database.GetAgent())
	// if err != nil {
	// 	logger.GetLogger().Panic(err)
	// }
	// reverseResult, err := simulate.GetResultByCond(int(global.ReverseCond.ID), database.GetAgent())
	// if err != nil {
	// 	logger.GetLogger().Panic(err)
	// }
	// if float64(forwardResult.Balance)/float64(reverseResult.Balance) < 0.5 {
	// 	global.TradeSwitch.Buy = false
	// 	logger.GetLogger().Warn("TradeSwitch Buy is OFF")
	// }
	// if float64(forwardResult.Balance)/float64(reverseResult.Balance) > 2 {
	// 	global.TradeSwitch.SellFirst = false
	// 	logger.GetLogger().Warn("TradeSwitch SellFirst is OFF")
	// }
	simulateprocess.ClearAllNotBest()
}

func simulationPrompt() []string {
	prompt := promptui.Prompt{
		Label: "Balance type?(a: forward, b: reverse)",
	}
	balanceTypeAns, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	prompt = promptui.Prompt{
		Label: "Discard over time trade?(y/n)",
	}
	discardAns, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	prompt = promptui.Prompt{
		Label: "Use global cond?(y/n)",
	}
	useDefault, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	prompt = promptui.Prompt{
		Label: "N days?",
	}
	countAns, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	return []string{balanceTypeAns, discardAns, useDefault, countAns}
}
