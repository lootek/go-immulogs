package service

import (
	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
)

type Storage interface {
	Start() error
	Stop() error

	WriteOne(b bucket.Bucket, e log.Entry) error
	WriteBatch(b bucket.Bucket, e []log.Entry) error

	All(b bucket.Bucket) ([]log.Entry, error)
	Last(b bucket.Bucket, n uint64) ([]log.Entry, error)
	Count(b bucket.Bucket) (uint64, error)
}
