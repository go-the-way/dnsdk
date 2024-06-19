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
	"errors"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"

	dnspodErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

var dnspodLineDef = LineListRespLine{"0", "默认"}

func DnspodApi(client *dnspod.Client) Api { return &dnspodApi{client} }

type dnspodApi struct{ *dnspod.Client }

func (a *dnspodApi) LineList() (resp LineListResp) {
	lines := []LineListRespLine{
		dnspodLineDef,
		{"10=1", "电信"},
		{"10=0", "联通"},
		{"10=3", "移动"},
		{"3=0", "境外"},
	}
	return LineListResp{lines}
}

func (a *dnspodApi) LineDefault() (resp LineListRespLine) { return dnspodLineDef }

func (a *dnspodApi) DomainList(req DomainListReq) (resp DomainListResp, err error) {
	req0 := dnspod.NewDescribeDomainListRequest()
	req0.Keyword = tea.String(req.Domain)
	req0.Offset = tea.Int64(int64((req.Page - 1) * req.Limit))
	req0.Limit = tea.Int64(int64(req.Limit))
	return resp.transformFromDnspod(a.DescribeDomainList(req0))
}

func (a *dnspodApi) DomainAdd(req DomainAddReq) (resp DomainAddResp, err error) {
	req0 := dnspod.NewCreateDomainRequest()
	req0.Domain = tea.String(req.Domain)
	return resp.transformFromDnspod(a.CreateDomain(req0))
}

func (a *dnspodApi) DomainDelete(req DomainDeleteReq) (err error) {
	req0 := dnspod.NewDeleteDomainRequest()
	req0.Domain = tea.String(req.Domain)
	_, err = a.DeleteDomain(req0)
	return
}

func (a *dnspodApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	req0 := dnspod.NewDescribeRecordListRequest()
	req0.Domain = tea.String(req.Domain)
	req0.DomainId = toUint64Ptr(req.DomainId)
	req0.Subdomain = tea.String(req.Record)
	req0.RecordType = tea.String(req.Type)
	req0.Keyword = tea.String(req.Value)
	req0.RecordLineId = tea.String(req.Line)
	req0.SortField = tea.String(req.Order)
	req0.SortType = tea.String(strings.ToUpper(req.Direction))
	req0.Offset = tea.Uint64(uint64((req.Page - 1) * req.Limit))
	req0.Limit = tea.Uint64(uint64(req.Limit))
	return resp.transformFromDnspod(a.DescribeRecordList(req0))
}

func (a *dnspodApi) RecordAdd(req RecordAddReq) (resp RecordAddResp, err error) {
	req0 := dnspod.NewCreateRecordRequest()
	req0.Domain = tea.String(req.Domain)
	req0.DomainId = toUint64Ptr(req.DomainId)
	req0.SubDomain = tea.String(req.Record)
	req0.RecordType = tea.String(req.Type)
	req0.Value = tea.String(req.Value)
	req0.RecordLineId = tea.String(req.Line)
	req0.RecordLine = tea.String("")
	req0.TTL = tea.Uint64(uint64(req.TTL))
	req0.MX = tea.Uint64(1)
	req0.Weight = tea.Uint64(uint64(req.Weight))
	req0.Remark = tea.String(req.Remark)
	return resp.transformFromDnspod(a.CreateRecord(req0))
}

func (a *dnspodApi) RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) {
	req0 := dnspod.NewModifyRecordRequest()
	req0.RecordId = toUint64Ptr(req.RecordId)
	req0.Domain = tea.String(req.Domain)
	req0.DomainId = toUint64Ptr(req.DomainId)
	req0.SubDomain = tea.String(req.Record)
	req0.RecordType = tea.String(req.Type)
	req0.Value = tea.String(req.Value)
	req0.RecordLineId = tea.String(req.Line)
	req0.RecordLine = tea.String("")
	req0.TTL = tea.Uint64(uint64(req.TTL))
	req0.MX = tea.Uint64(1)
	req0.Weight = tea.Uint64(uint64(req.Weight))
	req0.Remark = tea.String(req.Remark)
	return resp.transformFromDnspod(a.ModifyRecord(req0))
}

func (a *dnspodApi) RecordDelete(req RecordDeleteReq) (err error) {
	req0 := dnspod.NewDeleteRecordRequest()
	req0.RecordId = toUint64Ptr(req.RecordId)
	req0.DomainId = toUint64Ptr(req.DomainId)
	req0.Domain = tea.String("")
	_, err = a.DeleteRecord(req0)
	return
}

func (a *dnspodApi) recordStatus(domainId, recordId string, status string) (err error) {
	req0 := dnspod.NewModifyRecordStatusRequest()
	req0.RecordId = toUint64Ptr(recordId)
	req0.DomainId = toUint64Ptr(domainId)
	req0.Domain = tea.String("")
	req0.Status = tea.String(status)
	_, err = a.ModifyRecordStatus(req0)
	return
}

