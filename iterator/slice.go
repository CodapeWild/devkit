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

package iterator

func Contains(src []string, target string) bool {
	for i := range src {
		if src[i] == target {
			return true
		}
	}

	return false
}

func Include(src, target []string) bool {
	for i := range target {
		find := false
		for j := range src {
			if target[i] == src[j] {
				find = true
				break
			}
		}
		if !find {
			return false
		}
	}

	return true
}

func AtLeast(src, target []string, n int) bool {
	c := 0
	for i := range target {
		for j := range src {
			if target[i] == src[j] {
				if c++; c == n {
					return true
				}
			}
		}
	}

	return false
}

func AtMost(src, target []string, n int) bool {
	c := 0
	for i := range target {
		for j := range src {
			if target[i] == src[j] {
				if c++; c > n {
					return false
				}
			}
		}
	}

	return true
}

func Range(src, target []string, min, max int) bool {
	c := 0
	for i := range target {
		for j := range src {
			if target[i] == src[j] {
				if c++; c > max {
					return false
				}
			}
		}
	}

	return c >= min
}
