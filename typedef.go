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

package dnsdk

import "github.com/go-the-way/dnsdk/internal"

type (
	Api = internal.Api

	DomainListReq   = internal.DomainListReq
	DomainAddReq    = internal.DomainAddReq
	DomainDeleteReq = internal.DomainDeleteReq

	RecordListReq    = internal.RecordListReq
	RecordAddReq     = internal.RecordAddReq
	RecordUpdateReq  = internal.RecordUpdateReq
	RecordDeleteReq  = internal.RecordDeleteReq
	RecordEnableReq  = internal.RecordEnableReq
	RecordDisableReq = internal.RecordDisableReq

	LineListResp         = internal.LineListResp
	LineListRespLine     = internal.LineListRespLine
	DomainListResp       = internal.DomainListResp
	DomainListRespDomain = internal.DomainListRespDomain
	DomainAddResp        = internal.DomainAddResp

	RecordListResp       = internal.RecordListResp
	RecordListRespRecord = internal.RecordListRespRecord
	RecordAddResp        = internal.RecordAddResp
	RecordUpdateResp     = internal.RecordUpdateResp
)
