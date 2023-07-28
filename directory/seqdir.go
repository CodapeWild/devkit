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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/CodapeWild/devkit/comerr"
	"github.com/CodapeWild/devkit/id"
	"github.com/CodapeWild/devkit/set"
)

var _ Directory = (*SequentialDirectory)(nil)

type SequentialDirectory struct {
	path  string // directory path
	idflk *id.IDFlaker
	stque *set.SingleThreadQueue
}

func (seqdir *SequentialDirectory) List() ([]fs.DirEntry, error) {
	return os.ReadDir(seqdir.path)
}

func (seqdir *SequentialDirectory) Open(_ string) (fs.File, error) {
	value := seqdir.stque.Peek()
	if value == nil {
		return nil, ErrDirEmpty
	}
	id, ok := value.(*id.ID)
	if !ok {
		return nil, comerr.ErrAssertFailed
	}

	return os.Open(seqdir.formatPath(id.String('-')))
}

func (seqdir *SequentialDirectory) OpenWithID(index string) (fs.File, error) {
	return os.Open(index)
}

func (seqd *SequentialDirectory) Save(_ string, r io.Reader) error {
	id := seqd.idflk.NextID()
	f, err := os.Create(seqd.formatPath(id.String('-')))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	seqd.stque.Push(id)

	return err
}

func (seqdir *SequentialDirectory) Delete(_ string) error {
	value, _ := seqdir.stque.Pop()
	id, ok := value.(*id.ID)
	if !ok {
		return comerr.ErrAssertFailed
	}

	return os.Remove(seqdir.formatPath(id.String('-')))
}

func (seqdir *SequentialDirectory) formatPath(id string) string {
	return fmt.Sprintf("%s/.%s", seqdir.path, id)
}

func OpenSequentialDirectory(path string) (*SequentialDirectory, error) {
	if err := Exist(path); err != nil {
		if !errors.Is(err, ErrNotDir) {
			if err = os.MkdirAll(path, 0755); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	seqdir := &SequentialDirectory{path: path}
	entries, err := seqdir.List()
	if err != nil {
		return nil, err
	}
	seqdir.idflk = id.NewIDFlaker()
	seqdir.stque = set.NewSingleThreadQueue(10)

	ids, err := dirEntriesToIDs(entries)
	if err != nil {
		return nil, err
	}
	sort.Sort(ids)
	for _, id := range ids {
		if err = seqdir.stque.Push(id); err != nil {
			log.Println(err.Error())
		}
	}

	return seqdir, nil
}

func dirEntriesToIDs(entries []fs.DirEntry) (id.IDs, error) {
	var ids id.IDs
	for _, entry := range entries {
		if id, err := id.FromString(strings.TrimPrefix(entry.Name(), "."), '-'); err != nil {
			return nil, err
		} else {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
