package storage

import (
	"errors"
	"sync"

	"github.com/lootek/go-immulogs/pkg/storage/bucket"
	"github.com/lootek/go-immulogs/pkg/storage/log"
)

type Memory struct {
	dataMu sync.RWMutex
	data   map[bucket.Bucket][]log.Entry
}

func NewMemory() *Memory {
	return &Memory{data: map[bucket.Bucket][]log.Entry{}}
}

func (m *Memory) WriteOne(b bucket.Bucket, e log.Entry) error {
	m.dataMu.RLock()
	defer m.dataMu.RUnlock()

	entries := m.data[b]
	entries = append(entries, e)
	m.data[b] = entries

	return nil
}

func (m *Memory) WriteBatch(b bucket.Bucket, e []log.Entry) error {
	m.dataMu.RLock()
	defer m.dataMu.RUnlock()

	entries := m.data[b]
	entries = append(entries, e...)
	m.data[b] = entries

	return nil
}

func (m *Memory) All(b bucket.Bucket) ([]log.Entry, error) {
	m.dataMu.RLock()
	defer m.dataMu.RUnlock()

	var entries []log.Entry
	for _, e := range m.data {
		entries = append(entries, e...)
	}

	return entries, nil
}

func (m *Memory) Last(b bucket.Bucket, n uint64) ([]log.Entry, error) {
	m.dataMu.RLock()
	defer m.dataMu.RUnlock()

	if b.String() == "" {
		return nil, errors.New("empty bucket not supported in this context with in-memory storage")
	}

	entries := m.data[b]
	cnt := uint64(len(entries))
	if n > cnt {
		return entries, nil
	}

	return entries[cnt-n:], nil
}

func (m *Memory) Count(b bucket.Bucket) (uint64, error) {
	m.dataMu.RLock()
	defer m.dataMu.RUnlock()

	var cnt uint64
	for _, e := range m.data {
		cnt += uint64(len(e))
	}

	return cnt, nil
}
