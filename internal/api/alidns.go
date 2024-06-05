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
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/rwscode/dnsdk/internal/pkg"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
)

var (
	alidnsLines = []LineListRespLine{
		{"0", "default"},
		{"10=1", "unicom"},
		{"10=0", "telecom"},
		{"10=3", "mobile"},
		{"10=2", "edu"},
		{"3=0", "oversea"},
		{"10=22", "btvn"},
		{"80=0", "search"},
		{"7=0", "internal"},
	}
	alidnsLineMap = toLineMap(alidnsLines)
)

func AlidnsApi(client *alidns.Client) Api { return &alidnsApi{client} }

type alidnsApi struct{ *alidns.Client }

func (a *alidnsApi) DomainGet(req DomainGetReq) (resp DomainGetResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (a *alidnsApi) LineList(req LineListReq) (resp LineListResp, err error) {
	return resp.transformFromAlidns(a.DescribeSupportLines(&alidns.DescribeSupportLinesRequest{DomainName: tea.String(req.Domain)}))
}

func (a *alidnsApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	return resp.transformFromAlidns(a.DescribeDomainRecords(&alidns.DescribeDomainRecordsRequest{
		Direction:    tea.String(strings.ToUpper(req.Direction)),
		DomainName:   tea.String(req.Domain),
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
		Line:       tea.String(req.LineId),
		Priority:   tea.Int64(int64(req.MX)),
		RR:         tea.String(req.Record),
		TTL:        tea.Int64(int64(req.TTL)),
		Type:       tea.String(req.Type),
		Value:      tea.String(req.Value),
	}))
}

func (a *alidnsApi) RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) {
	return resp.transformFromAlidns(a.UpdateDomainRecord(&alidns.UpdateDomainRecordRequest{
		Line:     tea.String(req.LineId),
		Priority: tea.Int64(int64(req.MX)),
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

func alidnsRecordTransform(a *alidns.DescribeDomainRecordsResponseBodyDomainRecordsRecord) (record RecordListRespRecord) {
	return RecordListRespRecord{
		Id:         tea.StringValue(a.RecordId),
		Record:     tea.StringValue(a.RR),
		Name:       tea.StringValue(a.RecordId),
		Type:       tea.StringValue(a.Type),
		LineId:     tea.StringValue(a.Line), // FIXME? "0" or "default"
		LineName:   alidnsLineMap[tea.StringValue(a.Line)],
		Value:      tea.StringValue(a.Value),
		TTL:        uint(tea.Int64Value(a.TTL)),
		MX:         uint16(tea.Int64Value(a.Priority)),
		Weight:     uint(tea.Int32Value(a.Weight)),
		Remark:     tea.StringValue(a.Remark),
		CreateTime: pkg.FormatTime(time.UnixMilli(tea.Int64Value(a.CreateTimestamp))),
		UpdateTime: pkg.FormatTime(time.UnixMilli(tea.Int64Value(a.UpdateTimestamp))),
	}
}

func alidnsRecordTransformAdd(a *alidns.AddDomainRecordResponse) (record RecordListRespRecord) {
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

func (r *LineListResp) transformFromAlidns(a *alidns.DescribeSupportLinesResponse, err0 error) (resp LineListResp, err error) {
	if err = err0; err != nil {
		return
	}
	by := a.Body
	if by == nil {
		return
	}
	var list []LineListRespLine
	if dr := by.RecordLines; dr != nil {
		for _, aa := range dr.RecordLine {
			list = append(list, LineListRespLine{tea.StringValue(aa.LineCode), tea.StringValue(aa.LineName)})
		}
	}
	resp.List = list
	return
}

func (r *RecordListResp) transformFromAlidns(a *alidns.DescribeDomainRecordsResponse, err0 error) (resp RecordListResp, err error) {
	if err = err0; err != nil {
		return
	}
	by := a.Body
	if by == nil {
		return
	}
	resp.Total = uint(tea.Int64Value(by.TotalCount))
	var list []RecordListRespRecord
	if dr := by.DomainRecords; dr != nil {
		for _, aa := range dr.Record {
			list = append(list, alidnsRecordTransform(aa))
		}
	}
	resp.List = list
	return
}

func (r *RecordAddResp) transformFromAlidns(a *alidns.AddDomainRecordResponse, err0 error) (resp RecordAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = alidnsRecordTransformAdd(a)
	return
}

func (r *RecordUpdateResp) transformFromAlidns(a *alidns.UpdateDomainRecordResponse, err0 error) (resp RecordUpdateResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = alidnsRecordTransformUpdate(a)
	return
}
