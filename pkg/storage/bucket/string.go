package bucket

type Bucket interface {
	String() string
	Bytes() []byte
}

type bucket struct {
	string
}

func NewBucket(s string) Bucket {
	return bucket{s}
}

func (b bucket) String() string {
	return b.string
}

func (b bucket) Bytes() []byte {
	return []byte(b.string)
}
