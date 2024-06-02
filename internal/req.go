// Copyright 2024 dnsdk Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	q "github.com/google/go-querystring/query"
)

type (
	DomainListReq struct {
		PageNumber string `url:"page_number,omitempty"` // 页码
		PageSize   string `url:"page_size,omitempty"`   // 每页数量
		Name       string `url:"name,omitempty"`        // 域名名称
	}
	RecordTypeListReq     struct{}
	RecordLineListReq     struct{}
	RecordListReq         struct{}
	RecordAddReq          struct{}
	RecordUpdateReq       struct{}
	RecordDeleteReq       struct{}
	RecordStatusUpdateReq struct{}
)

func (r *DomainListReq) url() string { v, _ := q.Values(r); return v.Encode() }
