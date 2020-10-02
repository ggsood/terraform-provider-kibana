// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package kb

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// suppressEquivalentJSON permit to compare json string
func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {
	var oldObj, newObj interface{}
	if err := json.Unmarshal([]byte(old), &oldObj); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newObj); err != nil {
		return false
	}
	return reflect.DeepEqual(oldObj, newObj)
}

// suppressEquivalentNDJSON permit to compare ndjson string
func suppressEquivalentNDJSON(k, old, new string, d *schema.ResourceData) bool {

	// NDJSON mean sthat each line correspond to JSON struct
	oldSlice := strings.Split(old, "\n")
	newSlice := strings.Split(new, "\n")
	oldObjSlice := make([]map[string]interface{}, len(oldSlice))
	newObjSlice := make([]map[string]interface{}, len(newSlice))
	if len(oldSlice) != len(newSlice) {
		return false
	}

	// Convert string line to JSON
	for i, oldJSON := range oldSlice {
		jsonObj := make(map[string]interface{})
		if err := json.Unmarshal([]byte(oldJSON), &jsonObj); err != nil {
			return false
		}

		delete(jsonObj, "version")
		delete(jsonObj, "updated_at")

		oldObjSlice[i] = jsonObj
	}
	for i, newJSON := range newSlice {
		jsonObj := make(map[string]interface{})
		if err := json.Unmarshal([]byte(newJSON), &jsonObj); err != nil {
			return false
		}
		delete(jsonObj, "version")
		delete(jsonObj, "updated_at")

		newObjSlice[i] = jsonObj
	}

	// Compare json obj
	for _, oldJSON := range oldObjSlice {
		isFound := false
		for _, newJSON := range newObjSlice {
			if oldJSON["id"].(string) == newJSON["id"].(string) {
				if reflect.DeepEqual(oldJSON, newJSON) == false {
					return false
				}
				isFound = true
				break
			}
		}

		if isFound == false {
			return false
		}
	}

	return true

}
