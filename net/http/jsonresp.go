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

package http

import (
	"encoding/json"
	"net/http"
)

type JSONRespWriter interface {
	http.ResponseWriter
	WriteJSON(v interface{}) (int, error)
}

type JSONResponse struct {
	http.ResponseWriter
}

func (jresp *JSONResponse) Header() http.Header {
	return jresp.Header()
}

func (jresp *JSONResponse) Write(bts []byte) (int, error) {
	return jresp.ResponseWriter.Write(bts)
}

func (jresp *JSONResponse) WriteHeader(statusCode int) {
	jresp.ResponseWriter.WriteHeader(statusCode)
}

func (jresp *JSONResponse) WriteJSON(v interface{}) (int, error) {
	bts, err := json.Marshal(v)
	if err != nil {
		jresp.ResponseWriter.WriteHeader(http.StatusBadRequest)

		return 0, err
	}

	jresp.ResponseWriter.Header().Set("Content-Type", "application/json")

	return jresp.ResponseWriter.Write(bts)
}

type JSONRespMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Payload []byte `json:"payload"`
}

func (jmsg *JSONRespMessage) WriteBy(resp http.ResponseWriter) (int, error) {
	return (&JSONResponse{resp}).WriteJSON(jmsg)
}

func NewJSONRespMessage(status int, msg string, code string, payload []byte) *JSONRespMessage {
	return &JSONRespMessage{
		Status:  status,
		Message: msg,
		Code:    code,
		Payload: payload,
	}
}