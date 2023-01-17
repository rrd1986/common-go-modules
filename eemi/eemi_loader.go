package eemi

import (
	"encoding/json"

	"github.com/rrd1986/common-go-modules/utils"
)

type Config struct {
	MessageID      string `json:"messageId"`
	Message        string `json:"message"`
	ResponseAction string `json:"responseAction"`
	Category       string `json:"category"`
	SeverityLevel  string `json:"severity"`
	Status         int
}

func LoadEemiFromFile(source string, filesystem utils.FileSystemType) (map[string]Config, error) {
	byteValue, err := filesystem.ReadFile(source)
	if err != nil {
		return nil, err
	}
	var result map[string]Config
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}
