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

type SubscribeMessage interface {
	Subscribe(ctx context.Context, handler func(message proto.Message) *IOResponse) error
}

type SubscribeMessageBatch interface {
	SubscribeBatch(ctx context.Context, handler func(batch []proto.Message) *IOResponse) error
}

type SubscribeMessageStream interface {
	SubscribeStream(ctx context.Context, handler func(stream chan proto.Message) *IOResponse) error
}

type FetchMessage interface {
	Fetch(ctx context.Context) (proto.Message, *IOResponse, error)
}

type FetchMessageBatch interface {
	FetchBatch(ctx context.Context) ([]proto.Message, *IOResponse, error)
}

type FetchMessageStream interface {
	FetchStream(ctx context.Context) (chan proto.Message, *IOResponse, error)
}

type PublishAndSubscribeBatch interface {
	PublishMessage
	SubscribeMessageBatch
}
