package ziterate

type Iterator interface {
	Next() interface{}
}

type IPv4Iterator struct {
}
