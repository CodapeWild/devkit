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

package list

import "sync"

type StackNode struct {
	value any
	next  *StackNode
}

func (sn *StackNode) Next() *StackNode {
	return sn.next
}

func (sn *StackNode) Value() any {
	return sn.value
}

type Stack struct {
	top *StackNode
	l   int
}

func (s *Stack) Peek() *StackNode {
	return s.top
}

func (s *Stack) Push(v any) *StackNode {
	sn := &StackNode{value: v, next: s.top}
	s.top = sn
	s.l++

	return sn
}

func (s *Stack) Pop() *StackNode {
	top := s.top
	if s.top != nil {
		s.top = s.top.next
		s.l--
	}

	return top
}

func (s *Stack) Len() int {
	return s.l
}

func NewStack() *Stack {
	return &Stack{top: nil, l: 0}
}

type SyncStack struct {
	sync.Mutex
	*Stack
}

func (ss *SyncStack) Peek() *StackNode {
	ss.Lock()
	defer ss.Unlock()

	return ss.Stack.Peek()
}

func (ss *SyncStack) Push(v any) *StackNode {
	ss.Lock()
	defer ss.Unlock()

	return ss.Stack.Push(v)
}

func (ss *SyncStack) Pop() *StackNode {
	ss.Lock()
	defer ss.Unlock()

	return ss.Stack.Pop()
}

func (ss *SyncStack) Len() int {
	ss.Lock()
	ss.Unlock()

	return ss.Stack.Len()
}
