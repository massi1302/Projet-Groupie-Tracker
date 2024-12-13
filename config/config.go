// Description: This file contains the configuration of the application.
package config

import (
	"encoding/json"
	"os"
)

var App *AppConfig

type AppConfig struct {
	Name    string       `json:"app_name"`
	Version string       `json:"version"`
	Server  ServerConfig `json:"server"`
}

type ServerConfig struct {
	URL       string          `json:"url"`
	Port      string          `json:"port"`
	StaticWeb StaticWebConfig `json:"static_web"`
}

type StaticWebConfig struct {
	Template TemplateConfig `json:"template"`
	Assets   AssetsConfig   `json:"assets"`
}

type TemplateConfig struct {
	Dir string `json:"dir"`
}

type AssetsConfig struct {
	Dir string `json:"dir"`
}

func init() {
	var file []byte
	var err error
	if file, err = os.ReadFile("./resources/config.json"); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(file, &App); err != nil {
		panic(err)
	}
}
