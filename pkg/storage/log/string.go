package log

type Entry interface {
	String() string
}

type entry string

func (e *entry) UnmarshalJSON(bytes []byte) error {
	s := entry(bytes)
	e = &s
	return nil
}

func FromString(s string) Entry {
	return entry(s)
}

func FromBinary(b []byte) Entry {
	return entry(b)
}

func (e entry) String() string {
	return string(e)
}
