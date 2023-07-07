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

package msque

import (
	"encoding/binary"
	"errors"
	"time"
)

type MessageHandler interface {
	HandleMessage(bts []byte) (err error)
}

type MessageHandlerFunc func(bts []byte) (err error)

func (h MessageHandlerFunc) HandleMessage(bts []byte) (err error) {
	return h.HandleMessage(bts)
}

type MessageQueue interface {
	Publish(topic string, bts []byte) (err error)
	Subscribe(topic string, handler MessageHandler)
}

/*
	Memory model in message
| len            | ts                | retry         | payload      |
| -------------- | ----------------- | ------------- | ------------ |
| payload length | enqueue timestamp | enqueue times | message body |
| uint32         | int64             | uint8         | []byte       |
| 4294967295     | -                 | 255           | 4095MB       |
*/
const (
	MaxLen    = 1<<32 - 1
	MaxRetry  = 1<<8 - 1
	HeaderLen = 4 + 8 + 1
)

func newMessage(payload []byte) (message, error) {
	pl := len(payload)
	if pl > MaxLen {
		return nil, errors.New("max message size overflow")
	}

	bts := make([]byte, HeaderLen+pl)
	binary.BigEndian.PutUint32(bts, uint32(pl))
	binary.BigEndian.PutUint64(bts[4:], uint64(time.Now().UnixNano()))
	copy(bts[HeaderLen:], payload)

	return bts, nil
}

type message []byte

func (ms message) len() int {
	return len(ms)
}

func (ms message) valid() bool {
	return len(ms) > HeaderLen && (binary.BigEndian.Uint32(ms) == uint32(len(ms)-HeaderLen))
}

func (ms message) parse() (n uint32, ts time.Duration, retry uint8, payload []byte, err error) {
	if !ms.valid() {
		err = errors.New("invalid message format")

		return
	}

	n = binary.BigEndian.Uint32(ms)
	ts = time.Duration(binary.BigEndian.Uint64(ms[4:]))
	retry = uint8(ms[12])
	payload = ms[HeaderLen:]

	return
}

func (ms message) failOnce() (times uint8, ok bool) {
	times = ms[12] + 1
	ms[12] = times

	return times, times != 0
}
