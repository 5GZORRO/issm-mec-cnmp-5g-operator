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
