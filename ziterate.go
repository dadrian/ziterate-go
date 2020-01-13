package ziterate

type Iterator interface {
	Next() (interface{}, error)
}