func (a *dnspodApi) RecordEnable(req RecordEnableReq) (err error) {
	return a.recordStatus(req.DomainId, req.RecordId, "ENABLE")
}

func (a *dnspodApi) RecordDisable(req RecordDisableReq) (err error) {
	return a.recordStatus(req.DomainId, req.RecordId, "DISABLE")
}

func (*DomainListResp) transformFromDnspod(a *dnspod.DescribeDomainListResponse, err0 error) (resp DomainListResp, err error) {
	if err = err0; err != nil {
		return
	}
	aa := a.Response
	if aa == nil {
		return
	}
	if dci := aa.DomainCountInfo; dci != nil {
		resp.Total = uint(tea.Uint64Value(dci.AllTotal))
	}
	var list []DomainListRespDomain
	if dls := aa.DomainList; dls != nil {
		for _, aaa := range dls {
			list = append(list, DomainListRespDomain{
				Id:          fmt.Sprintf("%d", tea.Uint64Value(aaa.DomainId)),
				Name:        tea.StringValue(aaa.Name),
				DnsServer:   dnsServer(aaa.EffectiveDNS),
				RecordCount: uint(tea.Uint64Value(aaa.RecordCount)),
				Remark:      tea.StringValue(aaa.Remark),
				CreateTime:  tea.StringValue(aaa.CreatedOn),
			})
		}
	}
	return
}

func (*DomainAddResp) transformFromDnspod(a *dnspod.CreateDomainResponse, err0 error) (resp DomainAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	aa := a.Response
	if aa == nil {
		return
	}
	resp.Id = fmt.Sprintf("%d", tea.Uint64Value(aa.DomainInfo.Id))
	resp.DnsServer = dnsServer(aa.DomainInfo.GradeNsList)
	return
}

func (*RecordListRespRecord) dnspodRecordTransform(a *dnspod.RecordListItem) (record RecordListRespRecord) {
	return RecordListRespRecord{
		Id:         fmt.Sprintf("%d", tea.Uint64Value(a.RecordId)),
		Record:     tea.StringValue(a.Name),
		Name:       "", // TODO:
		Type:       tea.StringValue(a.Type),
		Value:      tea.StringValue(a.Value),
		Line:       tea.StringValue(a.LineId),
		TTL:        uint(tea.Uint64Value(a.TTL)),
		MX:         uint16(tea.Uint64Value(a.MX)),
		Weight:     uint(tea.Uint64Value(a.Weight)),
		Remark:     tea.StringValue(a.Remark),
		Status:     strings.ToLower(tea.StringValue(a.Status)),
		CreateTime: tea.StringValue(a.UpdatedOn),
		UpdateTime: tea.StringValue(a.UpdatedOn),
	}
}

func (*RecordListRespRecord) dnspodRecordTransformAdd(a *dnspod.CreateRecordResponse) (record RecordListRespRecord) {
	aa := a.Response
	if aa == nil {
		return
	}
	return RecordListRespRecord{Id: fmt.Sprintf("%d", tea.Uint64Value(aa.RecordId))}
}

func (*RecordListRespRecord) dnspodRecordTransformUpdate(a *dnspod.ModifyRecordResponse) (record RecordListRespRecord) {
	aa := a.Response
	if aa == nil {
		return
	}
	return RecordListRespRecord{Id: fmt.Sprintf("%d", tea.Uint64Value(aa.RecordId))}
}

func (r *RecordListResp) transformFromDnspod(a *dnspod.DescribeRecordListResponse, err0 error) (resp RecordListResp, err error) {
	ignoreCodesMap := map[string]struct{}{
		"ResourceNotFound.NoDataOfRecord": {},
	}
	if err0 != nil {
		var sdkError *dnspodErrors.TencentCloudSDKError
		if errors.As(err0, &sdkError) {
			if _, ignored := ignoreCodesMap[sdkError.Code]; ignored {
				// ignored
			} else {
				err = err0
				return
			}
		} else {
			err = err0
			return
		}
	}
	aa := a.Response
	if aa == nil {
		return
	}
	if rci := aa.RecordCountInfo; rci != nil {
		resp.Total = uint(tea.Uint64Value(rci.TotalCount))
	}
	var list = make([]RecordListRespRecord, 0)
	if dr := aa.RecordList; dr != nil {
		for _, aaa := range dr {
			list = append(list, (&RecordListRespRecord{}).dnspodRecordTransform(aaa))
		}
	}
	resp.List = list
	return
}

func (r *RecordAddResp) transformFromDnspod(a *dnspod.CreateRecordResponse, err0 error) (resp RecordAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = (&RecordListRespRecord{}).dnspodRecordTransformAdd(a)
	return
}

func (r *RecordUpdateResp) transformFromDnspod(a *dnspod.ModifyRecordResponse, err0 error) (resp RecordUpdateResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = (&RecordListRespRecord{}).dnspodRecordTransformUpdate(a)
	return
}
