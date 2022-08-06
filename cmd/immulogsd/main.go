package main

import (
	"context"
	"log"
	"os"
	"time"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/lootek/go-immulogs"
	"github.com/lootek/go-immulogs/pkg/service"
	"github.com/lootek/go-immulogs/pkg/storage"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "immulogsd",
		Flags: []cli.Flag{
			// Storage mode
			&cli.StringFlag{Name: "storage", Value: "memory"}, // memory|immudb

			// Service API mode
			&cli.StringFlag{Name: "api", Value: "rest"}, // rest

			// REST
			&cli.StringFlag{Name: "rest-address", Value: "0.0.0.0:8000"},
			&cli.Int64Flag{Name: "rest-timeout", Value: int64(3 * time.Second)},

			// ImmuDB
			&cli.IntFlag{Name: "immudb-port", Value: 3322},
			&cli.StringFlag{Name: "immudb-host", Value: "localhost"},
			&cli.StringFlag{Name: "immudb-username", Value: "immudb"},
			&cli.StringFlag{Name: "immudb-password", Value: "immudb"},
			&cli.StringFlag{Name: "immudb-database", Value: "defaultdb"},
			&cli.Int64Flag{Name: "immudb-timeout", Value: int64(3 * time.Second)},
		},
		Action: func(cliCtx *cli.Context) error {
			var storageService service.Storage
			switch cliCtx.String("storage") {
			case "immudb":
				immudbOpts := immudb.DefaultOptions().
					WithAddress(cliCtx.String("immudb-host")).
					WithPort(cliCtx.Int("immudb-port")).
					WithUsername(cliCtx.String("immudb-username")).
					WithPassword(cliCtx.String("immudb-password")).
					WithDatabase(cliCtx.String("immudb-database"))

				storageService = storage.NewImmuDB(immudbOpts)
			case "memory":
				storageService = storage.NewMemory()
			}

			var ioService immulogs.Service
			switch cliCtx.String("api") {
			case "rest":
				ioService = service.NewREST(storageService, cliCtx.String("rest-address"), time.Duration(cliCtx.Int64("rest-timeout")))
			}

			srv := immulogs.NewService(storageService, ioService)

			ctx := context.Background()
			if err := srv.Run(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
