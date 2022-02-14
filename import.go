package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/cachestore"
	"github.com/BuxOrg/bux/datastore"
	"github.com/BuxOrg/bux/importer"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/BuxOrg/bux/utils"
	"github.com/libsv/go-bk/bip32"
)

var WhatsOnChainApiKey = os.Getenv("WOC_API_KEY")

func main() {
	var depth int
	var gapLimit int
	var debug bool
	flag.IntVar(&depth, "depth", 20, "Depth of xpub to check")
	flag.IntVar(&gapLimit, "gap-limit", 10, "Gap limit for unused addresses")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Printf("missing xpub argument")
		os.Exit(1)
	}
	ctx := context.Background()
	buxClient, err := initBuxClient(ctx, debug)
	if err != nil {
		log.Printf("failed to init buxClient: %v", err)
		os.Exit(1)

	}

	defer func() {
		_ = buxClient.Close(context.Background())
	}()

	xpub, err := parseXpub(flag.Args()[0])
	if err != nil {
		log.Printf("not a valid xpub")
		os.Exit(1)
	}
	_, err = buxClient.NewXpub(
		ctx, xpub.String(),
	)
	if err != nil {
		log.Printf("failed to store xpub: %v", err)
		os.Exit(1)
	}
	err = importer.ImportXpub(ctx, buxClient, xpub, depth, gapLimit, "m")
	if err != nil {
		log.Printf("failed to import xpub: %v", err)
		os.Exit(1)
	}
}

func parseXpub(xpubStr string) (*bip32.ExtendedKey, error) {
	return utils.ValidateXPub(xpubStr)
}

func initBuxClient(ctx context.Context, debug bool) (bux.ClientInterface, error) {
	var options []bux.ClientOps
	if debug {
		options = append(options, bux.WithDebugging())
	}
	options = append(options, bux.WithITCDisabled())
	options = append(options, bux.WithInputUtxoCheckDisabled())
	options = append(options, bux.WithAutoMigrate(bux.BaseModels...))
	options = append(options, bux.WithRistretto(cachestore.DefaultRistrettoConfig()))
	options = append(options, bux.WithTaskQ(taskmanager.DefaultTaskQConfig("imp_queue"), taskmanager.FactoryMemory))
	options = append(options, bux.WithSQLite(&datastore.SQLiteConfig{
		CommonConfig: datastore.CommonConfig{
			Debug:       false,
			TablePrefix: "xapi",
		},
		DatabasePath: "./import.db", // "" for in memory
		Shared:       true,
	}))

	x, err := bux.NewClient(ctx, options...)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return x, err
}
