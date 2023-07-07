/*
 *   Copyright (c) 2023 CodapeWild
 *   All rights reserved.

 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:

 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.

 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
 */

package dataset

var _ DataSet = (*IntDataSet)(nil)

type IntDataSet []int

func (x *IntDataSet) Append(value interface{}) bool {
	v, ok := getIntValue(value)
	if !ok {
		return false
	}

	*x = append(*x, v)

	return true
}

func (x *IntDataSet) Replace(old, new interface{}, c int) bool {
	var (
		oldv, newv int
		ok         bool
	)
	if oldv, ok = getIntValue(old); !ok {
		return false
	}
	if newv, ok = getIntValue(new); !ok {
		return false
	}

	ok = false
	if c <= 0 {
		for i, v := range *x {
			if v == oldv {
				(*x)[i] = newv
				ok = true
			}
		}
	} else {
		for i, v := range *x {
			if v == oldv {
				(*x)[i] = newv
				ok = true
				if c--; c == 0 {
					break
				}
			}
		}
	}

	return ok
}

func (x *IntDataSet) Remove(value interface{}, c int) bool {
	v, ok := getIntValue(value)
	if !ok {
		return false
	}

	ok = false
	if c <= 0 {
		for i, u := range *x {
			if v == u {
				(*x) = append((*x)[:i], (*x)[i+1:]...)
				ok = true
			}
		}
	} else {
		for i, u := range *x {
			if v == u {
				(*x) = append((*x)[:i], (*x)[i+1:]...)
				ok = true
				if c--; c == 0 {
					break
				}
			}
		}
	}

	return ok
}

func (x *IntDataSet) Find(value interface{}) (int, bool) {
	v, ok := getIntValue(value)
	if !ok {
		return -1, false
	}

	for i, u := range *x {
		if v == u {
			return i, true
		}
	}

	return -1, false
}

func getIntValue(value interface{}) (int, bool) {
	switch t := value.(type) {
	case int:
		return t, true
	case *int:
		return *t, true
	default:
		return -1, false
	}
}
