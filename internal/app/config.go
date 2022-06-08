package app

import (
	"Coeus/internal/helper"
	"encoding/json"
	"os"
)

type CoeusConfig struct {
	TotalCallNum uint   `json:"totalCallNum"`
	Concurrent   int    `json:"concurrent"`
	TargetHost   string `json:"targetHost"`
	Insecure     bool   `json:"insecure"`
	ProtoFile    string `json:"protoFile"`
	MethodName   string `json:"methodName"`
}

func LoadConfigFromFile(path *string, config *CoeusConfig) error {
	jsonConfig, err := os.ReadFile(*path)
	if err != nil {
		return helper.ErrConfigLoadFailed(path)
	}

	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		return helper.ErrConfigReadFailed
	}
	println(config.Insecure)

	return nil
}
