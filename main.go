package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	tmrpc "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/types"
)

var (
	LogLevel      string
	TendermintRpc string

	BlocksDiffInThePast int64 = 100

	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

var rootCmd = &cobra.Command{
	Use:  "wenupgrade",
	Long: "Tool to estimate how much time left till a specific block.",
	Args: cobra.ExactArgs(1),
	Run:  Execute,
}

func Execute(cmd *cobra.Command, args []string) {
	logLevel, err := zerolog.ParseLevel(LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse log level")
	}

	zerolog.SetGlobalLevel(logLevel)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse log level")
	}

	block, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse block")
	}

	latestBlock := getBlock(nil)
	latestHeight := latestBlock.Height
	beforeLatestBlockHeight := latestBlock.Height - BlocksDiffInThePast
	beforeLatestBlock := getBlock(&beforeLatestBlockHeight)

	heightDiff := float64(latestHeight - beforeLatestBlockHeight)
	timeDiff := latestBlock.Time.Sub(beforeLatestBlock.Time).Seconds()

	avgBlockTime := timeDiff / heightDiff

	log.Info().
		Float64("heightDiff", heightDiff).
		Float64("timeDiff", timeDiff).
		Float64("avgBlockTime", avgBlockTime).
		Msg("Average block time")

	blocksToCalculate := block - latestHeight

	log.Info().
		Int64("diff", blocksToCalculate).
		Msg("Blocks till the specified block")

	latestTime := latestBlock.Time
	timeToAddAsSeconds := int64(avgBlockTime * float64(blocksToCalculate))
	timeToAddAsDuration := time.Duration(timeToAddAsSeconds) * time.Second
	calculatedBlockTime := latestTime.Add(timeToAddAsDuration)

	log.Info().
		Time("diff", calculatedBlockTime).
		Msg("Estimated block time")

	log.Info().
		Str("diff", timeToAddAsDuration.String()).
		Msg("Time till block")
}

func getBlock(height *int64) *ctypes.Block {
	client, err := tmrpc.New(TendermintRpc, "/websocket")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create Tendermint client")
	}

	block, err := client.Block(context.Background(), height)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not query Tendermint status")
	}

	return block.Block
}

func main() {
	rootCmd.PersistentFlags().StringVar(&LogLevel, "log-level", "info", "Logging level")
	rootCmd.PersistentFlags().StringVar(&TendermintRpc, "tendermint-rpc", "http://localhost:26657", "Tendermint RPC address")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Could not start application")
	}
}
