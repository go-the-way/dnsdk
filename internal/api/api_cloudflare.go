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
	"context"
	"github.com/cloudflare/cloudflare-go"
)

type cloudflareApi struct{ cApi *cloudflare.API }

func CloudflareApi(cApi *cloudflare.API) Api { return &cloudflareApi{cApi} }

func (a *cloudflareApi) ctx() context.Context { return context.Background() }

func (a *cloudflareApi) rc(domainId string) *cloudflare.ResourceContainer {
	return &cloudflare.ResourceContainer{Identifier: domainId, Type: cloudflare.AccountType}
}

func (a *cloudflareApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	return resp.transformFromCloudflare(a.cApi.ListDNSRecords(
		context.Background(),
		a.rc(req.DomainId),
		cloudflare.ListDNSRecordsParams{
			Type:      req.Type,
			Name:      "", // TODO for record: www => www.example.com
			Content:   req.Value,
			Comment:   req.Remark,
			Order:     req.Order,
			Direction: cloudflare.ListDirection(req.Direction),
			ResultInfo: cloudflare.ResultInfo{
				Page:    int(req.Page),
				PerPage: int(req.Limit),
			},
		}),
	)
}

func (a *cloudflareApi) RecordAdd(req RecordAddReq) (resp RecordAddResp, err error) {
	return resp.transformFromCloudflare(a.cApi.CreateDNSRecord(
		context.Background(),
		a.rc(req.DomainId),
		cloudflare.CreateDNSRecordParams{
			Type:    req.Type,
			Name:    req.Name,
			Content: req.Value,
			TTL:     int(req.TTL),
			Comment: req.Record,
		},
	))
}

func (a *cloudflareApi) RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) {
	return resp.transformFromCloudflare(a.cApi.UpdateDNSRecord(
		context.Background(),
		a.rc(req.DomainId),
		cloudflare.UpdateDNSRecordParams{
			Type:    req.Type,
			Name:    req.Name,
			Content: req.Value,
			ID:      req.RecordId,
			TTL:     int(req.TTL),
			Comment: &req.Remark,
		},
	))
}

func (a *cloudflareApi) RecordDelete(req RecordDeleteReq) (err error) {
	return a.cApi.DeleteDNSRecord(context.Background(), a.rc(req.DomainId), req.RecordId)
}

func (a *cloudflareApi) RecordEnable(_ RecordEnableReq) (err error) {
	return ErrNotSupportedOperation
}

func (a *cloudflareApi) RecordDisable(_ RecordDisableReq) (err error) {
	return ErrNotSupportedOperation
}
