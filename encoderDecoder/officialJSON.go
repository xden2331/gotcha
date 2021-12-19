package encoderDecoder

import "encoding/json"

type officialJSON int

func (officialJSON) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (officialJSON) Decode(bs []byte, v interface{}) error {
	return json.Unmarshal(bs, v)
}
