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

package directory

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestSeqDirSave(t *testing.T) {
	seqDir, err := OpenSequentialDirectory("./test")
	if err != nil {
		t.Fatal(err.Error())
	}

	var (
		bufSize   = 1000
		saveTimes = 10
	)
	for i := 0; i < saveTimes; i++ {
		t.Run("save:"+strconv.Itoa(i), func(t *testing.T) {
			buf := make([]byte, bufSize)
			_, err := rand.Read(buf)
			if err != nil {
				t.Fatal(err.Error())
			}
			if err := seqDir.Save("", bytes.NewBuffer(buf)); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}

func TestSeqDirDelete(t *testing.T) {
	seqDir, err := OpenSequentialDirectory("./test")
	if err != nil {
		t.Fatal(err.Error())
	}

	var (
		bufSize     = 1000
		deleteTimes = 5
	)
	for i := 0; i < deleteTimes; i++ {
		t.Run("remove:"+strconv.Itoa(i), func(t *testing.T) {
			f, err := seqDir.Open("")
			if err != nil {
				if errors.Is(err, ErrDirEmpty) {
					return
				}
				t.Fatal(err.Error())
			}
			defer f.Close()

			buf := make([]byte, bufSize)
			if _, err = f.Read(buf); err != nil {
				t.Fatal(err.Error())
			}
			fmt.Println(base64.RawStdEncoding.EncodeToString(buf))

			if err := seqDir.Delete(""); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}

func TestSeqDirOpenAndDelete(t *testing.T) {
	seqDir, err := OpenSequentialDirectory("./test")
	if err != nil {
		t.Fatal(err.Error())
	}

	for {
		bts, err := seqDir.OpenAndDelete("")
		if err != nil {
			if errors.Is(err, ErrDirEmpty) {
				break
			}
			t.Fatal(err.Error())
		}
		fmt.Println(base64.RawStdEncoding.EncodeToString(bts.Bytes()))
	}
}
