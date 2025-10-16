package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/ggymm/gnet"
)

// 协议
// head      2
// length    4
// content   length

var (
	bin  = binary.BigEndian
	head = make([]byte, headSize)
)

var (
	errIncompletePacket    = errors.New("incomplete packet")
	errInvalidPacketHeader = errors.New("invalid packet header")
)

const (
	headSize   = 2
	lengthSize = 4
	offsetSize = headSize + lengthSize
)

func init() {
	bin.PutUint16(head, 10808)
}

type SocketCodec struct {
}

func (codec SocketCodec) Decode(c gnet.Conn) ([]byte, error) {
	data, err := c.Peek(offsetSize)
	if err != nil {
		if errors.Is(err, io.ErrShortBuffer) {
			err = errIncompletePacket
		}
		return nil, err
	}

	// 读取协议头
	if !bytes.Equal(head, data[:headSize]) {
		return nil, errInvalidPacketHeader
	}

	// 读取消息长度
	bodySize := bin.Uint32(data[headSize:offsetSize])
	dataSize := offsetSize + int(bodySize)

	// 读取消息内容
	data, err = c.Peek(dataSize)
	if err != nil {
		if errors.Is(err, io.ErrShortBuffer) {
			err = errIncompletePacket
		}
		return nil, err
	}
	_, _ = c.Discard(dataSize)

	// 只返回消息体内容
	body := make([]byte, bodySize)
	copy(body, data[offsetSize:dataSize])
	return body, nil
}

func (codec SocketCodec) Encode(buf []byte) []byte {
	size := offsetSize + len(buf)
	data := make([]byte, size)

	// 写入协议头
	copy(data, head)

	// 写入消息长度
	bin.PutUint32(data[headSize:offsetSize], uint32(len(buf)))

	// 写入消息内容
	copy(data[offsetSize:], buf)
	return data
}
