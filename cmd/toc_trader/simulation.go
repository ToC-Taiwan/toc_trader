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
	deployment := os.Getenv("DEPLOYMENT")
	tmpChan := make(chan string)
	var err error
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
		} else {
			logger.GetLogger().Warn("Please run in container mode")
			<-tmpChan
		}
	} else {
		global.ForwardCond, err = simulate.GetBestForwardCondByTradeDay(global.TradeDay, database.GetAgent())
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		if global.ForwardCond.Model.ID == 0 {
			simulateprocess.Simulate("a", "n", "n", "2")
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
			simulateprocess.Simulate("b", "n", "n", "2")
			global.ReverseCond, err = simulate.GetBestReverseCondByTradeDay(global.TradeDay, database.GetAgent())
			if err != nil {
				logger.GetLogger().Panic(err)
			}
		}
		if global.ForwardCond.Model.ID == 0 || global.ReverseCond.Model.ID == 0 {
			logger.GetLogger().Warn("no cond to trade")
			<-tmpChan
		}
		logger.GetLogger().Warnf("BestForward is %+v", global.ForwardCond)
		logger.GetLogger().Warnf("BestReverse is %+v", global.ReverseCond)
		simulateprocess.ClearAllNotBest()
	}
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
