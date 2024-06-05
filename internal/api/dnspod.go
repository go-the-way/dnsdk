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
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"

	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func DnspodApi(client *dnspod.Client) Api { return &dnspodApi{client} }

type dnspodApi struct{ *dnspod.Client }

func (a *dnspodApi) DomainGet(req DomainGetReq) (resp DomainGetResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (a *dnspodApi) LineList(req LineListReq) (resp LineListResp, err error) {
	req0 := dnspod.NewDescribeRecordLineListRequest()
	req0.DomainId = toUint64Ptr(req.DomainId)
	req0.Domain = tea.String(req.Domain)
	// TODO? 域名等级。
	// + 旧套餐：D_FREE、D_PLUS、D_EXTRA、D_EXPERT、D_ULTRA 分别对应免费套餐、个人豪华、企业1、企业2、企业3。
	// + 新套餐：DP_FREE、DP_PLUS、DP_EXTRA、DP_EXPERT、DP_ULTRA 分别对应新免费、个人专业版、企业创业版、企业标准版、企业旗舰版。
	req0.DomainGrade = tea.String("DP_FREE")
	_, err = a.DescribeRecordLineList(req0)
	return resp.transformFromDnspod(a.DescribeRecordLineList(req0))
}

func (a *dnspodApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	req0 := dnspod.NewDescribeRecordListRequest()
	req0.Domain = tea.String(req.Domain)
	req0.DomainId = toUint64Ptr(req.DomainId)
	req0.Subdomain = tea.String(req.Record)
	req0.RecordType = tea.String(req.Type)
	req0.Keyword = tea.String(req.Value)
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
	req0.TTL = tea.Uint64(uint64(req.TTL))
	req0.MX = tea.Uint64(uint64(req.MX))
	req0.Weight = tea.Uint64(uint64(req.Weight))
	req0.Remark = tea.String(req.Remark)
	return resp.transformFromDnspod(a.CreateRecord(req0))
}

func (a *dnspodApi) RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) {
	req0 := dnspod.NewModifyRecordRequest()
	req0.Domain = tea.String(req.Domain)
	req0.DomainId = toUint64Ptr(req.DomainId)
	req0.SubDomain = tea.String(req.Record)
	req0.RecordType = tea.String(req.Type)
	req0.Value = tea.String(req.Value)
	req0.TTL = tea.Uint64(uint64(req.TTL))
	req0.MX = tea.Uint64(uint64(req.MX))
	req0.Weight = tea.Uint64(uint64(req.Weight))
	req0.Remark = tea.String(req.Remark)
	return resp.transformFromDnspod(a.ModifyRecord(req0))
}

func (a *dnspodApi) RecordDelete(req RecordDeleteReq) (err error) {
	req0 := dnspod.NewDeleteRecordRequest()
	req0.RecordId = toUint64Ptr(req.RecordId)
	req0.DomainId = toUint64Ptr(req.DomainId)
	_, err = a.DeleteRecord(req0)
	return
}

func (a *dnspodApi) recordStatus(recordId string, status string) (err error) {
	req0 := dnspod.NewModifyRecordStatusRequest()
	req0.RecordId = toUint64Ptr(recordId)
	req0.Status = tea.String(status)
	_, err = a.ModifyRecordStatus(req0)
	return
}

func (a *dnspodApi) RecordEnable(req RecordEnableReq) (err error) {
	return a.recordStatus(req.RecordId, "ENABLE")
}

func (a *dnspodApi) RecordDisable(req RecordDisableReq) (err error) {
	return a.recordStatus(req.RecordId, "DISABLE")
}

func dnspodRecordTransform(a *dnspod.RecordListItem) (record RecordListRespRecord) {
	return RecordListRespRecord{
		Id:         fmt.Sprintf("%d", tea.Uint64Value(a.RecordId)),
		Record:     "", // TODO: replace from name
		Name:       tea.StringValue(a.Name),
		Type:       tea.StringValue(a.Type),
		LineId:     "", // TODO:线路id
		LineName:   "", // TODO:线路名称
		Value:      tea.StringValue(a.Value),
		TTL:        uint(tea.Uint64Value(a.TTL)),
		MX:         uint16(tea.Uint64Value(a.MX)),
		Weight:     uint(tea.Uint64Value(a.Weight)),
		Remark:     tea.StringValue(a.Remark),
		CreateTime: tea.StringValue(a.UpdatedOn),
		UpdateTime: tea.StringValue(a.UpdatedOn),
	}
}

func dnspodRecordTransformAdd(a *dnspod.CreateRecordResponse) (record RecordListRespRecord) {
	by := a.Response
	if by == nil {
		return
	}
	return RecordListRespRecord{Id: fmt.Sprintf("%d", tea.Uint64Value(by.RecordId))}
}

func dnspodRecordTransformUpdate(a *dnspod.ModifyRecordResponse) (record RecordListRespRecord) {
	by := a.Response
	if by == nil {
		return
	}
	return RecordListRespRecord{Id: fmt.Sprintf("%d", tea.Uint64Value(by.RecordId))}
}

func (r *RecordListResp) transformFromDnspod(a *dnspod.DescribeRecordListResponse, err0 error) (resp RecordListResp, err error) {
	if err = err0; err != nil {
		return
	}
	by := a.Response
	if by == nil {
		return
	}
	if rci := by.RecordCountInfo; rci != nil {
		resp.Total = uint(tea.Uint64Value(rci.TotalCount))
	}
	var list = make([]RecordListRespRecord, 0)
	if dr := by.RecordList; dr != nil {
		for _, aa := range dr {
			list = append(list, dnspodRecordTransform(aa))
		}
	}
	resp.List = list
	return
}

func (r *RecordAddResp) transformFromDnspod(a *dnspod.CreateRecordResponse, err0 error) (resp RecordAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = dnspodRecordTransformAdd(a)
	return
}

func (r *RecordUpdateResp) transformFromDnspod(a *dnspod.ModifyRecordResponse, err0 error) (resp RecordUpdateResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = dnspodRecordTransformUpdate(a)
	return
}

func (r *LineListResp) transformFromDnspod(a *dnspod.DescribeRecordLineListResponse, err0 error) (resp LineListResp, err error) {
	if err = err0; err != nil {
		return
	}
	var list = make([]LineListRespLine, 0)
	if resp0 := a.Response; resp0 != nil {
		if lls := resp0.LineList; lls != nil {
			for _, line := range lls {
				list = append(list, LineListRespLine{tea.StringValue(line.LineId), tea.StringValue(line.Name)})
			}
		}
	}
	resp.List = list
	return
}
