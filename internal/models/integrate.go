/*
Copyright 2021 The AtomCI Group Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import (
	"encoding/base64"
	"github.com/go-atomci/atomci/utils"
)

// IntegrateSetting the Basic Data of stages based on commpany
type IntegrateSetting struct {
	Addons
	Name        string `orm:"column(name);size(64)" json:"name"`
	Type        string `orm:"column(type);size(64)" json:"type"`
	Config      string `orm:"column(config);type(text)" json:"config"`
	Description string `orm:"column(description);size(256)" json:"description"`
	Creator     string `orm:"column(creator);size(64)" json:"creator"`
}

// TableName ...
func (t *IntegrateSetting) TableName() string {
	return "sys_integrate_setting"
}

func (t *IntegrateSetting) CryptoConfig(raw string) {
	t.Config = t.crypto(raw)
}

func (t *IntegrateSetting) DecryptConfig() string {
	return t.decrypt()
}

func (t *IntegrateSetting) crypto(raw string) string {
	plainText := []byte(raw)
	return base64.StdEncoding.EncodeToString(utils.AesEny(plainText))
}

func (t *IntegrateSetting) decrypt() string {
	cfg, _ := base64.StdEncoding.DecodeString(t.Config)
	return string(utils.AesEny(cfg))
}
