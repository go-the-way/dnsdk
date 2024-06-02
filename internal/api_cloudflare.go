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
	"fmt"
	"net/http"
)

const (
	cloudflareApiBaseUrl     = "https://api.cloudflare.com/client/v4"
	cloudflareApiContentType = "application/json"
)

type cloudflareApi struct{ baseUrl, email, apiKey string }

func CloudflareApi(email, apiKey string) *cloudflareApi {
	return &cloudflareApi{cloudflareApiBaseUrl, email, apiKey}
}

func (c *cloudflareApi) joinReqUrl(url string) string { return fmt.Sprintf("%s%s", c.baseUrl, url) }

func (c *cloudflareApi) getHeaders() map[string]string {
	return map[string]string{"Content-Type": cloudflareApiContentType, "X-Auth-Email": c.email, "X-Auth-Key": c.apiKey}
}

func (c *cloudflareApi) DomainList(req DomainListReq) (resp DomainListResp, err error) {
	reqUrl := c.joinReqUrl(fmt.Sprintf("/zones?%s", req.url()))
	reqResp, err := doReq[any, cloudflareApiZonesResp](http.MethodGet, reqUrl, c.getHeaders(), nil)
	resp = reqResp.transform()
	return
}

func (c *cloudflareApi) RecordTypeList(req RecordTypeListReq) (resp RecordTypeListResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c *cloudflareApi) RecordLineList(req RecordLineListReq) (resp RecordLineListResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c *cloudflareApi) RecordList(req RecordListReq) (resp RecordListResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c *cloudflareApi) RecordAdd(req RecordAddReq) (err error) {
	// TODO implement me
	panic("implement me")
}

func (c *cloudflareApi) RecordUpdate(req RecordUpdateReq) (err error) {
	// TODO implement me
	panic("implement me")
}

func (c *cloudflareApi) RecordDelete(req RecordDeleteReq) (err error) {
	// TODO implement me
	panic("implement me")
}

func (c *cloudflareApi) RecordStatusUpdate(req RecordStatusUpdateReq) (err error) {
	// TODO implement me
	panic("implement me")
}
