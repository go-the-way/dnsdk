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

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/cloudflare/cloudflare-go"
)

var cloudflareLineDef = LineListRespLine{"0", "默认"}

func CloudflareApi(cApi *cloudflare.API) Api { return &cloudflareApi{cApi} }

type cloudflareApi struct{ *cloudflare.API }

func (a *cloudflareApi) ctx() context.Context { return context.Background() }

func (a *cloudflareApi) rc(domainId string) *cloudflare.ResourceContainer {
	return &cloudflare.ResourceContainer{Identifier: domainId, Type: cloudflare.AccountType}
}

func (a *cloudflareApi) toParams(str string) []string {
	if str == "" {
		return []string{}
	}
	return []string{str}
}

func (a *cloudflareApi) LineList() (resp LineListResp) {
	return LineListResp{[]LineListRespLine{cloudflareLineDef}}
}

func (a *cloudflareApi) LineDefault() (resp LineListRespLine) {
	return cloudflareLineDef
}

func (a *cloudflareApi) DomainList(req DomainListReq) (resp DomainListResp, err error) {
	return resp.transformFromCloudflare(a.ListZones(context.Background(), a.toParams(req.Domain)...))
}

func (a *cloudflareApi) DomainAdd(req DomainAddReq) (resp DomainAddResp, err error) {
	return resp.transformFromCloudflare(a.CreateZone(context.Background(), req.Domain, true, cloudflare.Account{}, ""))
}

func (a *cloudflareApi) DomainDelete(req DomainDeleteReq) (err error) {
	_, err = a.DeleteZone(context.Background(), req.DomainId)
	return
}

func (a *cloudflareApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	name := ""
	if req.Record != "" && req.Domain != "" {
		name = fmt.Sprintf("%s.%s", req.Record, req.Domain)
	}
	return resp.transformFromCloudflare(a.ListDNSRecords(
		context.Background(),
		a.rc(req.DomainId),
		cloudflare.ListDNSRecordsParams{
			Type:       req.Type,
			Name:       name,
			Content:    req.Value,
			Comment:    req.Remark,
			Order:      req.Order,
			Direction:  cloudflare.ListDirection(req.Direction),
			ResultInfo: cloudflare.ResultInfo{Page: int(req.Page), PerPage: int(req.Limit)},
		}),
	)
}

func (a *cloudflareApi) RecordAdd(req RecordAddReq) (resp RecordAddResp, err error) {
	return resp.transformFromCloudflare(a.CreateDNSRecord(
		context.Background(),
		a.rc(req.DomainId),
		cloudflare.CreateDNSRecordParams{
			Type:     req.Type,
			Name:     fmt.Sprintf("%s.%s", req.Record, req.Domain),
			Content:  req.Value,
			Priority: tea.Uint16(uint16(1)),
			TTL:      int(req.TTL),
			Comment:  req.Remark,
		},
	))
}

func (a *cloudflareApi) RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) {
	return resp.transformFromCloudflare(a.UpdateDNSRecord(
		context.Background(),
		a.rc(req.DomainId),
		cloudflare.UpdateDNSRecordParams{
			Type:     req.Type,
			Name:     fmt.Sprintf("%s.%s", req.Record, req.Domain),
			Content:  req.Value,
			ID:       req.RecordId,
			Priority: tea.Uint16(uint16(1)),
			TTL:      int(req.TTL),
			Comment:  &req.Remark,
		},
	))
}

func (a *cloudflareApi) RecordDelete(req RecordDeleteReq) (err error) {
	return a.DeleteDNSRecord(context.Background(), a.rc(req.DomainId), req.RecordId)
}

func (a *cloudflareApi) RecordEnable(_ RecordEnableReq) (err error) {
	return ErrNotSupportedOperation
}

func (a *cloudflareApi) RecordDisable(_ RecordDisableReq) (err error) {
	return ErrNotSupportedOperation
}

func (_ *DomainListResp) transformFromCloudflare(zones []cloudflare.Zone, err0 error) (resp DomainListResp, err error) {
	if err = err0; err != nil {
		return
	}
	var list []DomainListRespDomain
	for _, zone := range zones {
		list = append(list, (&DomainListRespDomain{}).transformFromCloudflare(zone))
	}
	resp.List = list
	return
}

func (_ *DomainListRespDomain) transformFromCloudflare(a cloudflare.Zone) (domain DomainListRespDomain) {
	return DomainListRespDomain{
		Id:         a.ID,
		Name:       a.Name,
		DnsServer:  a.NameServers,
		CreateTime: formatTime(a.CreatedOn),
	}
}

func (_ *DomainAddResp) transformFromCloudflare(a cloudflare.Zone, err0 error) (resp DomainAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp = DomainAddResp{a.ID, a.NameServers}
	return
}

func (_ *RecordListRespRecord) transformFromCloudflare(a cloudflare.DNSRecord) (record RecordListRespRecord) {
	return RecordListRespRecord{
		Id:         a.ID,
		Record:     strings.TrimSuffix(strings.ReplaceAll(a.Name, a.ZoneName, ""), "."),
		Name:       a.Name,
		Type:       a.Type,
		Value:      a.Content,
		Line:       cloudflareLineDef.Id,
		TTL:        uint(a.TTL),
		MX:         tea.Uint16Value(a.Priority),
		Remark:     a.Comment,
		CreateTime: formatTime(a.CreatedOn),
		UpdateTime: formatTime(a.ModifiedOn),
	}
}

func (_ *RecordListResp) transformFromCloudflare(records []cloudflare.DNSRecord, resultInfo *cloudflare.ResultInfo, err0 error) (resp RecordListResp, err error) {
	if err = err0; err != nil {
		return
	}
	if resultInfo != nil {
		resp.Total = uint(resultInfo.Total)
	}
	var list []RecordListRespRecord
	if len(records) > 0 {
		for _, a := range records {
			list = append(list, (&RecordListRespRecord{}).transformFromCloudflare(a))
		}
	}
	resp.List = list
	return
}

func (r *RecordAddResp) transformFromCloudflare(a cloudflare.DNSRecord, err0 error) (resp RecordAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = (&RecordListRespRecord{}).transformFromCloudflare(a)
	return
}

func (r *RecordUpdateResp) transformFromCloudflare(a cloudflare.DNSRecord, err0 error) (resp RecordUpdateResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = (&RecordListRespRecord{}).transformFromCloudflare(a)
	return
}
