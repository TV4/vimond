package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TV4/vimond/restapi"
)

const (
	vimondProd  = "https://restapi-vimond-prod.b17g.net/"
	vimondStage = "https://restapi-vimond-stage.b17g.net/"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr)
		printUsage()
	}

	fAuth := flag.String("auth", "", "API key and secret <key>:<secret>")
	fStage := flag.Bool("stage", false, "Use staging environment instead of prod")

	flag.Parse()

	remainingArgs := flag.Args()

	if len(remainingArgs) < 1 {
		fmt.Fprintf(os.Stderr, "vimond: missing command\n\n")
		printUsage()
		os.Exit(1)
	}

	var opts []func(*restapi.Client)

	if *fStage {
		opts = append(opts, restapi.BaseURL(vimondStage))
	} else {
		opts = append(opts, restapi.BaseURL(vimondProd))
	}

	if *fAuth != "" {
		kv := strings.Split(*fAuth, ":")
		if len(kv) != 2 {
			die("error parsing auth flag")
		}
		opts = append(opts, restapi.Credentials(kv[0], kv[1]))
	}

	client := restapi.NewClient(opts...)

	switch cmd, args := remainingArgs[0], remainingArgs[1:]; cmd {
	case "assets":
		cmdAssets(client, args)
	case "current-orders":
		cmdCurrentOrders(client, args)
	case "orders":
		cmdOrders(client, args)
	case "platforms":
		cmdPlatforms(client)
	case "video-files":
		cmdVideoFiles(client, args)
	default:
		die("unknown command %q", cmd)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: vimond [-auth=<apikey>:<secret>] [-stage] <command> [<args>]")
	fmt.Fprintln(os.Stderr, `
  Commands
    assets <platform> <ids>...           Fetches one or more assets
    current-orders <platform> <user-id>  Fetches current orders for the given user
    orders <platform> <ids>...           Fetches one or more orders
    platforms                            Lists available platforms
    video-files <ids>...                 Fetches video file data for the given asset(s)`)
	fmt.Fprintln(os.Stderr)
}

func die(format string, v ...interface{}) {
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, v...)
	fmt.Fprintln(os.Stderr)
	printUsage()
	os.Exit(1)
}

func cmdAssets(client *restapi.Client, args []string) {
	if len(args) < 2 {
		die("need platform and at least one ID")
	}

	platform := args[0]
	ids := args[1:]

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelCtx()

	for _, id := range ids {
		res, err := client.Asset(ctx, platform, id)
		if err != nil {
			die("error fetching asset (%s): %v", id, err)
		}

		json.NewEncoder(os.Stdout).Encode(res)
	}
}

func cmdCurrentOrders(client *restapi.Client, args []string) {
	if len(args) != 2 {
		die("need platform and user ID")
	}

	platform := args[0]
	userID := args[1]

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelCtx()

	res, err := client.CurrentOrders(ctx, platform, userID)
	if err != nil {
		die("error fetching current orders: %v", err)
	}

	json.NewEncoder(os.Stdout).Encode(res)
}

func cmdOrders(client *restapi.Client, args []string) {
	if len(args) < 2 {
		die("need platform and at least one order ID")
	}

	platform := args[0]
	ids := args[1:]

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelCtx()

	for _, id := range ids {
		res, err := client.Order(ctx, platform, id)
		if err != nil {
			die("error fetching order (%s): %v", id, err)
		}

		json.NewEncoder(os.Stdout).Encode(res)
	}
}

func cmdPlatforms(client *restapi.Client) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelCtx()

	res, err := client.Platforms(ctx)
	if err != nil {
		die("error fetching platforms: %v", err)
	}

	json.NewEncoder(os.Stdout).Encode(res)
}

func cmdVideoFiles(client *restapi.Client, args []string) {
	if len(args) < 1 {
		die("need at least one asset ID")
	}

	ids := args

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelCtx()

	for _, id := range ids {
		res, err := client.Videofiles(ctx, id)
		if err != nil {
			die("error fetching video file data (%s): %v", id, err)
		}

		json.NewEncoder(os.Stdout).Encode(res)
	}
}
