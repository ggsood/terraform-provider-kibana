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
