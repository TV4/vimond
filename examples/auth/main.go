package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/TV4/vimond/restapi"
)

func main() {
	args := os.Args

	if len(args) != 6 {
		// 5 additional args are required: url, apiKey, secret, platform, id
		log.Fatalf("missing args: %v", args[1:])
	}

	c := restapi.NewClient(
		restapi.BaseURL(args[1]),
		restapi.Credentials(args[2], args[3]),
	)

	asset, err := c.Asset(context.Background(), args[4], args[5])
	if err != nil {
		log.Fatalf("error fetching asset: %v", err)
	}

	json.NewEncoder(os.Stdout).Encode(*asset)
}
