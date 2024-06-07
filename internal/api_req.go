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

type (
	DomainListReq struct {
		Page   uint   `form:"page"`   // 页码 => 1
		Limit  uint   `form:"limit"`  // 每页数量 => 10
		Domain string `form:"domain"` // 域名 => example.com
	}
	DomainAddReq struct {
		Domain string `json:"domain"` // 域名 => example.com
	}
	DomainDeleteReq struct {
		Domain   string `json:"domain"`    // 域名 => example.com
		DomainId string `form:"domain_id"` // 域名Id => xxxxxxxxxxxx
	}
	DomainEnableReq struct {
		Domain   string `json:"domain"`    // 域名 => example.com
		DomainId string `form:"domain_id"` // 域名Id => xxxxxxxxxxxx
	}
	DomainDisableReq struct {
		Domain   string `json:"domain"`    // 域名 => example.com
		DomainId string `form:"domain_id"` // 域名Id => xxxxxxxxxxxx
	}
	RecordListReq struct {
		Page      uint   `form:"page"`      // 页码 => 1
		Limit     uint   `form:"limit"`     // 每页数量 => 10
		DomainId  string `form:"domain_id"` // 域名Id => xxxxxxxxxxxx
		Domain    string `json:"domain"`    // 域名 => example.com
		Record    string `json:"record"`    // 主机记录 => www
		Type      string `form:"type"`      // 解析类型 => A
		Value     string `form:"value"`     // 记录值 => 1.1.1.1
		Remark    string `form:"remark"`    // 备注 => created by dnsdk
		Order     string `form:"order"`     // 排序 => type
		Direction string `form:"direction"` // 方向 => asc / desc
	}
	RecordAddReq struct {
		DomainId string `json:"domain_id"` // 域名Id => xxxxxxxxxxxx
		Domain   string `json:"domain"`    // 域名 => example.com
		Record   string `json:"record"`    // 主机记录 => www
		Type     string `json:"type"`      // 类型 => A
		Value    string `json:"value"`     // 记录值 => 1.1.1.1
		TTL      uint   `json:"ttl"`       // TTL => 60
		MX       uint   `json:"mx"`        // MX => only for mx type
		Weight   uint   `json:"weight"`    // 权重 => 100
		Remark   string `json:"remark"`    // 备注 => created by dnsdk
	}
	RecordUpdateReq struct {
		RecordId string `json:"record_id"` // 记录Id => xxxxxxxxxxxx
		DomainId string `json:"domain_id"` // 域名Id => xxxxxxxxxxxx
		Domain   string `json:"domain"`    // 域名 => example.com
		Record   string `json:"record"`    // 主机记录 => www
		Name     string `json:"name"`      // 名称 => www.example.com
		Type     string `json:"type"`      // 类型 => A
		Value    string `json:"value"`     // 记录值 => 1.1.1.1
		TTL      uint   `json:"ttl"`       // TTL => 60
		MX       uint   `json:"mx"`        // MX => only for mx type
		Weight   uint   `json:"weight"`    // 权重 => 100
		Remark   string `json:"remark"`    // 备注 => created by dnsdk
	}
	RecordDeleteReq struct {
		RecordId string `json:"record_id"` // 记录Id => xxxxxxxxxxxx
		DomainId string `json:"domain_id"` // 域名Id => xxxxxxxxxxxx
	}
	RecordEnableReq struct {
		RecordId string `json:"record_id"` // 记录Id => xxxxxxxxxxxx
	}
	RecordDisableReq RecordEnableReq
)
