// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"os"
	"strconv"
	"strings"
)

// FileExists determines if file exists.
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func IntSliceToString(i []int) []string {
	s := make([]string, len(i), len(i))
	for in, num := range i {
		s[in] = strconv.Itoa(num)
	}
	return s
}

func GenerateSlug(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.ToLower(raw)
	raw = strings.Replace(raw, "_", "-", -1)
	return strings.Replace(raw, " ", "-", -1)
}
