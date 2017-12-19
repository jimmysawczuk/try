package try

import (
	"log"
	"time"

	"github.com/pkg/errors"
)

var (
	// DefaultTimeout  = 15 * time.Second
	// DefaultInterval = 100 * time.Millisecond

	errTimedOut = errors.New("timed out")
)

// Try runs the provided function until either the function runs succcessfully (returns a nil error)
// or the timeout duration is exceeded. If the function doesn't run successfully, and there is still
// time remaining before the timeout, it'll sleep for the interval provided before trying to run the
// function again.
func Try(f func() error, timeout, interval time.Duration) error {
	// This channel will fire after the timeout duration has elapsed.
	timeoutCh := time.After(timeout)

	// We'll fire on this channel when the function returns a nil error.
	finishCh := make(chan bool)

	// Run the function in a goroutine; if the function is successful, write to the channel and terminate the go routine, otherwise, log
	// the error, sleep for an interval and try again.
	go func() {
		for {
			start := time.Now()
			err := f()
			if err == nil {
				finishCh <- true
				return
			}

			log.Println(errors.Wrapf(err, "try (attempt took %s)", time.Now().Sub(start).Truncate(time.Millisecond)))
			time.Sleep(interval)
		}
	}()

	// Whichever channel fires first determines whether or not the function ran successfully within the timeout duration.
	select {
	case <-timeoutCh:
		return errTimedOut

	case <-finishCh:
		return nil
	}
}
