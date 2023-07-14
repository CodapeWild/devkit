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
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"sync"
)

var _ Directory = (*SequentialDirectory)(nil)

type SequentialDirectory struct {
	sync.Mutex
	path               string // directory path
	minindex, maxindex int    // start from -1 both
}

func (seqd *SequentialDirectory) Index() (min, max int) {
	return seqd.minindex, seqd.maxindex
}

func (seqd *SequentialDirectory) List() ([]fs.DirEntry, error) {
	return os.ReadDir(seqd.path)
}

func (seqd *SequentialDirectory) Open(name string) (fs.File, error) {
	if len(name) != 0 || name[0] != '.' {
		return nil, os.ErrNotExist
	}

	return os.Open(fmt.Sprintf("%s/%s", seqd.path, name))
}

func (seqd *SequentialDirectory) OpenWithIndex(index int) (fs.File, error) {
	return seqd.Open(fmt.Sprintf(".%d", index))
}

func (seqd *SequentialDirectory) Save(_ string, r io.Reader) error {
	f, err := os.Create(fmt.Sprintf("%s/.%d", seqd.path, seqd.nextID()))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)

	return err
}

func (seqd *SequentialDirectory) Delete(_ string) error {
	if seqd.minindex == -1 {
		return os.ErrNotExist
	}

	seqd.Lock()
	defer seqd.Unlock()

	if err := os.Remove(fmt.Sprintf("%s/.%d", seqd.path, seqd.minindex)); err != nil {
		return err
	}
	seqd.minindex++

	if seqd.minindex > seqd.maxindex {
		seqd.maxindex = -1
		seqd.minindex = -1
	}

	return nil
}

func (seqd *SequentialDirectory) nextID() int {
	seqd.Lock()
	defer seqd.Unlock()

	seqd.maxindex++

	return seqd.maxindex
}

func OpenSequentialDirectory(path string) (*SequentialDirectory, error) {
	seqd := &SequentialDirectory{path: path}
	files, err := seqd.List()
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		seqd.minindex = -1
		seqd.maxindex = -1

		return seqd, nil
	}

	sort.Sort(seqFiles(files))
	seqd.minindex, _ = strconv.Atoi(files[0].Name()[1:])
	seqd.maxindex, _ = strconv.Atoi(files[len(files)-1].Name()[1:])

	return seqd, nil
}

type seqFiles []fs.DirEntry

func (seq seqFiles) Len() int {
	return len(seq)
}

func (seq seqFiles) Less(i, j int) bool {
	x, _ := strconv.Atoi(seq[i].Name()[1:])
	y, _ := strconv.Atoi(seq[i].Name()[1:])

	return x < y
}

func (seq seqFiles) Swap(i, j int) {
	seq[i], seq[j] = seq[j], seq[i]
}
