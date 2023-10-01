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
	"context"

	"github.com/CodapeWild/devkit/id"
	"github.com/CodapeWild/devkit/message"
)

var _ PubBatchAndFetchBatch = (*FileChain)(nil)

const (
	_headFileName = "head"
	_tailFileName = "tail"
)

type FileChain struct {
	path     string
	pageSize int
	idflk    *id.IDFlaker
}

func (fc *FileChain) PublishBatch(ctx context.Context, batch message.MessageList) *IOResponse {

}

func (fc *FileChain) FetchBatch(ctx context.Context) (message.MessageList, *IOResponse) {

}

func OpenFileChain(path string, pageSize int) (*FileChain, error) {

}
