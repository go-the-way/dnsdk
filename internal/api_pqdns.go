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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var pqdnsLineDef = LineListRespLine{"9065", "默认"}

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
	req.Header.Set("User-Agent", "dnsdk (https://github.com/go-the-way/dnsdk)")
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
	if reflect.ValueOf(respT).IsValid() {
		return json.Unmarshal(buf, respT)
	}
	return
}

func (a *pqdnsApi) Ping() (ok bool) {
	req, _ := http.NewRequest(http.MethodGet, a.baseUrl, nil)
	resp, _ := (&http.Client{Timeout: time.Second * 5}).Do(req)
	return resp.StatusCode == http.StatusOK
}

func (a *pqdnsApi) LineList() (resp LineListResp) {
	lines := []LineListRespLine{
		pqdnsLineDef,
		{"4", "电信"},
		{"2971", "联通"},
		{"5643", "移动"},
		{"8542", "境外"},
	}
	return LineListResp{lines}
}

func (a *pqdnsApi) LineDefault() (resp LineListRespLine) { return pqdnsLineDef }

func (a *pqdnsApi) DomainList(req DomainListReq) (resp DomainListResp, err error) {
	var rsp pqdnsDomainListResp
	apiUrl := fmt.Sprintf("/api/ext/dns/domain?domain=%s&page=%d&limit=%d", req.Domain, req.Page, req.Limit)
	err = a.req(apiUrl, http.MethodGet, nil, &rsp)
	return
}

func (a *pqdnsApi) DomainAdd(req DomainAddReq) (resp DomainAddResp, err error) {
	apiUrl := "/api/ext/dns/domain"
	err = a.req(apiUrl, http.MethodPost, (&pqdnsDomainAddReq{}).transform(a.username, a.secretKey, req), &resp)
	if err != nil {
		return
	}
	domainListResp, err0 := a.DomainList(DomainListReq{Page: 1, Limit: 1, Domain: req.Domain})
	if err0 != nil {
		err = err0
		return
	}
	for _, do0 := range domainListResp.List {
		resp.Id = do0.Id
		resp.DnsServer = do0.DnsServer
		break
	}
	return
}

func (a *pqdnsApi) DomainDelete(req DomainDeleteReq) (err error) {
	apiUrl := "/api/ext/dns/domain"
	err = a.req(apiUrl, http.MethodDelete, (&pqdnsDomainDeleteReq{}).transform(a.username, a.secretKey, req), nil)
	return
}

func (a *pqdnsApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	var rsp pqdnsRecordListResp
	apiUrl := fmt.Sprintf("/api/ext/dns/record?host_record=%s&record_value=%s&line_id=%s&page=%d&limit=%d", req.Record, req.Value, req.Line, req.Page, req.Limit)
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
	var rsp pqdnsRecordListResp
	apiUrl := "/api/ext/dns/record"
	err = a.req(apiUrl, http.MethodPut, (&pqdnsRecordUpdateReq{}).transform(a.username, a.secretKey, req), &rsp)
	return
}

func (a *pqdnsApi) RecordDelete(req RecordDeleteReq) (err error) {
	var rsp pqdnsRecordListResp
	apiUrl := "/api/ext/dns/record"
	err = a.req(apiUrl, http.MethodDelete, (&pqdnsRecordDeleteReq{}).transform(a.username, a.secretKey, req), &rsp)
	return
}

func (a *pqdnsApi) RecordEnable(_ RecordEnableReq) (err error) { return ErrNotSupportedOperation }

func (a *pqdnsApi) RecordDisable(_ RecordDisableReq) (err error) { return ErrNotSupportedOperation }

