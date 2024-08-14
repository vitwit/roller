package tx

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	cometclient "github.com/cometbft/cometbft/rpc/client/http"
	comettypes "github.com/cometbft/cometbft/types"
	"github.com/pterm/pterm"
)

func MonitorTransaction(wsURL, txHash string) error {
	// Create a new client
	client, err := cometclient.New(wsURL, "/websocket")
	if err != nil {
		return fmt.Errorf("error creating client: %v", err)
	}

	// Start the client
	err = client.Start()
	if err != nil {
		return fmt.Errorf("error starting client: %v", err)
	}

	// nolint errcheck
	defer client.Stop()

	// Convert txHash string to bytes
	txBytes, err := hex.DecodeString(txHash)
	if err != nil {
		return fmt.Errorf("error decoding txHash: %v", err)
	}

	// Create a query to filter transactions
	query := fmt.Sprintf("tm.event='Tx' AND tx.hash='%X'", txBytes)

	// Subscribe to the query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// nolint errcheck
	defer cancel()

	subscription, err := client.Subscribe(ctx, "tx-monitor", query, 100)
	if err != nil {
		return fmt.Errorf("error subscribing: %v", err)
	}
	// nolint errcheck
	defer client.Unsubscribe(ctx, "tx-monitor", query)

	fmt.Println("Monitoring transaction:", txHash)

	spinner, _ := pterm.DefaultSpinner.WithText(
		fmt.Sprintf(
			"waiting for tx with hash %s to finalize",
			pterm.FgYellow.Sprint(txHash),
		),
	).Start()

	// Listen for events
	for {
		select {
		case event := <-subscription:
			txEvent, ok := event.Data.(comettypes.EventDataTx)
			if !ok {
				fmt.Println("Received non-tx event")
				continue
			}

			if txEvent.Result.Code == 0 {
				spinner.Success("transaction succeeded")
				pterm.Info.Printf(
					"Gas wanted: %d, Gas used: %d\n",
					txEvent.Result.GasWanted,
					txEvent.Result.GasUsed,
				)
				return nil
			} else {
				return fmt.Errorf("transaction failed with code %d: %v", txEvent.Result.Code, txEvent.Result.Log)
			}

		case <-time.After(5 * time.Minute):
			return fmt.Errorf("timeout waiting for transaction")

		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		}
	}
}