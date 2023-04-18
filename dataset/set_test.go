package dataset

import (
	"log"
	"testing"
)

func TestIntDataSet(t *testing.T) {
	var temp = []int{2, 4, 6, 8, 3, 4, 3, 6, 7}
	ds := IntDataSet(temp)
	ok := ds.Remove(4, 0)
	log.Println(ok)
	log.Println(ds)
}
