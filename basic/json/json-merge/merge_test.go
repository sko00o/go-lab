package main

import (
	"reflect"
	"testing"
)

func TestOverride(t *testing.T) {
	type args struct {
		data1 []byte
		data2 []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				data1: []byte(`{"a":"asd","b":{"c":"qwe","d":{"e":1000,"f":"911"}}}`),
				data2: []byte(`{"a":123,"b":{"d":{"e":2333,"d":"qq"}}}`),
			},
			want: []byte(`{"a":123,"b":{"c":"qwe","d":{"d":"qq","e":2333,"f":"911"}}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Override(tt.args.data1, tt.args.data2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Override() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Override() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
