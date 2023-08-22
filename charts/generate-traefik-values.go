package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// ports:
//   node1:
//     port: 1001
//     expose: true
//     exposedPort: 1001
//     protocol: TCP

type valuesFile struct {
	Ports map[string]port `yaml:"ports"`
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
	for i := 1001; i <= 1100; i++ {
		valuesFileData.Ports[fmt.Sprintf("node%d", i-1000)] = port{
			Port:        i,
			Expose:      true,
			ExposedPort: i,
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
