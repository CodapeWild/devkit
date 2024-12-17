/*
 *   Copyright (c) 2024 CodapeWild
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

package message

import (
	"compress/flate"
	"io"
)

type Compression interface {
	Compress(dst io.Writer, src []byte, level int) error
	Decompress(dst []byte, src io.Reader) error
}

const (
	DefCompMethod NoCompression = NoCompression(0)
	FlateMethod   Flate         = Flate(1)
)

type NoCompression uint32

func (NoCompression) Compress(dst io.Writer, src []byte, level int) error {
	_, err := dst.Write(src)

	return err
}

func (NoCompression) Decompress(dst []byte, src io.Reader) error {
	_, err := io.ReadFull(src, dst)

	return err
}

type Flate uint32

func (Flate) Compress(compressing io.Writer, raw []byte, level int) error {
	zw, err := flate.NewWriter(compressing, level)
	if err != nil {
		return err
	}
	if _, err = zw.Write(raw); err != nil {
		return err
	}

	return zw.Close()
}

func (Flate) Decompress(raw []byte, compressing io.Reader) error {
	zr := flate.NewReader(compressing)
	if _, err := zr.Read(raw); err != nil {
		return err
	}

	return zr.Close()
}
