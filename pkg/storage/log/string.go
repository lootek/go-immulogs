package log

import (
	"encoding/json"
)

type Entry interface {
	String() string
	Bytes() []byte
}

type entry string

func (e *entry) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	*e = entry(s)
	return nil
}

func FromString(s string) Entry {
	return entry(s)
}

func FromBytes(b []byte) Entry {
	return entry(b)
}

func (e entry) String() string {
	return string(e)
}

func (e entry) Bytes() []byte {
	return []byte(e)
}
