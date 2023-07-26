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

package io

import (
	"bytes"
	"errors"
	"io"

	"github.com/CodapeWild/devkit/bufpool"
	"google.golang.org/protobuf/proto"
)

type IOMessageOption func(message *IOMessage)

func IOMessageWithDataType(dataType string) IOMessageOption {
	return func(message *IOMessage) {
		message.DataType = dataType
	}
}

func IOMessageWithCoding(coding string) IOMessageOption {
	return func(message *IOMessage) {
		message.Coding = coding
	}
}

func IOMessageWithCompress(compress string) IOMessageOption {
	return func(message *IOMessage) {
		message.Compress = compress
	}
}

func IOMessageWithPayload(payload []byte) IOMessageOption {
	return func(message *IOMessage) {
		message.Payload = payload
	}
}

func (iomsg *IOMessage) With(opts ...IOMessageOption) *IOMessage {
	for _, opt := range opts {
		opt(iomsg)
	}

	return iomsg
}

func NewIOMessage(opts ...IOMessageOption) *IOMessage {
	msg := &IOMessage{}
	for _, opt := range opts {
		opt(msg)
	}

	return msg
}

type IOMessageNativeOption func(message *IOMessageNative)

func IOMsgNativeWithDataType(dataType string) IOMessageNativeOption {
	return func(message *IOMessageNative) {
		message.DataType = dataType
	}
}

func IOMsgNativeWithCoding(coding string) IOMessageNativeOption {
	return func(message *IOMessageNative) {
		message.Coding = coding
	}
}

func IOMsgNativeWithCompress(compress string) IOMessageNativeOption {
	return func(message *IOMessageNative) {
		message.Compress = compress
	}
}

func IOMsgNativeWithPayload(payload interface{}) IOMessageNativeOption {
	return func(message *IOMessageNative) {
		message.Payload = payload
	}
}

type IOMessageNative struct {
	IOMessage
	Payload interface{}
}

func (iomsgn *IOMessageNative) With(opts ...IOMessageNativeOption) *IOMessageNative {
	for _, opt := range opts {
		opt(iomsgn)
	}

	return iomsgn
}

func NewIOMessageNative(opts ...IOMessageNativeOption) *IOMessageNative {
	msg := &IOMessageNative{}
	for _, opt := range opts {
		opt(msg)
	}

	return msg
}

func (x *IOMessageBatch) SetMessages(batch []*IOMessage) {
	x.IOMessageBatch = batch
}

func (x *IOMessageBatch) AppendMessages(batch []*IOMessage) {
	l := len(x.IOMessageBatch) + len(batch)
	buf := make([]*IOMessage, l)
	i := copy(buf, x.IOMessageBatch)
	copy(buf[i:], batch)
	x.IOMessageBatch = buf
}

const delim byte = '\r'

type IOMessageList []*IOMessage

func (x IOMessageList) Encode(list IOMessageList) (bts []byte, err error) {
	bufpool.MakeUseOfBuffer(func(buf *bytes.Buffer) {
		for i := range x {
			var p []byte
			if p, err = proto.Marshal(x[i]); err != nil {
				return
			}
			buf.Write(p)
			buf.WriteByte(delim)
		}
		for i := range list {
			var p []byte
			if p, err = proto.Marshal(list[i]); err != nil {
				return
			}
			buf.Write(p)
			buf.WriteByte(delim)
		}
		bts = buf.Bytes()
	})

	return
}

func (x IOMessageList) Decode(bts []byte) (list IOMessageList, err error) {
	bufpool.MakeUseOfBuffer(func(buf *bytes.Buffer) {
		if _, err = buf.Write(bts); err != nil {
			return
		}

		for {
			var line []byte
			if line, err = buf.ReadBytes(delim); err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
				}

				return
			}

			msg := &IOMessage{}
			if err = proto.Unmarshal(line[:len(line)-1], msg); err != nil {
				return
			}
			list = append(list, msg)
		}
	})

	return
}
