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

package id

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidFormat = errors.New("not valid ID format")

type ID struct {
	high, low int64
}

func (id *ID) Int64() (high, low int64) {
	return id.high, id.low
}

func (id *ID) String(sep byte) string {
	return fmt.Sprintf("%d%c%d", id.high, sep, id.low)
}

func (id *ID) Bytes() [16]byte {
	var bts = [16]byte{}
	bts[0] = byte(id.low)
	bts[1] = byte(id.low >> 8)
	bts[2] = byte(id.low >> 16)
	bts[3] = byte(id.low >> 24)
	bts[4] = byte(id.low >> 32)
	bts[5] = byte(id.low >> 40)
	bts[6] = byte(id.low >> 48)
	bts[7] = byte(id.low >> 56)
	bts[8] = byte(id.high)
	bts[9] = byte(id.high >> 8)
	bts[10] = byte(id.high >> 16)
	bts[11] = byte(id.high >> 24)
	bts[12] = byte(id.high >> 32)
	bts[13] = byte(id.high >> 40)
	bts[14] = byte(id.high >> 48)
	bts[15] = byte(id.high >> 56)

	return bts
}

func FromInt64(high, low int64) *ID {
	return &ID{high: high, low: low}
}

func FromString(id string, sep byte) (*ID, error) {
	c := strings.IndexByte(id, sep)
	if c < 0 {
		return nil, ErrInvalidFormat
	}

	high, err := strconv.ParseInt(id[:c], 10, 64)
	if err != nil {
		return nil, err
	}
	low, err := strconv.ParseInt(id[c+1:], 10, 64)
	if err != nil {
		return nil, err
	}

	return &ID{high: high, low: low}, nil
}

func FromBytes(bts [16]byte) *ID {
	highb := bts[8:]
	lowb := bts[:8]

	return &ID{
		high: int64(highb[0]) | int64(highb[1])<<8 | int64(highb[2])<<16 | int64(highb[3])<<24 | int64(highb[4])<<32 | int64(highb[5])<<40 | int64(highb[6])<<48 | int64(highb[7])<<56,
		low:  int64(lowb[0]) | int64(lowb[1])<<8 | int64(lowb[2])<<16 | int64(lowb[3])<<24 | int64(lowb[4])<<32 | int64(lowb[5])<<40 | int64(lowb[6])<<48 | int64(lowb[7])<<56,
	}
}

type IDs []*ID

func (ids IDs) Len() int {
	return len(ids)
}

func (ids IDs) Less(i, j int) bool {
	return (ids[i].high < ids[j].high) && (ids[i].low < ids[j].low)
}

func (ids IDs) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}
