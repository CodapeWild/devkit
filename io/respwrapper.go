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

var (
	IOSuccess     = NewIOResponse(IOStatus_IOSuccess)
	InputSuccess  = NewIOResponse(IOStatus_IOK, IORespWithMessage("input success"))
	InputBusy     = NewIOResponse(IOStatus_IBusy, IORespWithMessage("input busy"))
	InputTimeout  = NewIOResponse(IOStatus_ITimeout, IORespWithMessage("input timeout"))
	InputFailed   = NewIOResponse(IOStatus_IFailed, IORespWithMessage("input failed"))
	OutputSuccess = NewIOResponse(IOStatus_OOK, IORespWithMessage("output success"))
	OutputBusy    = NewIOResponse(IOStatus_OBusy, IORespWithMessage("output busy"))
	OutputTimeout = NewIOResponse(IOStatus_OTimeout, IORespWithMessage("output timeout"))
	OutputFailed  = NewIOResponse(IOStatus_OFailed, IORespWithMessage("output failed"))
)

type IOResponseOption func(ioresp *IOResponse)

func IORespWithStatus(status IOStatus) IOResponseOption {
	return func(ioresp *IOResponse) {
		ioresp.Status = status
	}
}

func IORespWithMessage(msg string) IOResponseOption {
	return func(ioresp *IOResponse) {
		ioresp.Message = msg
	}
}

func IORespWithPayload(coding string, payload []byte) IOResponseOption {
	return func(ioresp *IOResponse) {
		ioresp.Coding = coding
		ioresp.Payload = payload
	}
}

func (ioresp *IOResponse) With(opts ...IOResponseOption) *IOResponse {
	for _, opt := range opts {
		opt(ioresp)
	}

	return ioresp
}

func NewIOResponse(status IOStatus, opts ...IOResponseOption) *IOResponse {
	ioresp := &IOResponse{Status: status}
	for _, opt := range opts {
		opt(ioresp)
	}

	return ioresp
}
