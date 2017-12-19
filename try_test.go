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
