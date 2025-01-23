package configs

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
)

func override(a, b map[string]interface{}) {
	for k, v := range b {
		if av, ok := a[k]; ok {
			if reflect.TypeOf(v) == reflect.TypeOf(av) {
				switch v.(type) {
				case map[string]interface{}:
					override(av.(map[string]interface{}), v.(map[string]interface{}))
				case []interface{}:
					a[k] = v
				default:
					a[k] = v
				}
			} else {
				a[k] = v
			}
		} else {
			a[k] = v
		}
	}
}

func hashOf(data map[string]interface{}) string {
	bs, _ := json.MarshalIndent(data, "", "")
	return digestOf(bs)
}

func digestOf(bs []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bs))
}
