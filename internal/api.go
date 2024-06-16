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

package internal

type Api interface {
	LineList() (resp LineListResp)                                       // 线路列表
	LineDefault() (resp LineListRespLine)                                // 线路默认
	DomainList(req DomainListReq) (resp DomainListResp, err error)       // 域名列表
	DomainAdd(req DomainAddReq) (resp DomainAddResp, err error)          // 域名添加
	DomainDelete(req DomainDeleteReq) (err error)                        // 域名删除
	RecordList(req RecordListReq) (resp RecordListResp, err error)       // 记录列表
	RecordAdd(req RecordAddReq) (resp RecordAddResp, err error)          // 记录新增
	RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) // 记录修改
	RecordDelete(req RecordDeleteReq) (err error)                        // 记录删除
	RecordEnable(req RecordEnableReq) (err error)                        // 记录启用
	RecordDisable(req RecordDisableReq) (err error)                      // 记录暂停
}
