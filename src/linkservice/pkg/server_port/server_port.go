package server_port

import (
	"os"
	"strconv"
)

func GetServerPortFromEnv(envVariable string, defaultPort int) string {
	if 1 > defaultPort || defaultPort > 65535 {
		panic("incorrect port number, port must be in [1-65535] range")
	}
	if port, ok := os.LookupEnv(envVariable); ok {
		portNumber, err := strconv.Atoi(port)
		if err != nil {
			return ":" + strconv.Itoa(defaultPort)
		}

		if 1 > portNumber || portNumber > 65535 {
			return ":" + strconv.Itoa(defaultPort)
		}

		return ":" + port
	}

	return ":" + strconv.Itoa(defaultPort)
}
