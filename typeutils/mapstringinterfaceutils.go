// helps with type coercion and merging maps
package typeutils

/*
   Copyright 2013 Am Laher

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

import (
	"fmt"
)

// coerce interface{} to slice of strings.
func ToStringSlice(v interface{}, k string) ([]string, error) {
	ret := []string{}
	switch typedV := v.(type) {
	case []interface{}:
		for _, i := range typedV {
			ret = append(ret, i.(string))
		}
		return ret, nil
	}
	return ret, fmt.Errorf("%s should be a json array, not a %T", k, v)
}

// coerce interface{} to string.
func ToString(v interface{}, k string) (string, error) {
	switch typedV := v.(type) {
	case string:
		return typedV, nil
	}
	return "", fmt.Errorf("%s should be a json string, not a %T", k, v)
}

// coerce interface{} to bool
func ToBool(v interface{}, k string) (bool, error) {
	switch typedV := v.(type) {
	case bool:
		return typedV, nil
	}
	return false, fmt.Errorf("%s should be a json boolean, not a %T", k, v)
}

// coerce interface{} to map[string]interface{}
func ToMap(v interface{}, k string) (map[string]interface{}, error) {
	switch typedV := v.(type) {
	case map[string]interface{}:
		return typedV, nil
	}
	return nil, fmt.Errorf("%s should be a json map, not a %T", k, v)
}

// coerce interface{} to map[string]map[string]interface{}
func ToMapStringMapStringInterface(v interface{}, k string) (map[string]map[string]interface{}, error) {
	switch typedV := v.(type) {
	case map[string]interface{}:
		ret := make(map[string]map[string]interface{})
		for k, v := range typedV {
			typedSubV, err := ToMap(v, k)
			ret[k] = typedSubV
			if err != nil {
				return nil, fmt.Errorf("%s should be a json map[string]map[string]interface{}, not a %T", k, v)
			}
		}
		return ret, nil
	}
	return nil, fmt.Errorf("%s should be a json map[string]map[string]interface{}, not a %T", k, v)
}

// merge nested maps (first argument takes priority)
// note that lists are replaced, not merged
func MergeMapsStringMapStringInterface(high, low map[string]map[string]interface{}) map[string]map[string]interface{} {
	if high == nil {
		return low
	}
	for key, lowValTyped := range low {
		if highValTyped, keyExists := high[key]; keyExists {
			// NOTE: go deeper for maps.
			// (Slices and other types should not go deeper)
			high[key] = MergeMaps(highValTyped, lowValTyped)
		} else {
			high[key] = lowValTyped
		}
	}
	return high
}

// merge possibly-nested maps (first argument takes priority)
// note that lists are replaced, not merged
func MergeMaps(high, low map[string]interface{}) map[string]interface{} {
	if high == nil {
		return low
	}
	for key, lowVal := range low {
		if highVal, keyExists := high[key]; keyExists {
			// NOTE: go deeper for maps.
			// (Slices and other types should not go deeper)
			switch highValTyped := highVal.(type) {
			case map[string]interface{}:
				switch lowValTyped := lowVal.(type) {
				case map[string]interface{}:
					high[key] = MergeMaps(highValTyped, lowValTyped)
				}
			}
		} else {
			high[key] = lowVal
		}
	}
	return high
}
