/*
 *   Copyright (c) 2023 CodapeWild
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package set

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
