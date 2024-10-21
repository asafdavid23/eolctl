package output

import (
	"eolctl/internal"
)

func Handler(outputData []byte, outputType string, exportToFile bool, exportPath string) error {
	if exportToFile {
		helpers.ExportToFile(outputData, exportPath)
	} else {
		helpers.ConvertOutput(outputData, outputType)
	}

	return nil
}
