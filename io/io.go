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

	"github.com/CodapeWild/devkit/message"
)

type PublishMessage interface {
	Publish(ctx context.Context, msg message.Message) *IOResponse
}

type PublishMessageBatch interface {
	PublishBatch(ctx context.Context, batch message.MessageList) *IOResponse
}

type PublishMessageStream interface {
	PublishStream(ctx context.Context, stream chan message.Message) *IOResponse
}

type SubscribeMessageHandler func(ctx context.Context, msg message.Message) *IOResponse

func (h SubscribeMessageHandler) BindContext(ctx context.Context, msg message.Message) SubscribeMessageHandler {
	return func(_ context.Context, _ message.Message) *IOResponse {
		return h(ctx, msg)
	}
}

type SubscribeMessageBatchHandler func(ctx context.Context, batch message.MessageList) *IOResponse

func (h SubscribeMessageBatchHandler) BindContext(ctx context.Context, batch message.MessageList) SubscribeMessageBatchHandler {
	return func(_ context.Context, _ message.MessageList) *IOResponse {
		return h(ctx, batch)
	}
}

type SubscribeMessageStreamHandler func(ctx context.Context, stream chan message.Message, out chan *IOResponse)

func (h SubscribeMessageStreamHandler) BindContext(ctx context.Context, stream chan message.Message, out chan *IOResponse) SubscribeMessageStreamHandler {
	return func(_ context.Context, _ chan message.Message, _ chan *IOResponse) {
		h(ctx, stream, out)
	}
}

type SubscribeMessage interface {
	Subscribe(handler SubscribeMessageHandler) error
}

type SubscribeMessageBatch interface {
	SubscribeBatch(handler SubscribeMessageBatchHandler) error
}

type SubscribeMessageStream interface {
	SubscribeStream(handler SubscribeMessageStreamHandler) error
}

type FetchMessage interface {
	Fetch(ctx context.Context) (message.Message, *IOResponse)
}

type FetchMessageBatch interface {
	FetchBatch(ctx context.Context) (message.MessageList, *IOResponse)
}

type PubAndSub interface {
	PublishMessage
	SubscribeMessage
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

type PubBatchAndFetchBatch interface {
	PublishMessageBatch
	FetchMessageBatch
}

type PubPubBatchAndFetchFetchBatch interface {
	PublishMessage
	PublishMessageBatch
	FetchMessage
	FetchMessageBatch
}
