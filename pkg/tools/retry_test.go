package tools

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	assert := assert.New(t)

	// quieten retry logs
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	retryCount := 0
	err := Retry(func(cancel context.CancelFunc) error {
		retryCount++
		return errors.New("Rando error")
	}, 1*time.Millisecond, 10*time.Millisecond)
	assert.NotNil(err)

	assert.Greater(retryCount, 2, "Function should have run at least a few times")
}

func TestRetryWithReturn(t *testing.T) {
	assert := assert.New(t)

	// quieten retry logs
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	// return value but retries until out of time
	retryCount := 0
	returnVal, err := RetryWithReturn(func(cancel context.CancelFunc) (string, error) {
		retryCount++
		return "", errors.New("Rando error")
	}, 1*time.Millisecond, 10*time.Millisecond)
	assert.NotNil(err)
	assert.Empty(returnVal)
	assert.Greater(retryCount, 2, "Function should have run at least a few times")

	// return value and succeeds first time
	retryCount = 0
	returnVal, err = RetryWithReturn(func(cancel context.CancelFunc) (string, error) {
		retryCount++
		return "SUCCESS", nil
	}, 1*time.Millisecond, 10*time.Millisecond)
	assert.Nil(err)
	assert.Equal("SUCCESS", returnVal)
	assert.Equal(1, retryCount, "Function should have run only once")

	// return value but cancels after first time
	retryCount = 0
	returnVal, err = RetryWithReturn(func(cancel context.CancelFunc) (string, error) {
		retryCount++
		cancel()
		return "", errors.New("Giving up")
	}, 1*time.Millisecond, 10*time.Millisecond)
	assert.NotNil(err)
	assert.Empty(returnVal)
	assert.Equal(1, retryCount, "Function should have run only once")
}
