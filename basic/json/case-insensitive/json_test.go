package json_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCaseSense(t *testing.T) {
	assert := require.New(t)

	data := []byte(`{"UPPER":1,"lower":2}`)

	// https://github.com/golang/go/blob/master/src/encoding/json/decode.go#L45
	type J struct {
		Upper int `json:"upper"`
		Lower int `json:"lower"`
	}

	var j J
	assert.NoError(json.Unmarshal(data, &j))
	t.Log(j)

	// https://github.com/golang/go/blob/master/src/encoding/json/decode.go#L717
	type J1 struct {
		Upper int `json:"Upper"`
		UPper int `json:"UPper"`
		UpPer int `json:"UpPer"`
		UppEr int `json:"UppEr"`
		Lower int `json:"lower"`
	}

	var j1 J1
	assert.NoError(json.Unmarshal(data, &j1))
	t.Log(j1)
}
