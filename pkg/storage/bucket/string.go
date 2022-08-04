package bucket

type Bucket interface {
	String() string
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
