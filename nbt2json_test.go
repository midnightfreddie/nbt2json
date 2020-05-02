package nbt2json

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"testing"
)

// NOTE: Only testing round-trips for consistency, but errors will still fail the test.

// TODO: Add list/array types to test json

const testJson = `{
  "nbt": [
    {
      "tagType": 10,
      "name": "",
      "value": [
        {
          "tagType": 1,
          "name": "TestByte",
          "value": 127
        },
        {
          "tagType": 2,
          "name": "TestShort",
          "value": 32767
        },
        {
          "tagType": 3,
          "name": "TestInt",
          "value": 2147483647
        },
        {
          "tagType": 4,
          "name": "TestLong",
          "value": 9223372036854775807
        },
        {
          "tagType": 5,
          "name": "TestFloat",
          "value": 1.234567e+38
        },
        {
          "tagType": 6,
          "name": "TestDouble",
          "value": 1.23456789012345e+307
        },
        {
          "tagType": 8,
          "name": "TestString",
          "value": "This is a test string"
        },
        {
          "tagType": 0,
          "name": "",
          "value": null
        }
      ]
    }
  ]
}`

/*
        {
          "tagType": 7,
          "name": "TestByteArray",
          "value": 256
        },
				{
					"tagType": 9,
					"name": "TestList",
					"value": 256
				},
				{
					"tagType": 11,
					"name": "TestIntArray",
					"value": 256
				},
				{
					"tagType": 12,
					"name": "TestByteArray",
					"value": 256
				},
*/

const testNumberRangeJsonTemplate = `{
  "nbt": [
    {
      "tagType": %d,
      "name": "%s",
      "value": %v
    }
  ]
}`

func TestRoundTrip(t *testing.T) {
	h := sha1.New()

	// Get first nbt from test json, get hash
	nbtData, err := Json2Nbt([]byte(testJson), Bedrock)
	if err != nil {
		t.Fatal("Error converting test json:", err.Error())
	}
	nbtHash := h.Sum(nbtData)

	// Put that nbt through to json, get hash
	jsonOut, err := Nbt2Json(nbtData, Bedrock, "")
	if err != nil {
		t.Fatal("Error in first Nbt2Json conversion:", err.Error())
	}
	jsonHash := h.Sum(jsonOut)

	// Back to nbt again
	nbtData, err = Json2Nbt([]byte(testJson), Bedrock)
	if err != nil {
		t.Fatal("Error converting generated json back to nbt:", err.Error())
	}
	nbtHash2 := h.Sum(nbtData)

	// Compare first and second nbt hashes
	if !bytes.Equal(nbtHash, nbtHash2) {
		t.Fatal("Round trip NBT hashes don't match")
	}

	// Back to json again
	jsonOut, err = Nbt2Json(nbtData, Bedrock, "")
	if err != nil {
		t.Fatal("Error in second Nbt2Json conversion:", err.Error())
	}
	jsonHash2 := h.Sum(jsonOut)

	// Compare two generated json hashes
	if !bytes.Equal(jsonHash, jsonHash2) {
		t.Fatal("Round trip JSON hashes don't match")
	}
}

// TODO: Test array tags
func TestValueConversions(t *testing.T) {
	intTags := []struct {
		tagType int64
		value   int64
		nbt     []byte
	}{
		{1, math.MaxInt8, []byte{1, 0, 0, 0x7f}},
		{1, math.MinInt8, []byte{1, 0, 0, 0x80}},
		{2, math.MaxInt16, []byte{2, 0, 0, 0xff, 0x7f}},
		{2, math.MinInt16, []byte{2, 0, 0, 0x00, 0x80}},
		{3, math.MaxInt32, []byte{3, 0, 0, 0xff, 0xff, 0xff, 0x7f}},
		{3, math.MinInt32, []byte{3, 0, 0, 0x00, 0x00, 0x00, 0x80}},
		{4, math.MaxInt64, []byte{4, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
		{4, math.MinInt64, []byte{4, 0, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}},
	}

	for _, tag := range intTags {
		nbtData, err := Json2Nbt([]byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value)), Bedrock)
		if err != nil {
			t.Error("Error in json conversion during range tests", err.Error())
		} else if !bytes.Equal(nbtData, tag.nbt) {
			t.Error(fmt.Sprintf("Tag type %d value %d, expected \n%s\n, got \n%s\n", tag.tagType, tag.value, hex.Dump(tag.nbt), hex.Dump(nbtData)))
		}
	}

	floatTags := []struct {
		tagType int64
		value   float64
		nbt     []byte
	}{
		{5, 0, []byte{5, 0, 0, 0x00, 0x00, 0x00, 0x00}},
		{5, math.MaxFloat32, []byte{5, 0, 0, 0xff, 0xff, 0x7f, 0x7f}},
		{5, math.SmallestNonzeroFloat32, []byte{5, 0, 0, 0x01, 0x00, 0x00, 0x00}},
		{6, 0, []byte{6, 0, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{6, math.MaxFloat64, []byte{6, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0x7f}},
		{6, math.SmallestNonzeroFloat64, []byte{6, 0, 0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, tag := range floatTags {
		nbtData, err := Json2Nbt([]byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value)), Bedrock)
		if err != nil {
			t.Error("Error in json conversion during range tests", err.Error())
		} else if !bytes.Equal(nbtData, tag.nbt) {
			t.Error(fmt.Sprintf("Tag type %d value %g, expected \n%s\n, got \n%s\n", tag.tagType, tag.value, hex.Dump(tag.nbt), hex.Dump(nbtData)))
		}
	}
}

// NOTE: Tested function should throw error to pass
func TestOutOfRange(t *testing.T) {
	intTags := []struct {
		tagType int64
		value   int64
	}{
		{1, math.MaxInt8 + 1},
		{1, math.MinInt8 - 1},
		{2, math.MaxInt16 + 1},
		{2, math.MinInt16 - 1},
		{3, math.MaxInt32 + 1},
		{3, math.MinInt32 - 1},
	}

	for _, tag := range intTags {
		_, err := Json2Nbt([]byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value)), Bedrock)
		if err == nil {
			t.Error(fmt.Sprintf("Tag type %d value %d failed to throw out of range error", tag.tagType, tag.value))
		}
	}

	floatTags := []struct {
		tagType int64
		value   float64
	}{
		{5, math.MaxFloat32 * 1.1},
		{5, -(math.MaxFloat32 * 1.1)},
		{5, math.SmallestNonzeroFloat32 / 1.1},
		{5, -(math.SmallestNonzeroFloat32 / 1.1)},
	}

	for _, tag := range floatTags {
		_, err := Json2Nbt([]byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value)), Bedrock)
		if err == nil {
			t.Error(fmt.Sprintf("Tag type %d value %g failed to throw out of range error", tag.tagType, tag.value))
		}
	}
}
