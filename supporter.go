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

const (
	ApiTypeAlidns     ApiType = "alidns"
	ApiTypeCloudflare ApiType = "cloudflare"
	ApiTypeDnspod     ApiType = "dnspod"
	ApiTypePqdns      ApiType = "pqdns"
)

type (
	ApiType             string
	supporter[T, R any] interface {
		Type() (at ApiType)
		Support(t T) (r R)
	}
	SupportFunc[T, R any]      func(t T) (r R)
	defaultSupporter[T, R any] struct {
		ApiType
		SupportFunc[T, R]
	}
)

func (d *defaultSupporter[T, R]) Type() (at ApiType) { return d.ApiType }
func (d *defaultSupporter[T, R]) Support(t T) (r R)  { return d.SupportFunc(t) }
