package dataset

import "sort"

type DataSet interface {
	Append(value interface{}) bool
	Replace(old, new interface{}, c int) bool
	Remove(value interface{}, c int) bool
	Find(value interface{}) (int, bool)
}

type OrderedDataSet interface {
	sort.Interface
	DataSet
}
