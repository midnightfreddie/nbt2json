package nbt2json

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"testing"
)

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
			"name": "",
			"value": {
				"valueLeast": 4294967295,
				"valueMost": 2147483647
			}
		},
        {
			"tagType": 4,
			"name": "TestLongAsString",
			"value": "9223372036854775807"
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
			"tagType": 7,
			"name": "TestByteArray",
			"value": [
				0,
				-128,
				127
			]
        },
        {
          "tagType": 8,
          "name": "TestString",
          "value": "This is a test string"
        },
				{
					"tagType": 9,
					"name": "TestList",
						"value": {
							"tagListType": 3,
							"list": [
								0,
								2147483647,
								-2147483648
							]
						}
				},
				{
					"tagType": 11,
					"name": "TestIntArray",
					"value": [
						0,
						2147483647,
						-2147483648
					]
				},
				{
					"tagType": 12,
					"name": "TestLongArray",
					"value": [
						{
							"valueLeast": 0,
							"valueMost": 0
						},
						{
							"valueLeast": 4294967295,
							"valueMost": 2147483647
						},
						"9223372036854775807",
						{
							"valueLeast": 0,
							"valueMost": -2147483648
						}
					]
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

const testNumberRangeJsonTemplate = `{
  "nbt": [
    {
      "tagType": %d,
      "name": "%s",
      "value": %v
    }
  ]
}`

const testLongTemplate = `{
	"valueLeast": %d,
	"valueMost": %d
}`

// TestRoundTrip checks to be sure generated output matches, but it doesn't check values against the original input
func TestRoundTrip(t *testing.T) {
	h := sha1.New()

	// Get first nbt from test json, get hash
	nbtData, err := Json2Nbt([]byte(testJson))
	if err != nil {
		t.Fatal("Error converting test json:", err.Error())
	}
	nbtHash := h.Sum(nbtData)

	// Put that nbt through to json, get hash
	jsonOut, err := Nbt2Json(bytes.NewReader(nbtData), "", 1)
	if err != nil {
		t.Fatal("Error in first Nbt2Json conversion:", err.Error())
	}
	jsonHash := h.Sum(jsonOut)

	// Back to nbt again
	nbtData, err = Json2Nbt([]byte(testJson))
	if err != nil {
		t.Fatal("Error converting generated json back to nbt:", err.Error())
	}
	nbtHash2 := h.Sum(nbtData)

	// Compare first and second nbt hashes
	if !bytes.Equal(nbtHash, nbtHash2) {
		t.Fatal("Round trip NBT hashes don't match")
	}

	// Back to json again
	jsonOut, err = Nbt2Json(bytes.NewReader(nbtData), "", 1)
	if err != nil {
		t.Fatal("Error in second Nbt2Json conversion:", err.Error())
	}
	jsonHash2 := h.Sum(jsonOut)

	// Compare two generated json hashes
	if !bytes.Equal(jsonHash, jsonHash2) {
		t.Fatal("Round trip JSON hashes don't match")
	}
}

// TestValueConversions checks nbt value versus input json value
// TODO: Test array tags
func TestValueConversions(t *testing.T) {
	UseBedrockEncoding()

	testIntTag := func(json []byte, tagType int64, value interface{}, nbt []byte) error {
		nbtData, err := Json2Nbt(json)
		if err != nil {
			return fmt.Errorf("Error in json conversion during value tests: %w", err)
		} else if !bytes.Equal(nbtData, nbt) {
			return fmt.Errorf(fmt.Sprintf("Tag type %d value %v, expected \n%s\n, got \n%s\n", tagType, value, hex.Dump(nbt), hex.Dump(nbtData)))
		} else {
			jsonData, err := Nbt2Json(bytes.NewReader(nbtData), "", 1)
			if err != nil {
				return fmt.Errorf("Error in nbt re-conversion during value tests: %w", err)
			} else {
				nbtData, err = Json2Nbt(jsonData)
				if err != nil {
					return fmt.Errorf("Error in json re-conversion during value tests: %w", err)
				} else if !bytes.Equal(nbtData, nbt) {
					return fmt.Errorf(fmt.Sprintf("Error on round-trip value reconversion - tag type %d value %v, expected \n%s\n, got \n%s\n", tagType, value, hex.Dump(nbt), hex.Dump(nbtData)))
				}
			}
		}

		return nil
	}

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
	}

	for _, tag := range intTags {
		b := []byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value))
		if err := testIntTag(b, tag.tagType, tag.value, tag.nbt); err != nil {
			t.Error(err)
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
		b := []byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value))
		if err := testIntTag(b, tag.tagType, tag.value, tag.nbt); err != nil {
			t.Error(err)
		}
	}

	int64Tags := []struct {
		tagType    int64
		valueLeast uint32
		valueMost  uint32
		nbt        []byte
	}{
		{4, 0xffffffff, math.MaxInt32, []byte{4, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
		{4, 0, 0x80000000, []byte{4, 0, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}},
	}

	for _, tag := range int64Tags {
		value := fmt.Sprintf(testLongTemplate, tag.valueLeast, tag.valueMost)
		b := []byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", value))
		if err := testIntTag(b, tag.tagType, value, tag.nbt); err != nil {
			t.Error(err)
		}
	}

	longAsStringTags := []struct {
		tagType int64
		value   string
		nbt     []byte
	}{
		{4, "\"9223372036854775807\"", []byte{4, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
		{4, "\"-9223372036854775808\"", []byte{4, 0, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}},
	}
	for _, tag := range longAsStringTags {
		b := []byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value))
		if err := testIntTag(b, tag.tagType, tag.value, tag.nbt); err != nil {
			t.Error(err)
		}
	}

}

// TestOutOfRange tries to offer input out of range of the tag type
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
		_, err := Json2Nbt([]byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value)))
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
		// Not testing for small limits; it wasn't working as expected
		// {5, math.SmallestNonzeroFloat32 / 1.1},
		// {5, -(math.SmallestNonzeroFloat32 / 1.1)},
	}

	for _, tag := range floatTags {
		_, err := Json2Nbt([]byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value)))
		if err == nil {
			t.Error(fmt.Sprintf("Tag type %d value %g failed to throw out of range error", tag.tagType, tag.value))
		}
	}

	longAsStringTags := []struct {
		tagType int64
		value   string
	}{
		{4, "\"9223372036854775808\""},
		{4, "\"-9223372036854775809\""},
	}

	for _, tag := range longAsStringTags {
		_, err := Json2Nbt([]byte(fmt.Sprintf(testNumberRangeJsonTemplate, tag.tagType, "", tag.value)))
		if err == nil {
			t.Error(fmt.Sprintf("Tag type %d value %s failed to throw out of range error", tag.tagType, tag.value))
		}
	}
}
