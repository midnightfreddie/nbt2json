package nbt2json

import (
	"bytes"
	"crypto/sha1"
	"testing"
)

// NOTE: Only testing round-trips for consistency, but errors will still fail the test.

// TODO: Add list/array types to test json

var testJson = `{
  "nbt": [
    {
      "tagType": 10,
      "name": "",
      "value": [
        {
          "tagType": 1,
          "name": "TestByte",
          "value": 64
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
