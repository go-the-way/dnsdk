// Copyright 2024 dnsdk Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/alibabacloud-go/tea/tea"
)

var ErrNotSupportedOperation = errors.New("不支持的操作")

func toLineMap(lines []LineListRespLine) map[string]string {
	m := make(map[string]string)
	for _, l := range lines {
		m[l.Id] = l.Name
	}
	return m
}

func i2s(i int) string { return fmt.Sprintf("%d", i) }

func toUint(str string) uint {
	i, _ := strconv.ParseUint(str, 10, 64)
	return uint(i)
}

func toUint64Ptr(str string) *uint64 {
	i, _ := strconv.ParseUint(str, 10, 64)
	return tea.Uint64(i)
}
