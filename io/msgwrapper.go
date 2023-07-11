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

type IOMessageOption func(msg *IOMessage)

func IOMessageWithDataType(dataType string) IOMessageOption {
	return func(msg *IOMessage) {
		msg.DataType = dataType
	}
}

func IOMessageWithCoding(coding string) IOMessageOption {
	return func(msg *IOMessage) {
		msg.Coding = coding
	}
}

func IOMessageWithCompress(compress string) IOMessageOption {
	return func(msg *IOMessage) {
		msg.Compress = compress
	}
}

func IOMessageWithPayload(payload []byte) IOMessageOption {
	return func(msg *IOMessage) {
		msg.Payload = payload
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

type IOMessageNativeOption func(msg *IOMessageNative)

func IOMsgNativeWithDataType(dataType string) IOMessageNativeOption {
	return func(msg *IOMessageNative) {
		msg.DataType = dataType
	}
}

func IOMsgNativeWithCoding(coding string) IOMessageNativeOption {
	return func(msg *IOMessageNative) {
		msg.Coding = coding
	}
}

func IOMsgNativeWithCompress(compress string) IOMessageNativeOption {
	return func(msg *IOMessageNative) {
		msg.Compress = compress
	}
}

func IOMsgNativeWithPayload(payload interface{}) IOMessageNativeOption {
	return func(msg *IOMessageNative) {
		msg.Payload = payload
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
