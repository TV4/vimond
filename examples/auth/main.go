package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
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

	resp, err := c.GetJSON(context.Background(), args[4], args[5], url.Values{"expand": []string{"metadata", "category"}})
	if err != nil {
		log.Fatalf("error fetching asset: %v", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
