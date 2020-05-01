package nbt2json

import (
	"bytes"
	"crypto/sha1"
	"testing"
)

// NOTE: Only testing round-trips for consistency, but errors will still fail the test.

// TODO: More thorough test json

var testJson = `{
  "nbt": [
    {
      "tagType": 10,
      "name": "",
      "value": [
        {
          "tagType": 2,
          "name": "Test",
          "value": 256
        }
      ]
    }
  ]
}`

func TestRoundTrip(t *testing.T) {
	h := sha1.New()

	// Get first nbt from test json, get hash
	nbtData, err := Json2Nbt([]byte(testJson), Bedrock)
	if err != nil {
		t.Error(err.Error())
	}
	nbtHash := h.Sum(nbtData)

	// Put that nbt through to json, get hash
	jsonOut, err := Nbt2Json(nbtData, Bedrock, "")
	if err != nil {
		t.Error(err.Error())
	}
	jsonHash := h.Sum(jsonOut)

	// Back to nbt again
	nbtData, err = Json2Nbt([]byte(testJson), Bedrock)
	if err != nil {
		t.Error(err.Error())
	}
	nbtHash2 := h.Sum(nbtData)

	// Compare first and second nbt hashes
	if !bytes.Equal(nbtHash, nbtHash2) {
		t.Error("Round trip NBT hashes don't match")
	}

	// Back to json again
	jsonOut, err = Nbt2Json(nbtData, Bedrock, "")
	if err != nil {
		t.Error(err.Error())
	}
	jsonHash2 := h.Sum(jsonOut)

	// Compare two generated json hashes
	if !bytes.Equal(jsonHash, jsonHash2) {
		t.Error("Round trip JSON hashes don't match")
	}
}
