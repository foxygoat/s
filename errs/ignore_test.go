package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIgnore(t *testing.T) {
	executed := false
	f := func() error {
		executed = true
		return errors.New("💥")
	}
	Ignore(f)
	require.True(t, executed)
}
