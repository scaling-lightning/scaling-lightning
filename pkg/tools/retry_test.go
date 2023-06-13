package tools

import (
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
	Retry(func() error {
		retryCount++
		return errors.New("Rando error")
	}, 1*time.Millisecond, 10*time.Millisecond)

	assert.Greater(retryCount, 2, "Function should have run at least a few times")
}
