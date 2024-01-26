// This file is part of CycloneDX Go
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
// Copyright (c) OWASP Foundation. All Rights Reserved.

package cyclonedx

import (
	"encoding/json"
)

func (sv SpecVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(sv.String())
}

func (sv *SpecVersion) UnmarshalJSON(bytes []byte) error {
	var v string
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}

	switch v {
	case SpecVersion1_0.String():
		*sv = SpecVersion1_0
	case SpecVersion1_1.String():
		*sv = SpecVersion1_1
	case SpecVersion1_2.String():
		*sv = SpecVersion1_2
	case SpecVersion1_3.String():
		*sv = SpecVersion1_3
	case SpecVersion1_4.String():
		*sv = SpecVersion1_4
	case SpecVersion1_5.String():
		*sv = SpecVersion1_5
	default:
		return ErrInvalidSpecVersion
	}

	return nil
}

var jsonSchemas = map[SpecVersion]string{
	SpecVersion1_0: "",
	SpecVersion1_1: "",
	SpecVersion1_2: "http://cyclonedx.org/schema/bom-1.2.schema.json",
	SpecVersion1_3: "http://cyclonedx.org/schema/bom-1.3.schema.json",
	SpecVersion1_4: "http://cyclonedx.org/schema/bom-1.4.schema.json",
	SpecVersion1_5: "http://cyclonedx.org/schema/bom-1.5.schema.json",
}
