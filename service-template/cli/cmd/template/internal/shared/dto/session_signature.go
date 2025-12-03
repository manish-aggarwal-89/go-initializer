package dto

import "encoding/json"

type SessionSignatureDto struct {
	AuthToken string `json:"auth_token"`
}

func (r SessionSignatureDto) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SessionSignatureDto) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &r)
}
