package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type MyUser struct {
	ID       int64     `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	LastSeen time.Time `json:"last_seen,omitempty"`
}

func (u *MyUser) MarshalJSON() ([]byte, error) {
	type Alias MyUser
	return json.Marshal(&struct {
		LastSeen int64 `json:"last_seen,omitempty"`
		*Alias
	}{
		LastSeen: u.LastSeen.Unix(),
		Alias:    (*Alias)(u),
	})
}

func (u *MyUser) UnmarshalJSON(data []byte) error {
	type Alias MyUser
	aux := &struct {
		LastSeen int64 `json:"last_seen,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	u.LastSeen = time.Unix(aux.LastSeen, 0)
	return nil
}

func main() {
	var buf bytes.Buffer

	u := MyUser{ID: 1, Name: "L", LastSeen: time.Now()}
	_ = json.NewEncoder(&buf).Encode(&u)

	fmt.Println(buf.String())

	var uu MyUser
	json.NewDecoder(&buf).Decode(&uu)
	fmt.Printf("%+v\n", uu)
}
