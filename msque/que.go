/*
 *   Copyright (c) 2023 CodapeWild
 *   All rights reserved.

 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:

 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.

 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
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
