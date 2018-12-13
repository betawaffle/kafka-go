package kafka

import (
	"testing"
	"time"
)

var (
	produceRequestEmpty = []byte{
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}

	produceRequestHeader = []byte{
		0x01, 0x23,
		0x00, 0x00, 0x04, 0x44,
		0x00, 0x00, 0x00, 0x00}

	produceRequestOneMessage = []byte{
		0x01, 0x23,
		0x00, 0x00, 0x04, 0x44,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x05, 't', 'o', 'p', 'i', 'c',
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0xAD,
		0x00, 0x00, 0x00, 0x1C,
		// messageSet
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x10,
		// message
		0x23, 0x96, 0x4a, 0xf7, // CRC
		0x00,
		0x00,
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x02, 0x00, 0xEE}

	produceRequestOneRecord = []byte{
		0xFF, 0xFF, // Transaction ID
		0x01, 0x23, // Required Acks
		0x00, 0x00, 0x04, 0x44, // Timeout
		0x00, 0x00, 0x00, 0x01, // Number of Topics
		0x00, 0x05, 't', 'o', 'p', 'i', 'c', // Topic
		0x00, 0x00, 0x00, 0x01, // Number of Partitions
		0x00, 0x00, 0x00, 0xAD, // Partition
		0x00, 0x00, 0x00, 0x52, // Records length
		// recordBatch
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x46,
		0x00, 0x00, 0x00, 0x00,
		0x02,
		0xCA, 0x33, 0xBC, 0x05,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x01, 0x58, 0x8D, 0xCD, 0x59, 0x38,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01,
		// record
		0x28,
		0x00,
		0x0A,
		0x00,
		0x08, 0x01, 0x02, 0x03, 0x04,
		0x06, 0x05, 0x06, 0x07,
		0x02,
		0x06, 0x08, 0x09, 0x0A,
		0x04, 0x0B, 0x0C,
	}
)

func TestProduceRequest(t *testing.T) {
	request := new(ProduceRequest)
	testRequest(t, "empty", request, produceRequestEmpty)

	request.RequiredAcks = 0x123
	request.Timeout = 0x444
	testRequest(t, "header", request, produceRequestHeader)

	request.AddMessage("topic", 0xAD, &Message{Codec: CompressionNone, Key: nil, Value: []byte{0x00, 0xEE}})
	testRequest(t, "one message", request, produceRequestOneMessage)

	request.Version = 3
	batch := &RecordBatch{
		LastOffsetDelta: 1,
		Version:         2,
		FirstTimestamp:  time.Unix(1479847795, 0),
		MaxTimestamp:    time.Unix(0, 0),
		Records: []*Record{{
			TimestampDelta: 5 * time.Millisecond,
			Key:            []byte{0x01, 0x02, 0x03, 0x04},
			Value:          []byte{0x05, 0x06, 0x07},
			Headers: []*RecordHeader{{
				Key:   []byte{0x08, 0x09, 0x0A},
				Value: []byte{0x0B, 0x0C},
			}},
		}},
	}
	request.AddBatch("topic", 0xAD, batch)
	packet := testRequestEncode(t, "one record", request, produceRequestOneRecord)
	// compressRecords field is not populated on decoding because consumers
	// are only interested in decoded records.
	batch.compressedRecords = nil
	testRequestDecode(t, "one record", request, packet)
}
