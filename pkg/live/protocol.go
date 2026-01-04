package live

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	HeaderLength = 16

	ProtoRaw    = 0
	ProtoInt    = 1
	ProtoZlib   = 2
	ProtoBrotli = 3

	OpHeartbeat      = 2
	OpHeartbeatReply = 3
	OpSendMsgReply   = 5
	OpAuth           = 7
	OpAuthReply      = 8
)

type Header struct {
	PacketLen uint32
	HeaderLen uint16
	ProtoVer  uint16
	Operation uint32
	Sequence  uint32
}

type AuthPayload struct {
	UID      int    `json:"uid"`
	RoomID   int    `json:"roomid"`
	ProtoVer int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
	Buvid    string `json:"buvid"`
}

func Pack(op uint32, payload []byte) []byte {
	header := Header{
		PacketLen: uint32(HeaderLength + len(payload)),
		HeaderLen: uint16(HeaderLength),
		ProtoVer:  1,
		Operation: op,
		Sequence:  1,
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, header)
	buf.Write(payload)

	return buf.Bytes()
}

func Unpack(data []byte) (*Header, []byte, error) {
	if len(data) < HeaderLength {
		return nil, nil, fmt.Errorf("data too short")
	}

	var header Header
	buf := bytes.NewReader(data[:HeaderLength])
	if err := binary.Read(buf, binary.BigEndian, &header); err != nil {
		return nil, nil, err
	}

	return &header, data[HeaderLength:], nil
}
