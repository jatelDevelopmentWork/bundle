package authenticate

import (
	"encoding/base64"

	"github.com/ethereum/go-ethereum/ethclient"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func SetBasicAuth(client *ethclient.Client, username, password string) {
	client.Client().SetHeader("Authorization", "Basic "+basicAuth(username, password))
}
