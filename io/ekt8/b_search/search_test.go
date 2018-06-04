package b_search

import (
	"fmt"
	"sort"
	"testing"
)

type Object struct {
	Integer int
}

type Objects struct {
	Objects []Object
}

func (a Object) Value() string {
	return fmt.Sprintf("%d", a.Integer)
}

func (objs Objects) Len() int {
	return len(objs.Objects)
}

func (objs Objects) Swap(i, j int) {
	objs.Objects[i], objs.Objects[j] = objs.Objects[j], objs.Objects[i]
}

func (objs Objects) Less(i, j int) bool {
	return objs.Objects[i].Integer < objs.Objects[j].Integer
}

func (objs Objects) Index(index int) Interface {
	if index > objs.Len() || index < 0 {
		panic("Index out of bound.")
	}
	return objs.Objects[index]
}

func TestSearch(t *testing.T) {
	a := Object{1}
	b := Object{2}
	c := Object{3}
	d := Object{5}
	e := Object{4}
	objs := Objects{make([]Object, 0)}
	objs.Objects = append(objs.Objects, a)
	objs.Objects = append(objs.Objects, d)
	objs.Objects = append(objs.Objects, c)
	objs.Objects = append(objs.Objects, e)
	objs.Objects = append(objs.Objects, b)
	sort.Sort(objs)
	if Search(e, objs) != 3 {
		t.Fail()
	}
	f := Object{7}
	if Search(f, objs) != -1 {
		t.Fail()
	}
}
