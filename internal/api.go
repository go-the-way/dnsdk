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
	"io"
	"net/http"
	"reflect"
	"time"
)

type (
	Api interface {
		DomainList(req DomainListReq) (resp DomainListResp, err error)             // 域名列表
		RecordTypeList(req RecordTypeListReq) (resp RecordTypeListResp, err error) // 记录类型列表
		RecordLineList(req RecordLineListReq) (resp RecordLineListResp, err error) // 记录线路列表
		RecordList(req RecordListReq) (resp RecordListResp, err error)             // 记录列表
		RecordAdd(req RecordAddReq) (err error)                                    // 记录新增
		RecordUpdate(req RecordUpdateReq) (err error)                              // 记录修改
		RecordDelete(req RecordDeleteReq) (err error)                              // 记录删除
		RecordStatusUpdate(req RecordStatusUpdateReq) (err error)                  // 记录状态修改
	}
)

func doReq[REQ, RESP any](method, reqUrl string, header map[string]string, req REQ) (resp RESP, err error) {
	var (
		reader io.Reader
		buf    []byte
	)
	if reflect.ValueOf(req).IsValid() {
		if buf, err = json.Marshal(req); err != nil {
			return
		}
		reader = bytes.NewBuffer(buf)
	}
	request, err := http.NewRequest(method, reqUrl, reader)
	request.Header = make(http.Header)
	if header != nil {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	buf, err = io.ReadAll(response.Body)
	err = json.Unmarshal(buf, &resp)
	return
}
