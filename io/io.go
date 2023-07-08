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

type PublishData interface {
	Publish(data any) error
}

type PublishDataStream interface {
	Publish(stream chan any) error
}

type SubscribeData interface {
	Subscribe(func(data any)) error
}

type SubscribeDataStream interface {
	Subscribe(func(stream chan any)) error
}

type FechData interface {
	Fetch() (any, error)
}

type FetchDataStream interface {
	Fetch() (chan any, error)
}
