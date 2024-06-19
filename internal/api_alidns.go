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
	"strings"
	"time"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

var alidnsLineDef = LineListRespLine{"default", "默认"}

func AlidnsApi(client *alidns.Client) Api { return &alidnsApi{client} }

type alidnsApi struct{ *alidns.Client }

func (a *alidnsApi) LineList() (resp LineListResp) {
	lines := []LineListRespLine{
		alidnsLineDef,
		{"telecom", "电信"},
		{"unicom", "联通"},
		{"mobile", "移动"},
		{"oversea", "境外"},
	}
	return LineListResp{lines}
}

func (a *alidnsApi) LineDefault() (resp LineListRespLine) { return alidnsLineDef }

func (a *alidnsApi) DomainList(req DomainListReq) (resp DomainListResp, err error) {
	return resp.transformFromAlidns(a.DescribeDomains(&alidns.DescribeDomainsRequest{
		KeyWord:    tea.String(req.Domain),
		PageNumber: tea.Int64(int64(req.Page)),
		PageSize:   tea.Int64(int64(req.Limit)),
		SearchMode: tea.String("EXACT"),
		Starmark:   tea.Bool(false),
	}))
}

func (a *alidnsApi) DomainAdd(req DomainAddReq) (resp DomainAddResp, err error) {
	return resp.transformFromAlidns(a.AddDomain(&alidns.AddDomainRequest{DomainName: tea.String(req.Domain)}))
}

func (a *alidnsApi) DomainDelete(req DomainDeleteReq) (err error) {
	_, err = a.DeleteDomain(&alidns.DeleteDomainRequest{DomainName: tea.String(req.Domain)})
	return
}

func (a *alidnsApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	return resp.transformFromAlidns(a.DescribeDomainRecords(&alidns.DescribeDomainRecordsRequest{
		Direction:    tea.String(strings.ToUpper(req.Direction)),
		DomainName:   tea.String(req.Domain),
		Line:         tea.String(req.Line),
		OrderBy:      tea.String(req.Order),
		PageNumber:   tea.Int64(int64(req.Page)),
		PageSize:     tea.Int64(int64(req.Limit)),
		Type:         tea.String(req.Type),
		ValueKeyWord: tea.String(req.Value),
	}))
}

func (a *alidnsApi) RecordAdd(req RecordAddReq) (resp RecordAddResp, err error) {
	return resp.transformFromAlidns(a.AddDomainRecord(&alidns.AddDomainRecordRequest{
		DomainName: tea.String(req.Domain),
		Line:       tea.String(req.Line),
		Priority:   tea.Int64(1),
		RR:         tea.String(req.Record),
		TTL:        tea.Int64(int64(req.TTL)),
		Type:       tea.String(req.Type),
		Value:      tea.String(req.Value),
	}))
}

func (a *alidnsApi) RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) {
	return resp.transformFromAlidns(a.UpdateDomainRecord(&alidns.UpdateDomainRecordRequest{
		Line:     tea.String(req.Line),
		Priority: tea.Int64(1),
		RR:       tea.String(req.Record),
		RecordId: tea.String(req.RecordId),
		TTL:      tea.Int64(int64(req.TTL)),
		Type:     tea.String(req.Type),
		Value:    tea.String(req.Value),
	}))
}

func (a *alidnsApi) RecordDelete(req RecordDeleteReq) (err error) {
	_, err = a.DeleteDomainRecord(&alidns.DeleteDomainRecordRequest{RecordId: tea.String(req.RecordId)})
	return
}

func (a *alidnsApi) recordStatus(recordId string, status string) (err error) {
	_, err = a.SetDomainRecordStatus(&alidns.SetDomainRecordStatusRequest{RecordId: tea.String(recordId), Status: tea.String(status)})
	return
}

func (a *alidnsApi) RecordEnable(req RecordEnableReq) (err error) {
	return a.recordStatus(req.RecordId, "Enable")
}

func (a *alidnsApi) RecordDisable(req RecordDisableReq) (err error) {
	return a.recordStatus(req.RecordId, "Disable")
}

