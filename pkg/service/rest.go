package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
)

type REST struct {
	srv     *http.Server
	storage Storage

	address string
	timeout time.Duration
}

func NewREST(s Storage, address string, timeout time.Duration) *REST {
	globalRouter := gin.New()
	globalRouter.Use(
		gin.Logger(),
		gin.Recovery(),
		// gin.BasicAuth(),
	)

	// TODO: Is there a better way to make bucket an optional parameter?
	for _, router := range []gin.IRoutes{globalRouter, globalRouter.Group("/:bucket")} {
		router.POST("/add", ginWrapper(func(c *gin.Context) (gin.H, error) {
			var entry log.Entry
			if err := c.Bind(&entry); err != nil {
				return nil, err
			}

			if entry == nil {
				data, err := c.GetRawData()
				if err != nil {
					return nil, err
				}

				entry = log.FromBytes(data)
			}

			b := bucket.NewBucket(c.Param("bucket"))
			return addLog(s, b, entry)
		}))
		router.POST("/batch", ginWrapper(func(c *gin.Context) (gin.H, error) {
			var entries []log.Entry

			// try parsing input by gin framework
			if err := c.Bind(&entries); err != nil {
				return nil, err
			}

			// fallback to manual parsing
			if entries == nil {
				data, err := c.GetRawData()
				if err != nil {
					return nil, err
				}

				// try JSON-encoded list
				var logs []string
				err = json.Unmarshal(data, &logs)
				if err != nil {

					// fallback to splitting the given string by the newline character
					logs = strings.Split(string(data), "\n")
				}

				for _, s := range logs {
					entries = append(entries, log.FromString(s))
				}
			}

			b := bucket.NewBucket(c.Param("bucket"))
			return addLogsBatch(s, b, entries)
		}))
		router.GET("/last/:n", ginWrapper(func(c *gin.Context) (gin.H, error) {
			n, err := strconv.ParseInt(c.Param("n"), 10, 64)
			if err != nil {
				return nil, err
			}

			b := bucket.NewBucket(c.Param("bucket"))
			entries, err := lastN(s, b, n)
			if err != nil {
				return nil, err
			}

			return map[string]any{"entries": entries}, err
		}))
		router.GET("/count", ginWrapper(func(c *gin.Context) (gin.H, error) {
			b := bucket.NewBucket(c.Param("bucket"))
			cnt, err := count(s, b)
			if err != nil {
				return nil, err
			}

			return map[string]any{"count": cnt}, err
		}))
	}

	r := &REST{
		srv: &http.Server{
			Addr:         address,
			Handler:      globalRouter,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	}

	return r
}

func ginWrapper(fn func(c *gin.Context) (gin.H, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		res, err := fn(c)

		if err != nil {
			if res != nil {
				c.JSON(http.StatusBadRequest, res)
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func (r REST) Start(context.Context) error {
	return r.srv.ListenAndServe()
}

func (r REST) Stop() error {
	return r.srv.Close()
}

func addLog(s Storage, b bucket.Bucket, e log.Entry) (map[string]any, error) {
	return s.WriteOne(b, e)
}

func addLogsBatch(s Storage, b bucket.Bucket, e []log.Entry) (map[string]any, error) {
	return s.WriteBatch(b, e)
}

func lastN(s Storage, b bucket.Bucket, n int64) ([]log.Entry, error) {
	if n > 0 {
		return s.Last(b, uint64(n))
	}

	return s.All(b)
}

func count(s Storage, b bucket.Bucket) (uint64, error) {
	return s.Count(b)
}
