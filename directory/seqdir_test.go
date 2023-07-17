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
	"fmt"
	"strconv"
	"testing"
)

func TestSeqDirReadAndWrite(t *testing.T) {
	seqdir, err := OpenSequentialDirectory("./test")
	if err != nil {
		t.Fatal(err.Error())
	}

	for i := 0; i < 10; i++ {
		t.Run("TestSeqDirReadAndWrite-Write:"+strconv.Itoa(i), func(t *testing.T) {
			buf := make([]byte, 1000)
			_, err := rand.Read(buf)
			if err != nil {
				t.Fatal(err.Error())
			}
			if err := seqdir.Save("", bytes.NewBuffer(buf)); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
	for i := 0; i < 5; i++ {
		t.Run("TestSeqDirReadAndWrite-Read:"+strconv.Itoa(i), func(t *testing.T) {
			f, err := seqdir.Open("")
			if err != nil {
				t.Fatal(err.Error())
			}
			buf := make([]byte, 1000)
			if _, err = f.Read(buf); err != nil {
				t.Fatal(err.Error())
			}
			f.Close()
			fmt.Println(base64.RawStdEncoding.EncodeToString(buf))
			if err := seqdir.Delete(""); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}
