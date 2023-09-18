package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// ports:
//   endpoint1:
//     port: 1001
//     expose: true
//     exposedPort: 1001
//     protocol: TCP
// logs:
//   general:
//     # -- By default, the logs use a text format (common), but you can
//     # also ask for the json format in the format option
//     # format: json
//     # By default, the level is set to ERROR.
//     # -- Alternative logging levels are DEBUG, PANIC, FATAL, ERROR, WARN, and INFO.
//     level: ERROR
//   access:
//     # -- To enable access logs
//     enabled: false

type valuesFile struct {
	Ports map[string]port `yaml:"ports"`
	Logs  struct {
		General struct {
			Level string `yaml:"level"`
		} `yaml:"general"`
		Access struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"access"`
	} `yaml:"logs"`
}

type port struct {
	Port        int    `yaml:"port"`
	Expose      bool   `yaml:"expose"`
	ExposedPort int    `yaml:"exposedPort"`
	Protocol    string `yaml:"protocol"`
}

func main() {
	valuesFileData := valuesFile{}
	valuesFileData.Ports = make(map[string]port)

	valuesFileData.Logs.General.Level = "ERROR"
	valuesFileData.Logs.Access.Enabled = false

	valuesFileData.Ports["grpc"] = port{
		Port:        28100,
		Expose:      true,
		ExposedPort: 28100,
		Protocol:    "TCP",
	}

	amountOfNodes := 39

	for i := 1; i <= amountOfNodes; i++ {
		valuesFileData.Ports[fmt.Sprintf("endpoint%d", i)] = port{
			Port:        i + 28100,
			Expose:      true,
			ExposedPort: i + 28100,
			Protocol:    "TCP",
		}
	}

	file, err := os.OpenFile("traefik-values.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("error opening/creating file: %v", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)

	err = enc.Encode(valuesFileData)
	if err != nil {
		log.Fatalf("error encoding: %v", err)
	}

}
