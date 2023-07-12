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

type NetMessageOption func(msg *NetMessage)

func NetMessageWithProtocol(uri, schema string) NetMessageOption {
	return func(msg *NetMessage) {
		msg.URI = uri
		msg.Schema = schema
	}
}

func NetMessageWithAddress(host string, port int, path string) NetMessageOption {
	return func(msg *NetMessage) {
		msg.Host = host
		msg.Port = int32(port)
		msg.Path = path
	}
}

func NetMessageWithHeaders(key string, values []string) NetMessageOption {
	return func(msg *NetMessage) {
		if msg.Headers == nil {
			msg.Headers = make(map[string]*StringList)
		}
		msg.Headers[key] = &StringList{List: values}
	}
}

func NetMessageWithBody(coding, compress string, body []byte) NetMessageOption {
	return func(msg *NetMessage) {
		msg.Coding = coding
		msg.Compress = compress
		msg.Body = body
	}
}

func (nmsg *NetMessage) With(opts ...NetMessageOption) *NetMessage {
	for _, opt := range opts {
		opt(nmsg)
	}

	return nmsg
}

func NewNetMessage(opts ...NetMessageOption) *NetMessage {
	msg := &NetMessage{}
	for _, opt := range opts {
		opt(msg)
	}

	return msg
}

type NetMessageNativeOption func(msg *NetMessageNative)

func NetMessageNativeWithProtocol(uri, schema string) NetMessageNativeOption {
	return func(msg *NetMessageNative) {
		msg.URI = uri
		msg.Schema = schema
	}
}

func NetMessageNativeWithAddress(host string, port int, path string) NetMessageNativeOption {
	return func(msg *NetMessageNative) {
		msg.Host = host
		msg.Port = int32(port)
		msg.Path = path
	}
}

func NetMessageNativeWithHeaders(key string, values []string) NetMessageNativeOption {
	return func(msg *NetMessageNative) {
		if msg.Headers == nil {
			msg.Headers = make(map[string]*StringList)
		}
		msg.Headers[key] = &StringList{List: values}
	}
}

func NetMessageNativeWithBody(coding, compress string, body interface{}) NetMessageNativeOption {
	return func(msg *NetMessageNative) {
		msg.Coding = coding
		msg.Compress = compress
		msg.Body = body
	}
}

type NetMessageNative struct {
	NetMessage
	Body interface{}
}

func (nmsgn *NetMessageNative) With(opts ...NetMessageNativeOption) *NetMessageNative {
	for _, opt := range opts {
		opt(nmsgn)
	}

	return nmsgn
}

func NetNetMessageNative(opts ...NetMessageNativeOption) *NetMessageNative {
	msg := &NetMessageNative{}
	for _, opt := range opts {
		opt(msg)
	}

	return msg
}
