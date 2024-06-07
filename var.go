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

package dnsdk

import (
	"errors"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rwscode/dnsdk/internal"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func GetSupportApi[T, R any](t T, supporter0 supporter[T, R]) (a Api, err error) {
	if supporter0 == nil {
		return nil, errors.New("nil supporter error")
	}
	switch at := supporter0.Type(); at {
	default:
		return nil, errors.New("not supported:" + string(at))
	case ApiTypeAlidns:
		spr, ok := supporter0.(supporter[T, *AlidnsSupportOpts])
		if !ok {
			return nil, errors.New("invalid alidns opts definition")
		}
		return newAlidnsApi(spr.Support(t))
	case ApiTypeCloudflare:
		spr, ok := supporter0.(supporter[T, *CloudflareSupportOpts])
		if !ok {
			return nil, errors.New("invalid cloudflare opts definition")
		}
		return newCloudflareApi(spr.Support(t))
	case ApiTypeDnspod:
		spr, ok := supporter0.(supporter[T, *DnspodSupportOpts])
		if !ok {
			return nil, errors.New("invalid dnspod opts definition")
		}
		return newDnspodApi(spr.Support(t))
	case ApiTypePqdns:
		spr, ok := supporter0.(supporter[T, *PqdnsSupportOpts])
		if !ok {
			return nil, errors.New("invalid pqdns opts definition")
		}
		return newPqdnsApi(spr.Support(t))
	}
	return
}

func newAlidnsApi(opts *AlidnsSupportOpts) (a Api, err error) {
	client, err0 := alidns.NewClient(&openapi.Config{
		AccessKeyId:        tea.String(opts.accessKeyId),
		AccessKeySecret:    tea.String(opts.accessKeySecret),
		Endpoint:           tea.String(opts.endpoint),
		SignatureVersion:   tea.String(opts.signatureVersion),
		SignatureAlgorithm: tea.String(opts.signatureAlgorithm),
	})
	if err = err0; err != nil {
		return
	}
	a = internal.AlidnsApi(client)
	return
}

func newCloudflareApi(opts *CloudflareSupportOpts) (a Api, err error) {
	cApi, err0 := cloudflare.New(opts.apiKey, opts.email)
	if err = err0; err != nil {
		return
	}
	a = internal.CloudflareApi(cApi)
	return
}

func newDnspodApi(opts *DnspodSupportOpts) (a Api, err error) {
	client, err0 := dnspod.NewClient(common.NewCredential(opts.secretId, opts.secretKey), "", profile.NewClientProfile())
	if err = err0; err != nil {
		return
	}
	a = internal.DnspodApi(client)
	return
}

func newPqdnsApi(opts *PqdnsSupportOpts) (a Api, err error) {
	a = internal.PqdnsApi(opts.baseUrl, opts.username, opts.secretKey)
	return
}