func (_ *DomainListRespDomain) transformFromAlidns(a *alidns.DescribeDomainsResponseBodyDomainsDomain) (domain DomainListRespDomain) {
	return DomainListRespDomain{
		Id:   tea.StringValue(a.DomainId),
		Name: tea.StringValue(a.DomainName),
		DnsServer: func() []string {
			if ds := a.DnsServers; ds != nil {
				return dnsServer(a.DnsServers.DnsServer)
			}
			return []string{}
		}(),
		RecordCount: uint(tea.Int64Value(a.RecordCount)),
		Remark:      tea.StringValue(a.Remark),
		CreateTime:  formatTime(time.UnixMilli(tea.Int64Value(a.CreateTimestamp))),
	}
}

func (_ *DomainAddResp) transformFromAlidns(a *alidns.AddDomainResponse, err0 error) (resp DomainAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	aa := a.Body
	if aa == nil {
		return
	}
	resp = DomainAddResp{
		Id: tea.StringValue(aa.DomainId),
		DnsServer: func() []string {
			if ds := aa.DnsServers; ds != nil {
				return dnsServer(ds.DnsServer)
			}
			return []string{}
		}(),
	}
	return
}

func (_ *RecordListRespRecord) transformFromAlidns(a *alidns.DescribeDomainRecordsResponseBodyDomainRecordsRecord) (record RecordListRespRecord) {
	return RecordListRespRecord{
		Id:         tea.StringValue(a.RecordId),
		Record:     tea.StringValue(a.RR),
		Name:       tea.StringValue(a.RecordId),
		Type:       tea.StringValue(a.Type),
		Value:      tea.StringValue(a.Value),
		Line:       tea.StringValue(a.Line),
		TTL:        uint(tea.Int64Value(a.TTL)),
		MX:         uint16(tea.Int64Value(a.Priority)),
		Weight:     uint(tea.Int32Value(a.Weight)),
		Remark:     tea.StringValue(a.Remark),
		Status:     strings.ToLower(tea.StringValue(a.Status)),
		CreateTime: formatTime(time.UnixMilli(tea.Int64Value(a.CreateTimestamp))),
		UpdateTime: formatTime(time.UnixMilli(tea.Int64Value(a.UpdateTimestamp))),
	}
}

func (_ *RecordListRespRecord) transformFromAlidnsAdd(a *alidns.AddDomainRecordResponse) (record RecordListRespRecord) {
	if a.Body == nil {
		return
	}
	return RecordListRespRecord{Id: tea.StringValue(a.Body.RecordId)}
}

func alidnsRecordTransformUpdate(a *alidns.UpdateDomainRecordResponse) (record RecordListRespRecord) {
	if a.Body == nil {
		return
	}
	return RecordListRespRecord{Id: tea.StringValue(a.Body.RecordId)}
}

func (_ *DomainListResp) transformFromAlidns(a *alidns.DescribeDomainsResponse, err0 error) (resp DomainListResp, err error) {
	if err = err0; err != nil {
		return
	}
	aa := a.Body
	if aa == nil {
		return
	}
	resp.Total = uint(tea.Int64Value(aa.TotalCount))
	var list []DomainListRespDomain
	if dr := aa.Domains; dr != nil {
		for _, dd := range dr.Domain {
			list = append(list, (&DomainListRespDomain{}).transformFromAlidns(dd))
		}
	}
	resp.List = list
	return
}

func (r *RecordListResp) transformFromAlidns(a *alidns.DescribeDomainRecordsResponse, err0 error) (resp RecordListResp, err error) {
	if err = err0; err != nil {
		return
	}
	aa := a.Body
	if aa == nil {
		return
	}
	resp.Total = uint(tea.Int64Value(aa.TotalCount))
	var list []RecordListRespRecord
	if dr := aa.DomainRecords; dr != nil {
		for _, rr := range dr.Record {
			list = append(list, (&RecordListRespRecord{}).transformFromAlidns(rr))
		}
	}
	resp.List = list
	return
}

func (_ *RecordAddResp) transformFromAlidns(a *alidns.AddDomainRecordResponse, err0 error) (resp RecordAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = resp.transformFromAlidnsAdd(a)
	return
}

func (r *RecordUpdateResp) transformFromAlidns(a *alidns.UpdateDomainRecordResponse, err0 error) (resp RecordUpdateResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = alidnsRecordTransformUpdate(a)
	return
}
