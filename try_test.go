package try

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestTry(t *testing.T) {
	f := func(m int) func() error {
		i := 0
		return func() error {
			i++
			if i <= m {
				return errors.New("whoops")
			}
			return nil
		}
	}

	err := Try(f(3), 1*time.Second, 100*time.Millisecond)
	require.Nil(t, err)

	err = Try(f(50), 1*time.Second, 100*time.Millisecond)
	require.Equal(t, errTimedOut, err)
}

func TestFatalTry(t *testing.T) {
	f := func() error {
		return TerminableError(errors.New("wut"))
	}

	err := Try(f, 5*time.Second, 100*time.Millisecond)
	require.Error(t, err)
	require.IsType(t, terminableErr{}, errors.Cause(err))
}

func TestFunctionActuallyStops(t *testing.T) {
	i := 0
	f := func() error {
		i++
		return errors.New("whoops")
	}

	err := Try(f, 1*time.Second, 100*time.Millisecond)
	currentI := i
	require.Error(t, err)
	time.Sleep(1 * time.Second)
	require.Equal(t, currentI, i)
}
