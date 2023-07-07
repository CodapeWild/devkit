/*
 *   Copyright (c) 2023 CodapeWild
 *   All rights reserved.

 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:

 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.

 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
 */

package slice

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