type (
	pqdnsDomainListReq struct {
		Username  string `json:"username"`
		SecretKey string `json:"secret_key"`
		Domain    string `json:"domain"` // 域名
	}
	pqdnsDomainAddReq struct {
		Username  string `json:"username"`
		SecretKey string `json:"secret_key"`
		Domain    string `json:"domain"` // 域名
	}
	pqdnsDomainDeleteReq struct {
		Username  string `json:"username"`
		SecretKey string `json:"secret_key"`
		Ids       []uint `json:"ids"` // 域名id
	}
	pqdnsRecordAddReq struct {
		Username  string `json:"username"`
		SecretKey string `json:"secret_key"`
		DomainId  uint   `json:"domain_id"`
		Host      string `json:"host"`
		RecType   string `json:"rec_type"`
		RecValue  string `json:"rec_value"`
		LineId    uint   `json:"line_id"`
		MX        uint   `json:"mx"`
		Weight    uint   `json:"weight"`
		TTL       uint   `json:"ttl"`
	}
	pqdnsRecordUpdateReq struct {
		Username  string `json:"username"`
		SecretKey string `json:"secret_key"`
		RecordId  uint   `json:"record_id"`
		Host      string `json:"host"`
		RecType   string `json:"rec_type"`
		RecValue  string `json:"rec_value"`
		LineId    uint   `json:"line_id"`
		MX        uint   `json:"mx"`
		Weight    uint   `json:"weight"`
		TTL       uint   `json:"ttl"`
	}
	pqdnsRecordDeleteReq struct {
		Username  string `json:"username"`
		SecretKey string `json:"secret_key"`
		DomainId  uint   `json:"domain_id"`
		RecordIds []uint `json:"record_id"`
	}
)

type (
	pqdnsDomainListResp struct {
		Total uint                        `json:"total"`
		List  []pqdnsRecordListRespDomain `json:"list"`
	}
	pqdnsRecordListRespDomain struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		DnsServer   []string `json:"tip_ns_value"`
		RecordCount uint     `json:"record_count"`
		Remark      string   `json:"remark"`
		CreateTime  string   `json:"create_time"`
	}
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

func (a *pqdnsDomainAddReq) transform(username, secretKey string, req DomainAddReq) (resp pqdnsDomainAddReq) {
	return pqdnsDomainAddReq{username, secretKey, req.Domain}
}

func (a *pqdnsDomainDeleteReq) transform(username, secretKey string, req DomainDeleteReq) (resp pqdnsDomainDeleteReq) {
	return pqdnsDomainDeleteReq{username, secretKey, []uint{toUint(req.Domain)}}
}

func (a *pqdnsRecordListResp) transform() (resp RecordListResp) {
	var list []RecordListRespRecord
	for _, rc := range a.List {
		list = append(list, RecordListRespRecord{
			Id:         fmt.Sprintf("%d", rc.Id),
			Record:     rc.HostRecord,
			Name:       fmt.Sprintf("%s.%s", rc.HostRecord, rc.DomainName),
			Type:       rc.RecordType,
			Value:      rc.RecordValue,
			Line:       fmt.Sprintf("%d", rc.LineId),
			TTL:        uint(rc.TTL),
			MX:         uint16(rc.MX),
			Weight:     uint(rc.Weight),
			Remark:     "", // ignored
			CreateTime: formatTime(rc.CreateTime),
			UpdateTime: formatTime(rc.UpdateTime),
		})
	}
	return RecordListResp{Total: a.Total, List: list}
}
func (a *pqdnsRecordAddReq) transform(username, secretKey string, req RecordAddReq) *pqdnsRecordAddReq {
	return &pqdnsRecordAddReq{
		Username:  username,
		SecretKey: secretKey,
		DomainId:  toUint(req.DomainId),
		Host:      req.Record,
		RecType:   req.Type,
		RecValue:  req.Value,
		LineId:    toUint(req.Line),
		MX:        1,
		Weight:    req.Weight,
		TTL:       req.TTL,
	}
}

func (a *pqdnsRecordUpdateReq) transform(username, secretKey string, req RecordUpdateReq) *pqdnsRecordUpdateReq {
	return &pqdnsRecordUpdateReq{
		Username:  username,
		SecretKey: secretKey,
		RecordId:  toUint(req.RecordId),
		Host:      req.Record,
		RecType:   req.Type,
		RecValue:  req.Value,
		LineId:    toUint(req.Line),
		MX:        1,
		Weight:    req.Weight,
		TTL:       req.TTL,
	}
}

func (a *pqdnsRecordDeleteReq) transform(username, secretKey string, req RecordDeleteReq) *pqdnsRecordDeleteReq {
	return &pqdnsRecordDeleteReq{
		Username:  username,
		SecretKey: secretKey,
		DomainId:  toUint(req.DomainId),
		RecordIds: []uint{toUint(req.RecordId)},
	}
}
