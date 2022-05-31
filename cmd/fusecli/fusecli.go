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
	"os/exec"
	"strconv"

	"github.com/mdedys/fuse/fuse"
	"github.com/mdedys/fuse/lightning"
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
			resp, err := http.Get("http://localhost:1100/balance")
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

			resp, err := http.Post("http://localhost:1100/invoices", "application/json", &buf)
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

			resp, err := http.Post("http://localhost:1100/pay", "application/json", &buf)
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

	var (
		openChannelFlagSet = flag.NewFlagSet("channels:open", flag.ExitOnError)
		openChannelAddr    = openChannelFlagSet.String("node", "", "<pubkey>@<host>")
	)

	openChannel := &ffcli.Command{
		Name:       "open",
		ShortUsage: "fusecli channels open <localSat> <pushSat>",
		ShortHelp:  "commands to manage channels",
		LongHelp:   "Manage channels on the lightning network",
		FlagSet:    openChannelFlagSet,
		Exec: func(ctx context.Context, args []string) error {

			if len(args) != 2 {
				return errors.New("channels open requires localSat and pushSat")
			}

			local, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			push, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			data := fuse.OpenChannelRequest{
				Addr:        lightning.LightningAddress(*openChannelAddr),
				LocalAmount: local,
				PushAmount:  push,
			}

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(data); err != nil {
				return err
			}

			resp, err := http.Post("http://localhost:1100/channels", "application/json", &buf)
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

	channels := &ffcli.Command{
		Name:        "channels",
		ShortUsage:  "fusecli channels <sub-command>",
		ShortHelp:   "commands to manage channels",
		LongHelp:    "Manage channels on the lightning network",
		Subcommands: []*ffcli.Command{openChannel},
		Exec:        noop,
	}

	lnNewAddress := &ffcli.Command{
		Name:       "newaddress",
		ShortUsage: "fusecli ln newaddress",
		ShortHelp:  "run commands against lightning node",
		LongHelp:   "CLI interface for interacting with lightning node",
		Exec: func(ctx context.Context, args []string) error {
			str := "docker exec -t fuse_lnd /bin/bash -c \"lncli --tlscertpath /lnd/tls.cert --macaroonpath /lnd/data/chain/bitcoin/regtest/admin.macaroon newaddress np2wkh\""
			cmd := exec.Command("/bin/bash", "-c", str)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Fprint(os.Stderr, string(output))
				return err
			}

			fmt.Fprint(os.Stdout, string(output))
			return nil
		},
	}

	lightning := &ffcli.Command{
		Name:        "ln",
		ShortUsage:  "fusecli ln <sub-command>",
		ShortHelp:   "run commands against lightning node",
		LongHelp:    "CLI interface for interacting with lightning node",
		Subcommands: []*ffcli.Command{lnNewAddress},
		Exec:        noop,
	}

	var (
		btcdFlagSet = flag.NewFlagSet("btcd", flag.ExitOnError)
		btcNode     = btcdFlagSet.String("node", "bitcoind", "Docker container running bitcoind")

		mineFlagSet = flag.NewFlagSet("mine", flag.ExitOnError)
		blocks      = mineFlagSet.Int("blocks", 10, "Number of blocks to mine")
	)

	mine := &ffcli.Command{
		Name:       "mine",
		ShortUsage: "fusecli btcd mine <address>",
		ShortHelp:  "mine some blocks",
		LongHelp:   "Mine blocks immediately to a specified address",
		FlagSet:    mineFlagSet,
		Exec: func(ctx context.Context, args []string) error {

			if len(args) != 1 {
				return errors.New("mine takes 1 arg")
			}

			str := fmt.Sprintf("docker exec -t %s /bin/bash -c \"bitcoin-cli -chain=regtest -rpcuser=regtest -rpcpassword=regtest -rpcwait\" generatetoaddress %v %v", *btcNode, *blocks, args[0])
			cmd := exec.Command("/bin/sh", "-c", str)

			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Fprint(os.Stderr, string(output))
				return err
			}

			fmt.Fprint(os.Stdout, string(output))
			return nil
		},
	}

	btcd := &ffcli.Command{
		Name:        "btcd",
		ShortUsage:  "fusecli btcd <sub-command>",
		ShortHelp:   "run commands against bitcoind",
		LongHelp:    "CLI interface for interacting with bitcoind",
		FlagSet:     btcdFlagSet,
		Subcommands: []*ffcli.Command{mine},
		Exec:        noop,
	}

	fuse := &ffcli.Command{
		ShortUsage:  "fusecli [global flags] <subcommand> [subcommand flags] [subcommand args]",
		ShortHelp:   "client for the Fuse Wallet API",
		LongHelp:    "",
		FlagSet:     cliFlagSet,
		Subcommands: []*ffcli.Command{balance, channels, createInvoice, pay, btcd, lightning},
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
