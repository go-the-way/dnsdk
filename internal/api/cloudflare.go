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
	"fmt"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rwscode/dnsdk/internal/pkg"
)

const (
	cloudflareLineDef = iota
	cloudflareLinePx
)

const (
	cloudflareLineDefDesc = "default"
	cloudflareLinePxDesc  = "proxy"
)

const (
	cloudflarePxTrue  = true
	cloudflarePxFalse = false
)

var (
	cloudflareLines = []LineListRespLine{
		{i2s(cloudflareLineDef), cloudflareLineDefDesc},
		{i2s(cloudflareLinePx), cloudflareLinePxDesc},
	}
	cloudflareLineMap         = map[byte]string{cloudflareLineDef: cloudflareLineDefDesc, cloudflareLinePx: cloudflareLinePxDesc}
	cloudflareProxiedMap      = map[bool]byte{cloudflarePxFalse: cloudflareLineDef, cloudflarePxTrue: cloudflareLinePx}
	cloudflareLine2ProxiedMap = map[byte]bool{cloudflareLineDef: cloudflarePxFalse, cloudflareLinePx: cloudflarePxTrue}
)

func i2s(i int) string { return fmt.Sprintf("%d", i) }

func CloudflareApi(cApi *cloudflare.API) Api { return &cloudflareApi{cApi} }

type cloudflareApi struct{ *cloudflare.API }

func (a *cloudflareApi) ctx() context.Context { return context.Background() }

func (a *cloudflareApi) rc(domainId string) *cloudflare.ResourceContainer {
	return &cloudflare.ResourceContainer{Identifier: domainId, Type: cloudflare.AccountType}
}

func (a *cloudflareApi) px(lineId string) bool {
	px := false
	if lineId != "" {
		i, err0 := strconv.Atoi(lineId)
		if err0 == nil {
			px = cloudflareLine2ProxiedMap[byte(i)]
		}
	}
	return px
}

func (a *cloudflareApi) DomainGet(req DomainGetReq) (resp DomainGetResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (a *cloudflareApi) LineList(_ LineListReq) (resp LineListResp, err error) {
	resp.List = cloudflareLines
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
			Priority: tea.Uint16(uint16(req.MX)),
			TTL:      int(req.TTL),
			Proxied:  tea.Bool(a.px(req.LineId)),
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
			Name:     req.Name,
			Content:  req.Value,
			ID:       req.RecordId,
			Priority: tea.Uint16(uint16(req.MX)),
			TTL:      int(req.TTL),
			Proxied:  tea.Bool(a.px(req.LineId)),
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

func cloudflareDnsRecordTransform(a cloudflare.DNSRecord) (record RecordListRespRecord) {
	return RecordListRespRecord{
		Id:         a.ID,
		Record:     strings.TrimRight(a.Name, a.ZoneName),
		Name:       a.Name,
		Type:       a.Type,
		LineId:     fmt.Sprintf("%d", cloudflareProxiedMap[tea.BoolValue(a.Proxied)]),
		LineName:   cloudflareLineMap[cloudflareProxiedMap[tea.BoolValue(a.Proxied)]],
		Value:      a.Content,
		TTL:        uint(a.TTL),
		MX:         tea.Uint16Value(a.Priority),
		Weight:     0, // ignored
		Remark:     a.Comment,
		CreateTime: pkg.FormatTime(a.CreatedOn),
		UpdateTime: pkg.FormatTime(a.ModifiedOn),
	}
}

func (r *RecordListResp) transformFromCloudflare(records []cloudflare.DNSRecord, resultInfo *cloudflare.ResultInfo, err0 error) (resp RecordListResp, err error) {
	if err = err0; err != nil {
		return
	}
	if resultInfo != nil {
		resp.Total = uint(resultInfo.Total)
	}
	var list []RecordListRespRecord
	if len(records) > 0 {
		for _, a := range records {
			list = append(list, cloudflareDnsRecordTransform(a))
		}
	}
	resp.List = list
	return
}

func (r *RecordAddResp) transformFromCloudflare(a cloudflare.DNSRecord, err0 error) (resp RecordAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = cloudflareDnsRecordTransform(a)
	return
}

func (r *RecordUpdateResp) transformFromCloudflare(a cloudflare.DNSRecord, err0 error) (resp RecordUpdateResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = cloudflareDnsRecordTransform(a)
	return
}
