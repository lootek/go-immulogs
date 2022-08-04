package bucket

type Bucket struct {
	string
}

func NewBucket(s string) Bucket {
	return Bucket{s}
}
