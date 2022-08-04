package main

import (
	"context"
	"log"
	"os"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/lootek/go-immulogs"
	"github.com/lootek/go-immulogs/pkg/service"
	"github.com/lootek/go-immulogs/pkg/storage"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "immulogs",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "port", Value: 3322, Aliases: []string{"p"}},
			&cli.StringFlag{Name: "address", Value: "localhost", Aliases: []string{"a"}},
			&cli.StringFlag{Name: "username", Value: "immudb", Aliases: []string{"u"}},
			&cli.StringFlag{Name: "password", Value: "immudb", Aliases: []string{"x"}},
			&cli.StringFlag{Name: "database", Value: "defaultdb", Aliases: []string{"d"}},
		},
		Action: func(cliCtx *cli.Context) error {
			immudbOpts := immudb.DefaultOptions().
				WithAddress(cliCtx.String("address")).
				WithPort(cliCtx.Int("port")).
				WithUsername(cliCtx.String("username")).
				WithPassword(cliCtx.String("password")).
				WithDatabase(cliCtx.String("database"))

			storageService := storage.NewImmuDB(immudbOpts)
			ioService := service.NewREST(storageService)
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
