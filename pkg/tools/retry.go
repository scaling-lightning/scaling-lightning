package tools

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

type NoRetryError struct{}

func (nre NoRetryError) Error() string {
	return "Fundamental problem, skipping retries"
}

func Retry(operation func() error, delay time.Duration, maxWait time.Duration) error {
	var totalWaited time.Duration
	for {
		if totalWaited > maxWait {
			return errors.New("Exceeded maximum wait period")
		}
		err := operation()
		var noRetry NoRetryError
		if errors.As(err, &noRetry) {
			return err
		}
		if err == nil {
			break
		}
		log.Trace().Err(err).Msg("Error was")
		log.Info().Msg("Retry...")
		time.Sleep(delay)
		totalWaited += delay
	}
	return nil
}
