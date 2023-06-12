package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	assert := assert.New(t)

	config := appConfig{}
	err := validateFlags(&config)

	assert.NotNil(err, "Didn't complain about missing cooke file location")

	config = appConfig{rpcCookieFile: "dummy"}
	err = validateFlags(&config)

	assert.NotNil(err, "Didn't complain about missing rpc host")
}
