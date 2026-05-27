package cmd

import (
	"fmt"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"

	ai "github.com/asafdavid23/eolctl/pkg/ai"
)

func printRiskNarrative(items []ai.RiskItem, logger *log.Logger) {
	if len(items) == 0 {
		logger.Warn("no risk data to summarize — no components were successfully scanned")
		return
	}
	narrative, err := ai.GenerateRiskNarrative(items)
	if err != nil {
		logger.Errorf("failed to generate risk narrative: %v", err)
		return
	}
	fmt.Println("\n--- AI Risk Summary ---")
	fmt.Println(narrative)
}

func printUpgradeSuggestions(items []ai.UpgradeItem, logger *log.Logger) {
	if len(items) == 0 {
		logger.Warn("no upgrade data available — no components were successfully scanned")
		return
	}
	suggestions, err := ai.SuggestUpgradePath(items)
	if err != nil {
		logger.Errorf("failed to generate upgrade suggestions: %v", err)
		return
	}
	fmt.Println("\n--- AI Upgrade Suggestions ---")
	fmt.Println(suggestions)
}

func riskLevelColor(level string) tablewriter.Colors {
	switch level {
	case "CRITICAL":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor}
	case "HIGH":
		return tablewriter.Colors{tablewriter.FgRedColor}
	case "MEDIUM":
		return tablewriter.Colors{tablewriter.FgYellowColor}
	case "LOW":
		return tablewriter.Colors{tablewriter.FgGreenColor}
	default:
		return tablewriter.Colors{}
	}
}

func renderRichRow(table *tablewriter.Table, row []string) {
	colors := make([]tablewriter.Colors, len(row))
	for i := range colors {
		colors[i] = tablewriter.Colors{}
	}
	colors[len(row)-1] = riskLevelColor(row[len(row)-1])
	table.Rich(row, colors)
}
