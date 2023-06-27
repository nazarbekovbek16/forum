package config

import (
	"bytes"
	"encoding/json"
	"os"
)

type Config struct {
	Http struct {
		Addr          string `json:"port"`
		ReadTimeout   int    `json:"readTimeout"`
		WriteTimeout  int    `json:"writeTimeout"`
		MaxHeaderByte int    `json:"oneMB"`
	}

	Db struct {
		Driver     string `json:"driver"`
		DBName     string `json:"dbname"`
		CtxTimeout int    `json:"timeout"`
	}
}

func NewConfig(cfgFilePath string) *Config {
	var cfg *Config
	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return nil
	}
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&cfg); err != nil {
		return nil
	}
	return cfg
}
