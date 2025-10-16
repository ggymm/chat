package json

import (
	"github.com/bytedance/sonic"
)

func Marshal(v interface{}) []byte {
	ret, _ := sonic.Marshal(v)
	return ret
}

func MarshalE(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) {
	_ = sonic.Unmarshal(data, v)
	return
}

func UnmarshalE(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}
