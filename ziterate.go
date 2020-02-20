package ziterate

type Iterator interface {
	Next() interface{}
}

func AssignmentTest() Iterator {
	it, _ := smallGroupIteratorFromGroup(zmapGroups[0])
	return it
}
