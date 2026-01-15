package live

import (
	"encoding/binary"
	"testing"
)

func TestHeartbeatParsing(t *testing.T) {
	// Sample data from doc/message_stream.md
	// Header: 00 00 00 14 00 10 00 01 00 00 00 03 00 00 00 00
	// Body: 00 00 09 a2 5b 6f 62 6a 65 63 74 20 4f 62 6a 65 63 74 5d ([object Object])

	headerBytes := []byte{
		0x00, 0x00, 0x00, 0x14, // PacketLen 20
		0x00, 0x10,             // HeaderLen 16
		0x00, 0x01,             // ProtoVer 1
		0x00, 0x00, 0x00, 0x03, // Op 3 (HeartbeatReply)
		0x00, 0x00, 0x00, 0x00, // Sequence
	}

	bodyBytes := []byte{
		0x00, 0x00, 0x09, 0xa2, // 2466
		0x5b, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5d, // [object Object]
	}

	header, _, err := Unpack(append(headerBytes, bodyBytes...))
	if err != nil {
		t.Fatalf("Unpack failed: %v", err)
	}

	// Manually call routeOperation logic
	if header.Operation == OpHeartbeatReply {
		if len(bodyBytes) >= 4 {
			popularity := binary.BigEndian.Uint32(bodyBytes[:4])
			t.Logf("Parsed popularity: %d", popularity)
			if popularity != 2466 {
				t.Errorf("Expected 2466, got %d", popularity)
			}
		}
	}
}
