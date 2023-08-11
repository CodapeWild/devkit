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

import "errors"

var (
	ErrIOClosed      = errors.New("io closed")
	ErrIOUncompleted = errors.New("io init uncompleted")
)

var (
	IOSuccess      = NewIOResponse(IOStatus_IOSuccess, IORespWithMessage("io success"))
	IOClosed       = NewIOResponse(IOStatus_IOClosed, IORespWithMessage("io closed"))
	IOUncompleted  = NewIOResponse(IOStatus_IOUncompleted, IORespWithMessage("io init uncompleted"))
	IOWrongMsgType = NewIOResponse(IOStatus_IOWrongMsgType, IORespWithMessage("message assertion failed"))
	InputSuccess   = NewIOResponse(IOStatus_IOK, IORespWithMessage("input success"))
	InputBusy      = NewIOResponse(IOStatus_IBusy, IORespWithMessage("input busy"))
	InputTimeout   = NewIOResponse(IOStatus_ITimeout, IORespWithMessage("input timeout"))
	InputFailed    = NewIOResponse(IOStatus_IFailed, IORespWithMessage("input failed"))
	OutputSuccess  = NewIOResponse(IOStatus_OOK, IORespWithMessage("output success"))
	OutputEmpty    = NewIOResponse(IOStatus_OEMPTY, IORespWithMessage("output empty"))
	OutputBusy     = NewIOResponse(IOStatus_OBusy, IORespWithMessage("output busy"))
	OutputTimeout  = NewIOResponse(IOStatus_OTimeout, IORespWithMessage("output timeout"))
	OutputFailed   = NewIOResponse(IOStatus_OFailed, IORespWithMessage("output failed"))
)

type IOResponseOption func(ioresp *IOResponse)

func IORespWithStatus(status IOStatus) IOResponseOption {
	return func(resp *IOResponse) {
		resp.Status = status
	}
}

func IORespWithMessage(s string) IOResponseOption {
	return func(resp *IOResponse) {
		resp.Message = s
	}
}

func IORespWithPayload(coding string, payload []byte) IOResponseOption {
	return func(resp *IOResponse) {
		resp.Coding = coding
		resp.Payload = payload
	}
}

func (x *IOResponse) With(opts ...IOResponseOption) *IOResponse {
	for _, opt := range opts {
		opt(x)
	}

	return x
}

func (x *IOResponse) IS(target *IOResponse) bool {
	return x.Status == target.Status
}

func NewIOResponse(status IOStatus, opts ...IOResponseOption) *IOResponse {
	resp := &IOResponse{Status: status}
	for _, opt := range opts {
		opt(resp)
	}

	return resp
}
