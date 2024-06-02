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

import "time"

type (
	cloudflareApiZonesResp struct {
		Result []struct {
			Id                  string      `json:"id"`
			Name                string      `json:"name"`
			Status              string      `json:"status"`
			Paused              bool        `json:"paused"`
			Type                string      `json:"type"`
			DevelopmentMode     int         `json:"development_mode"`
			NameServers         []string    `json:"name_servers"`
			OriginalNameServers interface{} `json:"original_name_servers"`
			OriginalRegistrar   interface{} `json:"original_registrar"`
			OriginalDnshost     interface{} `json:"original_dnshost"`
			ModifiedOn          time.Time   `json:"modified_on"`
			CreatedOn           time.Time   `json:"created_on"`
			ActivatedOn         time.Time   `json:"activated_on"`
			Meta                struct {
				Step                   int  `json:"step"`
				CustomCertificateQuota int  `json:"custom_certificate_quota"`
				PageRuleQuota          int  `json:"page_rule_quota"`
				PhishingDetected       bool `json:"phishing_detected"`
			} `json:"meta"`
			Owner struct {
				Id    interface{} `json:"id"`
				Type  string      `json:"type"`
				Email interface{} `json:"email"`
			} `json:"owner"`
			Account struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"account"`
			Tenant struct {
				Id   interface{} `json:"id"`
				Name interface{} `json:"name"`
			} `json:"tenant"`
			TenantUnit struct {
				Id interface{} `json:"id"`
			} `json:"tenant_unit"`
			Permissions []string `json:"permissions"`
			Plan        struct {
				Id                string `json:"id"`
				Name              string `json:"name"`
				Price             int    `json:"price"`
				Currency          string `json:"currency"`
				Frequency         string `json:"frequency"`
				IsSubscribed      bool   `json:"is_subscribed"`
				CanSubscribe      bool   `json:"can_subscribe"`
				LegacyId          string `json:"legacy_id"`
				LegacyDiscount    bool   `json:"legacy_discount"`
				ExternallyManaged bool   `json:"externally_managed"`
			} `json:"plan"`
		} `json:"result"`
		ResultInfo struct {
			Page       int `json:"page"`
			PerPage    int `json:"per_page"`
			TotalPages int `json:"total_pages"`
			Count      int `json:"count"`
			TotalCount int `json:"total_count"`
		} `json:"result_info"`
		Success  bool          `json:"success"`
		Errors   []interface{} `json:"errors"`
		Messages []interface{} `json:"messages"`
	}
)
