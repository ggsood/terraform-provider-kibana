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

import "encoding/json"

// optionalInterfaceJSON permit to convert string as json object
func optionalInterfaceJSON(input string) interface{} {
	if input == "" || input == "{}" {
		return nil
	}
	return json.RawMessage(input)

}

// convertArrayInterfaceToArrayString permit to convert an array of interface to an array of string
func convertArrayInterfaceToArrayString(raws []interface{}) []string {
	data := make([]string, len(raws))
	for i, raw := range raws {
		data[i] = raw.(string)
	}

	return data
}

// unused method
//func convertMapInterfaceToMapString(raws map[string]interface{}) map[string]string {
//	data := make(map[string]string)
//	for k, v := range raws {
//		data[k] = v.(string)
//	}
//
//	return data
//}
