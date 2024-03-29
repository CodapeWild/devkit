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

type Message interface {
	Encode() (p []byte, err error)
	Decode(p []byte) (err error)
}

type MessageList interface {
	Message
	Length() int
	Foreach(handler func(k int, msg Message) bool)
}

type MessageSet interface {
	Message
	Length() int
	Foreach(handler func(k any, msg Message) bool)
	Get(k any) (msg Message, ok bool)
	Set(k any, msg Message)
}
