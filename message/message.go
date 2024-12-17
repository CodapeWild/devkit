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

package message

import (
	"bytes"
	"compress/flate"
	"io"
	reflect "reflect"

	"github.com/CodapeWild/devkit/bufpool"
	"google.golang.org/protobuf/proto"
)

type Msg interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}

func (x *Message) Marshal(_ any) ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message) Unmarshal(_ []byte, _ any) error {
	return proto.Unmarshal(x.Content, x)
}

func Pack(p []byte, method Compression) (*Message, error) {
	var (
		compressedBuf []byte
		err           error
	)
	bufpool.MakeUseOfBuffer(func(buf *bytes.Buffer) {
		if err = method.Compress(buf, p, flate.DefaultCompression); err == nil {
			compressedBuf, err = io.ReadAll(buf)
		}
	})

	if err != nil {
		return nil, err
	} else {
		return &Message{
			Compression: uint32(reflect.ValueOf(method).Uint()),
			Length:      uint32(len(p)),
			Content:     compressedBuf,
		}, nil
	}
}

func Unpack(msg *Message) ([]byte, error) {

}
