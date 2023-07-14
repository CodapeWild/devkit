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

	"google.golang.org/protobuf/proto"
)

type PublishMessage interface {
	Publish(ctx context.Context, message proto.Message) (*IOResponse, error)
}

type PublishMessageBatch interface {
	PublishBatch(ctx context.Context, batch []proto.Message) (*IOResponse, error)
}

type PublishMessageStream interface {
	PublishStream(ctx context.Context, stream chan proto.Message) (*IOResponse, error)
}

type SubscribeMessageHandler func(message proto.Message) *IOResponse

type SubscribeMessage interface {
	Subscribe(handler SubscribeMessageHandler) error
}

type SubscribeMessageBatchHandler func(batch []proto.Message) *IOResponse

type SubscribeMessageBatch interface {
	SubscribeBatch(handler SubscribeMessageBatchHandler) error
}

type SubscribeMessageStreamHandler func(stream chan proto.Message, out chan *IOResponse)

type SubscribeMessageStream interface {
	SubscribeStream(handler SubscribeMessageStreamHandler) error
}

type FetchMessage interface {
	Fetch(ctx context.Context) (proto.Message, *IOResponse, error)
}

type FetchMessageBatch interface {
	FetchBatch(ctx context.Context) ([]proto.Message, *IOResponse, error)
}

type PubAndSubBatch interface {
	PublishMessage
	SubscribeMessageBatch
}

type PubPubBatchAndSubSubBatch interface {
	PublishMessage
	PublishMessageBatch
	SubscribeMessage
	SubscribeMessageBatch
}

type PubBatchAndSubBatch interface {
	PublishMessageBatch
	SubscribeMessageBatch
}

type PubStreamAndSubStream interface {
	PublishMessageStream
	SubscribeMessageStream
}
