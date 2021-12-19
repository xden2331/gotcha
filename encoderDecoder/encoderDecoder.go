package encoderDecoder

import "fmt"

type EncoderDecoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(bs []byte, v interface{}) error
}

type EncoderDecoderType int64

const (
	EncoderDecoderType_OfficialJSON EncoderDecoderType = 0
)

func (e EncoderDecoderType) String() string {
	switch e {
	case EncoderDecoderType_OfficialJSON:
		return "golang built-in JSON EncoderDecoder"
	default:
		return fmt.Sprintf("Unknown EncoderDecoderType [%v]", int64(e))
	}
}

func NewEncoderDecoder(edType EncoderDecoderType) EncoderDecoder {
	switch edType {
	default: // official json or unknown type would return the offical encoder-decoder
		return officialJSON(0)
	}
}