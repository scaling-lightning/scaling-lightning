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

	valuesFileData.Ports["web"] = port{
		Port:        8000,
		Expose:      false,
		ExposedPort: 80,
		Protocol:    "TCP",
	}

	valuesFileData.Ports["websecure"] = port{
		Port:        8443,
		Expose:      false,
		ExposedPort: 443,
		Protocol:    "TCP",
	}

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
