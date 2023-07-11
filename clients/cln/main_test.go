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
	assert.Contains(strings.ToLower(err.Error()), "clientcert", "Didn't complain about missing client cert path")

	config = appConfig{clientCertificate: "certpath"}
	err = validateFlags(&config)

	assert.NotNil(err)
	assert.Contains(strings.ToLower(err.Error()), "clientkey", "Didn't complain about missing client key file location")

	config = appConfig{clientCertificate: "certpath", clientKey: "clientkeypath"}
	err = validateFlags(&config)

	assert.NotNil(err)
	assert.Contains(strings.ToLower(err.Error()), "cacert", "Didn't complain about missing ca cert file location")

	config = appConfig{clientCertificate: "certpath", clientKey: "clientkeypath", caCert: "cacertpath"}
	err = validateFlags(&config)

	assert.NotNil(err)
	assert.Contains(strings.ToLower(err.Error()), "grpc", "Didn't complain about missing gRPC address")

}
