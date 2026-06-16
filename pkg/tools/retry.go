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

	var err error
	var returnVal T

	for {
		select {
		case <-ctx.Done():
			return returnVal, errors.Wrap(errors.CombineErrors(errors.Wrap(err, "last error"), ctx.Err()),
				"Retry ending")
		default:
			returnVal, err = operation(cancel)
			if err != nil {
				log.Debug().Err(err).Msg("Error was")
				log.Debug().Msg("Retrying...")
				// wait for delay
				time.Sleep(delay)
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

	var err error

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(errors.CombineErrors(errors.Wrap(err, "last error"), ctx.Err()),
				"Retry ending")
		default:
			err = operation(cancel)
			if err != nil {
				log.Debug().Err(err).Msg("Error was")
				log.Debug().Msg("Retrying...")
				// wait for delay
				time.Sleep(delay)
				// retry
				continue
			}
		}
		break
	}
	return nil
}
