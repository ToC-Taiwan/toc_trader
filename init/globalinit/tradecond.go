// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	// Back="38" BuyAt="09:50:43" Date="2021-11-05" Name="長榮航" OriginalBalance="108" SellAt="11:15:41"
	// Back="36" BuyAt="10:04:58" Date="2021-11-05" Name="友達" OriginalBalance="512" SellAt="10:42:51"
	// Back="34" BuyAt="09:14:07" Date="2021-11-05" Name="華航" OriginalBalance="-83" SellAt="10:11:13"
	// Back="51" BuyAt="09:10:28" Date="2021-11-05" Name="永光" OriginalBalance="1726" SellAt="09:30:29"
	// Date="2021-11-05" ForwardBalance="2422" ReverseBalance="0"
	// Back="36" BuyAt="09:26:32" Date="2021-11-04" Name="友達" OriginalBalance="-86" SellAt="09:27:20"
	// Back="26" BuyAt="10:05:50" Date="2021-11-04" Name="金寶" OriginalBalance="339" SellAt="10:15:12"
	// Date="2021-11-04" ForwardBalance="315" ReverseBalance="0"
	// Back="35" BuyAt="09:13:00" Date="2021-11-03" Name="華航" OriginalBalance="266" SellAt="09:59:04"
	// Back="36" BuyAt="09:11:55" Date="2021-11-03" Name="長榮航" OriginalBalance="211" SellAt="09:42:22"
	// Back="34" BuyAt="09:42:54" Date="2021-11-03" Name="友達" OriginalBalance="17" SellAt="09:52:42"
	// Date="2021-11-03" ForwardBalance="599" ReverseBalance="0"
	// Balance="3336" PositiveCount="3" TradeCount="9"
	// HistoryCloseCount="1900" TrimHistoryCloseCount="true"
	// OutInRatio="90" ReverseOutInRatio="0"
	// CloseChangeRatioHigh="3" CloseChangeRatioLow="-1" CloseDiff="0" OpenChangeRatio="-1"
	// ReverseRsiHigh="0.9" ReverseRsiLow="0.1" RsiHigh="0.9" RsiLow="0.1"
	// TicksPeriodCount="2" TicksPeriodLimit="10.4" TicksPeriodThreshold="8" VolumePerSecond="20"
	// Forward Condition
	global.ForwardCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     1900,
		TrimHistoryCloseCount: true,
		OutInRatio:            90,
		ReverseOutInRatio:     0,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   -1,
		CloseDiff:             0,
		OpenChangeRatio:       -1,
		ReverseRsiHigh:        0.9,
		ReverseRsiLow:         0.1,
		RsiHigh:               0.9,
		RsiLow:                0.1,
		TicksPeriodCount:      2,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodThreshold:  8,
		VolumePerSecond:       20,
	}
	// Back="46" BuyLaterAt="09:18:16" Date="2021-11-05" Name="建漢" OriginalBalance="40" SellFirstAt="09:16:21"
	// Back="70" BuyLaterAt="09:35:02" Date="2021-11-05" Name="元晶" OriginalBalance="85" SellFirstAt="09:31:47"
	// Back="34" BuyLaterAt="09:43:56" Date="2021-11-05" Name="華航" OriginalBalance="17" SellFirstAt="09:23:16"
	// Back="70" BuyLaterAt="09:21:02" Date="2021-11-05" Name="東森" OriginalBalance="-164" SellFirstAt="09:19:39"
	// Date="2021-11-05" ForwardBalance="0" ReverseBalance="198"
	// Back="36" BuyLaterAt="09:27:15" Date="2021-11-04" Name="友達" OriginalBalance="-86" SellFirstAt="09:27:14"
	// Back="68" BuyLaterAt="10:06:52" Date="2021-11-04" Name="東森" OriginalBalance="138" SellFirstAt="09:21:08"
	// Date="2021-11-04" ForwardBalance="0" ReverseBalance="156"
	// Back="80" BuyLaterAt="09:30:26" Date="2021-11-03" Name="光磊" OriginalBalance="459" SellFirstAt="09:24:23"
	// Back="34" BuyLaterAt="09:55:00" Date="2021-11-03" Name="友達" OriginalBalance="-83" SellFirstAt="09:54:59"
	// Back="34" BuyLaterAt="09:10:08" Date="2021-11-03" Name="華航" OriginalBalance="-82" SellFirstAt="09:09:59"
	// Back="36" BuyLaterAt="09:19:14" Date="2021-11-03" Name="長榮航" OriginalBalance="-138" SellFirstAt="09:13:27"
	// Back="71" BuyLaterAt="09:26:17" Date="2021-11-03" Name="凌陽" OriginalBalance="580" SellFirstAt="09:19:11"
	// Back="67" BuyLaterAt="09:26:17" Date="2021-11-03" Name="元晶" OriginalBalance="389" SellFirstAt="09:24:10"
	// Date="2021-11-03" ForwardBalance="0" ReverseBalance="1447"
	// Balance="1801" PositiveCount="3" TradeCount="12"
	// HistoryCloseCount="2500" TrimHistoryCloseCount="true"
	// OutInRatio="100" ReverseOutInRatio="6"
	// CloseChangeRatioHigh="3" CloseChangeRatioLow="0" CloseDiff="0" OpenChangeRatio="3"
	// ReverseRsiHigh="0.7" ReverseRsiLow="0.2" RsiHigh="0.7" RsiLow="0.2"
	// TicksPeriodCount="2" TicksPeriodLimit="5.2" TicksPeriodThreshold="4" VolumePerSecond="30"
	// Reverse Condition
	global.ReverseCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     2500,
		TrimHistoryCloseCount: true,
		OutInRatio:            100,
		ReverseOutInRatio:     6,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   0,
		CloseDiff:             0,
		OpenChangeRatio:       3,
		ReverseRsiHigh:        0.7,
		ReverseRsiLow:         0.2,
		RsiHigh:               0.7,
		RsiLow:                0.2,
		TicksPeriodCount:      2,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodThreshold:  4,
		VolumePerSecond:       30,
	}
}
