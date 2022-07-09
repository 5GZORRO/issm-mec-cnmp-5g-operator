/*
Copyright 2021 IBM.

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

package model

type Properties map[string]string

//type Properties map[string]interface{}
//func (v *Properties) DeepCopy() *Properties {
//	if v == nil {
//		return nil
//	}
//	out := make(map[string]interface{})
//
//	for key, val := range *v {
//		switch v := val.(type) {
//		case int:
//			out[key] = v
//		case string:
//			out[key] = v
//		case bool:
//			out[key] = v
//		case float64:
//			out[key] = v
//		case []interface{}:
//			out[key] = v
//		case map[string]interface{}:
//			out[key] = copyMap(v)
//		default:
//			// TODO
//		}
//	}
//
//	p := Properties(out)
//	return &p
//}
//
//func copyMap(m map[string]interface{}) map[string]interface{} {
//	cp := make(map[string]interface{})
//	for k, v := range m {
//		vm, ok := v.(map[string]interface{})
//		if ok {
//			cp[k] = copyMap(vm)
//		} else {
//			cp[k] = v
//		}
//	}
//
//	return cp
//}
