package telnet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDial(t *testing.T) {
	s, err := Dial("tcp", "192.168.1.2:23", Config{})
	require.NoError(t, err)
	err = s.Shell()
	require.NoError(t, err)
	s.Wait()
}
