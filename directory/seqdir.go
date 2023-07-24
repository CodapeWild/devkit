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

func (seqd *SequentialDirectory) List() ([]fs.DirEntry, error) {
	return os.ReadDir(seqd.path)
}

func (seqd *SequentialDirectory) Open(_ string) (fs.File, error) {
	t := seqd.stque.Peek()
	if t == nil {
		return nil, ErrDirEmpty
	}
	id, ok := t.(*id.ID)
	if !ok {
		return nil, comerr.ErrAssertFailed
	}

	return os.Open(seqd.formatPath(id.String('-')))
}

func (seqd *SequentialDirectory) OpenWithIndex(index string) (fs.File, error) {
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

func (seqd *SequentialDirectory) Delete(_ string) error {
	return seqd.stque.AsyncPop(func(value any) {
		id, ok := value.(*id.ID)
		if !ok {
			log.Println(comerr.ErrAssertFailed.Error())

			return
		}
		if err := os.Remove(seqd.formatPath(id.String('-'))); err != nil {
			log.Println(err.Error())
		}
	})
}

func (seqd *SequentialDirectory) formatPath(index string) string {
	return fmt.Sprintf("%s/.%s", seqd.path, index)
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

	seqd := &SequentialDirectory{path: path}
	entries, err := seqd.List()
	if err != nil {
		return nil, err
	}
	seqd.idflk = id.NewIDFlaker()
	seqd.stque = set.NewSingleThreadQueue(10)

	ids := dirEntriesToIDs(entries)
	sort.Sort(ids)
	for _, id := range ids {
		if err = seqd.stque.Push(id); err != nil {
			log.Println(err.Error())
		}
	}

	return seqd, nil
}

func dirEntriesToIDs(entries []fs.DirEntry) id.IDs {
	var ids id.IDs
	for _, entry := range entries {
		if id, err := id.FromString(entry.Name(), '-'); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}
