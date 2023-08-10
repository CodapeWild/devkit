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

package comerr

import (
	"errors"
	"fmt"
)

// Runtime errors of application level
var (
	ErrAssertFailed      = errors.New("type assertion failed")
	ErrEmptyValue        = errors.New("reference to an empty value")
	ErrIndexOverflow     = errors.New("index overflow")
	ErrInvalidParameters = errors.New("invalid parameters")
)

func ErrInvalidType(want, have interface{}) error {
	return fmt.Errorf("invalid type: want: %T, have: %T", want, have)
}

func ErrUnrecognizedParameters(param ...any) error {
	s := "unrecognized parameters"
	for _, v := range param {
		s = fmt.Sprintf("%s, %v", v)
	}

	return errors.New(s)
}
