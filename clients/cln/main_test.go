package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	assert := assert.New(t)

	config := appConfig{}
	err := validateFlags(&config)

	assert.NotNil(err)
	assert.Contains(strings.ToLower(err.Error()), "tls", "Didn't complain about missing TLS file path")

	config = appConfig{tlsFilePath: "tls"}
	err = validateFlags(&config)

	assert.NotNil(err)
	assert.Contains(strings.ToLower(err.Error()), "macaroon", "Didn't complain about missing Macaroon file location")

	config = appConfig{tlsFilePath: "tls", macaroonFilePath: "macaroon"}
	err = validateFlags(&config)

	assert.NotNil(err)
	assert.Contains(strings.ToLower(err.Error()), "grpc", "Didn't complain about missing gRPC address")

}
