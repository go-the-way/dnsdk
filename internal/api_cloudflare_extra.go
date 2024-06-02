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

func (c *cloudflareApiZonesResp) transform() (resp DomainListResp) {
	resp.Total = uint(c.ResultInfo.TotalCount)
	var list []DomainListRespDomain
	for _, rst := range c.Result {
		list = append(list, DomainListRespDomain{Id: rst.Id, Name: rst.Name})
	}
	resp.List = list
	return
}
