package try

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

var (
	// DefaultTimeout  = 15 * time.Second
	// DefaultInterval = 100 * time.Millisecond

	errTimedOut = errors.New("timed out")
)

// terminableErr represents an error that indicates that a try function shouldn't be
// attempted again.
type terminableErr struct {
	error
}

func (t terminableErr) Error() string {
	return fmt.Sprintf("terminable: %s", t.error.Error())
}

// Try runs the provided function until either the function runs succcessfully (returns a nil error)
// or the timeout duration is exceeded. If the function doesn't run successfully, and there is still
// time remaining before the timeout, it'll sleep for the interval provided before trying to run the
// function again.
func Try(f func() error, timeout, interval time.Duration) error {
	// This channel will fire after the timeout duration has elapsed.
	timeoutCh := time.After(timeout)

	// We'll fire on this channel when the function terminates by either returning a nil error
	// or a terminable error.
	finishCh := make(chan error)

	// Run the function in a goroutine; if the function is successful, write to the channel and terminate the go routine, otherwise, log
	// the error, sleep for an interval and try again.
	go func() {
		for {
			start := time.Now()
			err := f()
			if err == nil {
				finishCh <- nil
				return
			}

			switch err.(type) {
			case terminableErr:
				finishCh <- err
			}

			log.Println(errors.Wrapf(err, "try (attempt took %s)", time.Now().Sub(start).Truncate(time.Millisecond)))
			time.Sleep(interval)
		}
	}()

	// Whichever channel fires first determines whether or not the function ran successfully within the timeout duration.
	select {
	case <-timeoutCh:
		return errTimedOut

	case err := <-finishCh:
		return err
	}
}

// TerminableError wraps a new TerminableError around the provided error.
func TerminableError(e error) error {
	return terminableErr{e}
}
