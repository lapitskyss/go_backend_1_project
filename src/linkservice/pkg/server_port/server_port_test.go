package server_port

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultServerPort(t *testing.T) {
	serverPort := GetServerPortFromEnv("TEST_PORT", 3000)

	require.Equal(t, ":3000", serverPort)
}

func TestEnvServerPort(t *testing.T) {
	err := os.Setenv("TEST_PORT", "3001")
	require.NoError(t, err, "could not create env variable TEST_PORT")
	serverPort := GetServerPortFromEnv("TEST_PORT", 3000)

	require.Equal(t, ":3001", serverPort)
}
