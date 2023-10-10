package tools

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

func RetryWithReturn[T any](
	operation func(cancelFn context.CancelFunc) (T, error),
	delay time.Duration,
	timeout time.Duration,
) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	var returnVal T
	for {
		select {
		case <-ctx.Done():
			return returnVal, errors.Wrap(ctx.Err(), "Retry ending")
		case <-ticker.C:
			var err error
			returnVal, err = operation(cancel)
			if err != nil {
				log.Trace().Err(err).Msg("Error was")
				log.Debug().Msg("Retrying...")
				// retry
				continue
			}
		}
		break
	}
	return returnVal, nil
}

func Retry(
	operation func(cancelFn context.CancelFunc) error,
	delay time.Duration,
	timeout time.Duration,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "Retry ending")
		case <-ticker.C:
			var err error
			err = operation(cancel)
			if err != nil {
				log.Trace().Err(err).Msg("Error was")
				log.Debug().Msg("Retrying...")
				// retry
				continue
			}
		}
		break
	}
	return nil
}
