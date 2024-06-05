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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/rwscode/dnsdk/internal/pkg"
)

func PqdnsApi(baseUrl, username, secretKey string) Api {
	return &pqdnsApi{baseUrl, username, secretKey}
}

type pqdnsApi struct{ baseUrl, username, secretKey string }

func (a *pqdnsApi) getAuthUrl() string {
	return fmt.Sprintf("user_name=%s&secret_key=%s", a.username, a.secretKey)
}

func (a *pqdnsApi) req(apiUrl, apiMethod string, reqT, respT any) (err error) {
	var prefix string
	if strings.Contains(apiUrl, "?") {
		prefix = "&"
	} else {
		prefix = "?"
	}
	apiUrl += prefix + a.getAuthUrl()
	reqUrl := fmt.Sprintf("%s%s", a.baseUrl, apiUrl)
	client := &http.Client{Timeout: time.Second * 10}
	var reader io.Reader
	if reflect.ValueOf(reqT).IsValid() {
		buf, err0 := json.Marshal(reqT)
		if err0 != nil {
			return err0
		}
		reader = bytes.NewBuffer(buf)
	}
	req, _ := http.NewRequest(apiMethod, reqUrl, reader)
	req.Header = make(http.Header)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dnsdk (https://github.com/rwscode/dnsdk)")
	resp, err0 := client.Do(req)
	if err0 != nil {
		err = err0
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}
	buf, err0 := io.ReadAll(resp.Body)
	if err0 != nil {
		err = err0
		return
	}
	return json.Unmarshal(buf, respT)
}

func (a *pqdnsApi) DomainGet(req DomainGetReq) (resp DomainGetResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (a *pqdnsApi) LineList(req LineListReq) (resp LineListResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (a *pqdnsApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	var rsp pqdnsRecordListResp
	apiUrl := fmt.Sprintf("/api/ext/dns/record?host_record=%s&record_value=%s", req.Record, req.Value)
	err = a.req(apiUrl, http.MethodGet, nil, &rsp)
	resp = rsp.transform()
	return
}

func (a *pqdnsApi) RecordAdd(req RecordAddReq) (resp RecordAddResp, err error) {
	var rsp pqdnsRecordListResp
	apiUrl := "/api/ext/dns/record"
	err = a.req(apiUrl, http.MethodPost, (&pqdnsRecordAddReq{}).transform(a.username, a.secretKey, req), &rsp)
	return
}

func (a *pqdnsApi) RecordUpdate(req RecordUpdateReq) (resp RecordUpdateResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (a *pqdnsApi) RecordDelete(req RecordDeleteReq) (err error) {
	// TODO implement me
	panic("implement me")
}

func (a *pqdnsApi) RecordEnable(req RecordEnableReq) (err error) {
	// TODO implement me
	panic("implement me")
}

func (a *pqdnsApi) RecordDisable(req RecordDisableReq) (err error) {
	// TODO implement me
	panic("implement me")
}

type (
	pqdnsRecordListResp struct {
		Total uint                        `json:"total"`
		List  []pqdnsRecordListRespRecord `json:"list"`
	}
	pqdnsRecordListRespRecord struct {
		Id           uint32    `json:"id"`
		DomainId     uint32    `json:"domain_id"`
		HostRecord   string    `json:"host_record"`
		RecordType   string    `json:"record_type"`
		RecordValue  string    `json:"record_value"`
		IsCustomLine byte      `json:"is_custom_line"`
		LineId       uint32    `json:"line_id"`
		Weight       uint32    `json:"weight"`
		MX           int       `json:"mx"`
		TTL          uint32    `json:"ttl"`
		Status       byte      `json:"status"` // 0表示参与解析，1表示不参与解析
		TCModel      bool      `json:"tc_model"`
		CreateTime   time.Time `json:"create_time"`
		UpdateTime   time.Time `json:"update_time"`
		DomainName   string    `json:"domain_name"`
		LineName     string    `json:"line_name"`
	}
)

func (a *pqdnsRecordListResp) transform() (resp RecordListResp) {
	var list []RecordListRespRecord
	for _, rc := range a.List {
		list = append(list, RecordListRespRecord{
			Id:         fmt.Sprintf("%d", rc.Id),
			Record:     rc.HostRecord,
			Name:       fmt.Sprintf("%s.%s", rc.HostRecord, rc.DomainName),
			Type:       rc.RecordType,
			LineId:     fmt.Sprintf("%d", rc.LineId),
			LineName:   rc.LineName,
			Value:      rc.RecordValue,
			TTL:        uint(rc.TTL),
			MX:         uint16(rc.MX),
			Weight:     uint(rc.Weight),
			Remark:     "", // ignored
			CreateTime: pkg.FormatTime(rc.CreateTime),
			UpdateTime: pkg.FormatTime(rc.UpdateTime),
		})
	}
	return RecordListResp{Total: a.Total, List: list}
}

type pqdnsRecordAddReq struct {
	Username  string `json:"username"`
	SecretKey string `json:"secret_key"`

	DomainId uint `json:"domain_id"`

	Host     string `json:"host"`
	RecType  string `json:"rec_type"`
	RecValue string `json:"rec_value"`
	LineId   uint   `json:"line_id"`
	MX       uint   `json:"mx"`
	Weight   uint   `json:"weight"`
	TTL      uint   `json:"ttl"`
}

func (a *pqdnsRecordAddReq) transform(username, secretKey string, req RecordAddReq) *pqdnsRecordAddReq {
	return &pqdnsRecordAddReq{
		Username:  username,
		SecretKey: secretKey,
		DomainId:  toUint(req.DomainId),
		Host:      req.Record,
		RecType:   req.Type,
		RecValue:  req.Value,
		LineId:    toUint(req.LineId),
		MX:        req.MX,
		Weight:    req.Weight,
		TTL:       req.TTL,
	}
}
