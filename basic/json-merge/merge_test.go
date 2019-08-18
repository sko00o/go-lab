package main

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

func Test_merge(t *testing.T) {
	type args struct {
		x1 interface{}
		x2 interface{}
	}

	x1 := `{"a":"asd","b":{"c":"qwe","d":{"e":1000,"f":"911"}}}`
	x2 := `{"a":123,"b":{"d":{"e":2333,"d":"qq"}}}`
	x3 := `{"a":123,"b":{"c":"qwe","d":{"e":2333,"f":"911","d":"qq"}}}`
	var v1 interface{}
	var v2 interface{}
	json.Unmarshal([]byte(x1), &v1)
	json.Unmarshal([]byte(x2), &v2)
	var v3 interface{}
	if err := json.Unmarshal([]byte(x3), &v3); err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				x1: v1,
				x2: v2,
			},
			want: v3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := merge(tt.args.x1, tt.args.x2)
			if (err != nil) != tt.wantErr {
				t.Errorf("merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
