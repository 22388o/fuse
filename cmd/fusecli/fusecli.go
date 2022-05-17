package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	var (
		noop       = func(context.Context, []string) error { return flag.ErrHelp }
		cliFlagSet = flag.NewFlagSet("fusecli", flag.ExitOnError)
	)

	balance := &ffcli.Command{
		Name:       "balance",
		ShortUsage: "fusecli balance",
		ShortHelp:  "check balance of wallet",
		LongHelp:   "Check balance of wallet",
		Exec: func(ctx context.Context, args []string) error {
			resp, err := http.Get("http://localhost:1000/balance")
			if err != nil {
				return err
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, string(body))
			return nil
		},
	}

	fuse := &ffcli.Command{
		ShortUsage:  "fusecli [global flags] <subcommand> [subcommand flags] [subcommand args]",
		ShortHelp:   "client for the Fuse Wallet API",
		LongHelp:    "",
		FlagSet:     cliFlagSet,
		Subcommands: []*ffcli.Command{balance},
		Exec:        noop,
	}

	if err := fuse.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		// no need to output redundant message if printing help
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "mapi: %v\n", err)
		}
		os.Exit(1)
	}
}
