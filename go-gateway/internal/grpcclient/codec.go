package grpcclient

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const jsonCodecName = "json"

// jsonCodec JSON 编解码器（用于开发阶段，生产建议使用 protobuf）
type jsonCodec struct{}

func (j jsonCodec) Name() string {
	return jsonCodecName
}

func (j jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {
	encoding.RegisterCodec(jsonCodec{})
}
