package main

import (
	"encoding/json"
)

// use data2 to override data1
func Override(data1, data2 []byte) ([]byte, error) {
	// check json format
	var j1 interface{}
	if err := json.Unmarshal(data1, &j1); err != nil {
		return nil, err
	}
	var j2 interface{}
	if err := json.Unmarshal(data2, &j2); err != nil {
		return nil, err
	}

	out, err := json.Marshal(merge(j1, j2))
	if err != nil {
		return nil, err
	}

	return out, nil
}

func merge(oldJ, newJ interface{}) interface{} {
	switch newO := newJ.(type) {
	case map[string]interface{}:
		oldO, ok := oldJ.(map[string]interface{})
		if !ok {
			return newO
		}
		for k, oldV := range oldO {
			if newV, ok := newO[k]; ok {
				newO[k] = merge(oldV, newV)
			} else {
				newO[k] = oldV
			}
		}
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		oldO, ok := oldJ.(map[string]interface{})
		if ok {
			return oldO
		}
	}
	return newJ
}
