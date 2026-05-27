package helpers

import "time"

type RiskLevel string

const (
	RiskCritical RiskLevel = "CRITICAL"
	RiskHigh     RiskLevel = "HIGH"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskLow      RiskLevel = "LOW"
	RiskUnknown  RiskLevel = "UNKNOWN"
)

type RiskInfo struct {
	Level        RiskLevel
	DaysUntilEOL int
	EOLDate      string
}

func daysToRisk(ds int) RiskLevel {
	switch {
	case ds <= 0:
		return RiskCritical
	case ds <= 90:
		return RiskHigh
	case ds <= 180:
		return RiskMedium
	default:
		return RiskLow
	}
}

func CalculateRisk(eolValue interface{}) RiskInfo {
	switch v := eolValue.(type) {
	case bool:
		if v {
			return RiskInfo{Level: RiskCritical, DaysUntilEOL: 0, EOLDate: "already EOL"}
		} else {
			return RiskInfo{Level: RiskLow, DaysUntilEOL: -1, EOLDate: "no known EOL date"}
		}
	case string:
		eolDate, err := time.Parse("2006-01-02", v)

		if err != nil {
			return RiskInfo{Level: RiskUnknown, DaysUntilEOL: -1, EOLDate: v}
		}

		days := int(time.Until(eolDate).Hours() / 24)
		return RiskInfo{Level: daysToRisk(days), DaysUntilEOL: days, EOLDate: v}
	}
	return RiskInfo{Level: RiskUnknown, DaysUntilEOL: -1, EOLDate: "unknown"}
}
