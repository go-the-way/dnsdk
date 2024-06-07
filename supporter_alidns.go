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

func AlidnsSupporter[T any](supportFunc SupportFunc[T, *AlidnsSupportOpts]) supporter[T, *AlidnsSupportOpts] {
	return &defaultSupporter[T, *AlidnsSupportOpts]{ApiType: ApiTypeAlidns, SupportFunc: supportFunc}
}

const (
	alidnsEndpoint           = "alidns.aliyuncs.com"
	alidnsSignatureVersion   = "1.0"
	alidnsSignatureAlgorithm = "HMAC-MD5"
)

type AlidnsSupportOpts struct {
	accessKeyId        string
	accessKeySecret    string
	endpoint           string
	signatureVersion   string
	signatureAlgorithm string
}

func NewAlidnsSupportOpts(accessKeyId string, accessKeySecret string) *AlidnsSupportOpts {
	return NewAlidnsSupportOptsWithParams(accessKeyId, accessKeySecret, alidnsEndpoint, alidnsSignatureVersion, alidnsSignatureAlgorithm)
}

func NewAlidnsSupportOptsWithParams(accessKeyId string, accessKeySecret string, endpoint string, signatureVersion string, signatureAlgorithm string) *AlidnsSupportOpts {
	return &AlidnsSupportOpts{accessKeyId, accessKeySecret, endpoint, signatureVersion, signatureAlgorithm}
}
