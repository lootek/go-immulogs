package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	lorem "github.com/drhodes/golorem"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "logs-gen",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "rest-address", Value: "http://localhost:8000/generator/batch"},
			&cli.Int64Flag{Name: "rest-timeout", Value: int64(10 * time.Second)},

			&cli.Int64Flag{Name: "log-interval", Value: int64(10000 * time.Millisecond)},
			&cli.Int64Flag{Name: "log-count", Value: int64(5)},
		},
		Action: func(cliCtx *cli.Context) error {
			ticker := time.NewTicker(time.Duration(cliCtx.Int64("log-interval")))

			c := &http.Client{
				Timeout: time.Duration(cliCtx.Int64("rest-timeout")),
			}

			for {
				select {
				case <-cliCtx.Done():
					return nil
				case <-ticker.C:
					entries := generateLogEntries(cliCtx.Int64("log-count"))
					err := sendLogEntry(c, cliCtx.String("rest-address"), entries)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func sendLogEntry(cli *http.Client, uri string, entries []string) error {
	jsonData, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))

	return nil
}

func generateLogEntries(count int64) []string {
	var entries []string

	for i := int64(0); i < count; i++ {
		entries = append(entries, lorem.Sentence(5, 10))
	}

	return entries
}
