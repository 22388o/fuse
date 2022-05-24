package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mdedys/fuse/fuse"
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

	createInvoiceFlagSet := flag.NewFlagSet("invoices:create", flag.ExitOnError)
	createInvoiceMemo := createInvoiceFlagSet.String("memo", "", "Invoice memo")
	createInvoiceAmount := createInvoiceFlagSet.Int64("amt", 0, "Amount of sats to receive")

	createInvoice := &ffcli.Command{
		Name:       "invoices:create",
		ShortUsage: "fusecli invoices:create",
		ShortHelp:  "create lightning invoice",
		LongHelp:   "Create Lightning invoice",
		FlagSet:    createInvoiceFlagSet,
		Exec: func(ctx context.Context, args []string) error {

			data := fuse.CreateInvoiceRequest{
				Memo:   *createInvoiceMemo,
				Amount: *createInvoiceAmount,
			}

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(data); err != nil {
				return err
			}

			resp, err := http.Post("http://localhost:1000/invoices", "application/json", &buf)
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

	pay := &ffcli.Command{
		Name:       "pay",
		ShortUsage: "fusecli pay <payment request>",
		ShortHelp:  "pay lightning invoice",
		LongHelp:   "Complete a lightning payment",
		Exec: func(ctx context.Context, args []string) error {

			if len(args) != 1 {
				return errors.New("missing payment request")
			}

			data := fuse.PayRequest{
				Request: args[0],
			}

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(data); err != nil {
				return err
			}

			resp, err := http.Post("http://localhost:1000/pay", "application/json", &buf)
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
		Subcommands: []*ffcli.Command{balance, createInvoice, pay},
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
