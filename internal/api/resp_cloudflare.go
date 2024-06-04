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

	"github.com/cloudflare/cloudflare-go"
	"github.com/rwscode/dnsdk/internal/pkg"
)

/*
	{
	 "id": "xxxxxxxxxxxxxxxx",
	 "zone_id": "xxxxxxxxxxxxxx",
	 "zone_name": "example.com",
	 "name": "www.example.com,
	 "type": "A",
	 "content": x.y.z.w",
	 "proxiable": true,
	 "proxied": false,
	 "ttl": 60,
	 "locked": false,
	 "meta": {
	   "auto_added": false,
	   "managed_by_apps": false,
	   "managed_by_argo_tunnel": false
	 },
	 "comment": null,
	 "tags": [],
	 "created_on": "2022-09-27T08:09:25.556108Z",
	 "modified_on": "2022-11-02T08:47:13.410559Z"
	}
*/
func cloudflareDnsRecordTransform(a cloudflare.DNSRecord) (record RecordListRespRecord) {
	return RecordListRespRecord{
		Id:         a.ID,
		Record:     strings.TrimRight(a.Name, a.ZoneName),
		Name:       a.Name,
		Type:       a.Type,
		LineId:     "", // TODO
		LineName:   "", // TODO
		Value:      a.Content,
		TTL:        uint(a.TTL),
		MX:         "", // TODO
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

func (r *RecordAddResp) transformFromCloudflare(record cloudflare.DNSRecord, err0 error) (resp RecordAddResp, err error) {
	if err = err0; err != nil {
		return
	}
	resp.RecordListRespRecord = cloudflareDnsRecordTransform(record)
	return
}
