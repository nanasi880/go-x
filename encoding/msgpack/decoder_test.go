package msgpack_test

import (
	"bytes"
	"math"
	"reflect"
	"testing"
	"time"

	"go.nanasi880.dev/x/encoding/base64"
	"go.nanasi880.dev/x/encoding/msgpack"
)

func TestDecoder_DecodeFromMessagePackForCSharpBinary(t *testing.T) {
	binary, err := base64.DecodeFromFile(base64.StdEncoding, "testdata/test-data-base64.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("decode to struct", func(t *testing.T) {
		t.Parallel()
		type Data struct {
			Int8     int8
			Int16    int16
			Int32    int32
			Int64    int64
			UInt8    uint8
			UInt16   uint16
			UInt32   uint32
			UInt64   uint64
			String   string
			DateTime time.Time
		}
		var (
			decoded Data
			decoder = msgpack.NewDecoder(bytes.NewReader(binary))
		)
		decoder.SetStructKeyType(msgpack.StructKeyTypeString)
		err := decoder.Decode(&decoded)
		if err != nil {
			t.Fatal(err)
		}

		if decoded.Int8 != int8(math.MaxInt8/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.Int16 != int16(math.MaxInt16/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.Int32 != int32(math.MaxInt32/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.Int64 != int64(math.MaxInt64/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.UInt8 != uint8(math.MaxUint8/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.UInt16 != uint16(math.MaxUint16/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.UInt32 != uint32(math.MaxUint32/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.UInt64 != uint64(math.MaxUint64/2)+1 {
			t.Fatal(decoded)
		}
		if decoded.String != "Hello" {
			t.Fatal(decoded)
		}
		if !decoded.DateTime.Equal(time.Date(2021, 5, 25, 12, 34, 56, int((789 * time.Millisecond).Nanoseconds()), time.UTC)) {
			t.Fatal(decoded)
		}
	})
	t.Run("decode to struct pointer", func(t *testing.T) {
		t.Parallel()
		type Data struct {
			Int8     *int8
			Int16    *int16
			Int32    *int32
			Int64    *int64
			UInt8    *uint8
			UInt16   *uint16
			UInt32   *uint32
			UInt64   *uint64
			String   *string
			DateTime *time.Time
		}
		var (
			decoded Data
			decoder = msgpack.NewDecoder(bytes.NewReader(binary))
		)
		decoder.SetStructKeyType(msgpack.StructKeyTypeString)
		err := decoder.Decode(&decoded)
		if err != nil {
			t.Fatal(err)
		}

		if *decoded.Int8 != int8(math.MaxInt8/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.Int16 != int16(math.MaxInt16/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.Int32 != int32(math.MaxInt32/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.Int64 != int64(math.MaxInt64/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.UInt8 != uint8(math.MaxUint8/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.UInt16 != uint16(math.MaxUint16/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.UInt32 != uint32(math.MaxUint32/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.UInt64 != uint64(math.MaxUint64/2)+1 {
			t.Fatal(decoded)
		}
		if *decoded.String != "Hello" {
			t.Fatal(decoded)
		}
		if !decoded.DateTime.Equal(time.Date(2021, 5, 25, 12, 34, 56, int((789 * time.Millisecond).Nanoseconds()), time.UTC)) {
			t.Fatal(decoded)
		}
	})
	t.Run("decode to interface", func(t *testing.T) {
		t.Parallel()
		var (
			decoded interface{}
			decoder = msgpack.NewDecoder(bytes.NewReader(binary))
		)
		decoder.SetStructKeyType(msgpack.StructKeyTypeString)
		err := decoder.Decode(&decoded)
		if err != nil {
			t.Fatal(err)
		}

		decodedMap, ok := decoded.(map[interface{}]interface{})
		if !ok {
			t.Fatal(decoded)
		}
		if decodedMap["Int8"] != uint64(math.MaxInt8/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["Int16"] != uint64(math.MaxInt16/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["Int32"] != uint64(math.MaxInt32/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["Int64"] != uint64(math.MaxInt64/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt8"] != uint64(math.MaxUint8/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt16"] != uint64(math.MaxUint16/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt32"] != uint64(math.MaxUint32/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt64"] != uint64(math.MaxUint64/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["String"] != "Hello" {
			t.Fatal(decoded)
		}
		if !decodedMap["DateTime"].(time.Time).Equal(time.Date(2021, 5, 25, 12, 34, 56, int((789 * time.Millisecond).Nanoseconds()), time.UTC)) {
			t.Fatal(decoded)
		}
	})
	t.Run("decode to map", func(t *testing.T) {
		t.Parallel()
		var (
			decoded map[string]interface{}
			decoder = msgpack.NewDecoder(bytes.NewReader(binary))
		)
		decoder.SetStructKeyType(msgpack.StructKeyTypeString)
		err := decoder.Decode(&decoded)
		if err != nil {
			t.Fatal(err)
		}

		decodedMap := decoded
		if decodedMap["Int8"] != uint64(math.MaxInt8/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["Int16"] != uint64(math.MaxInt16/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["Int32"] != uint64(math.MaxInt32/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["Int64"] != uint64(math.MaxInt64/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt8"] != uint64(math.MaxUint8/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt16"] != uint64(math.MaxUint16/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt32"] != uint64(math.MaxUint32/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["UInt64"] != uint64(math.MaxUint64/2)+1 {
			t.Fatal(decoded)
		}
		if decodedMap["String"] != "Hello" {
			t.Fatal(decoded)
		}
		if !decodedMap["DateTime"].(time.Time).Equal(time.Date(2021, 5, 25, 12, 34, 56, int((789 * time.Millisecond).Nanoseconds()), time.UTC)) {
			t.Fatal(decoded)
		}
	})
}

func TestDecoder_DecodeTestSuite(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Nil {
			for binaryNo, binary := range suite.MsgPack {
				var decoded interface{}
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Bool", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Bool {
			for binaryNo, binary := range suite.MsgPack {
				var decoded bool
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Uint8", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint8 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded uint8
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Uint16", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint16 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded uint16
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Uint32", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint32 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded uint32
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Uint64", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint64 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded uint64
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Int8", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int8 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded int8
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Int16", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int16 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded int16
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Int32", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int32 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded int32
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Int64", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int64 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded int64
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Float32", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Float32 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded float32
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Float64", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Float64 {
			for binaryNo, binary := range suite.MsgPack {
				var decoded float64
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("String", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.String {
			for binaryNo, binary := range suite.MsgPack {
				var decoded string
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if decoded != suite.Value {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Binary", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Binary {
			for binaryNo, binary := range suite.MsgPack {
				var decoded []byte
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if bytes.Compare(decoded, suite.Value) != 0 {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Array", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Array {
			for binaryNo, binary := range suite.MsgPack {
				var decoded []int
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if !reflect.DeepEqual(decoded, suite.Value) {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Map", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Map {
			for binaryNo, binary := range suite.MsgPack {
				var decoded map[int]int
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if !reflect.DeepEqual(decoded, suite.Value) {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
	t.Run("Time", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Time {
			for binaryNo, binary := range suite.MsgPack {
				var decoded time.Time
				err := msgpack.Unmarshal(binary, &decoded)
				if err != nil {
					t.Logf("suite:%d binary:%d err:%v  got:%v", suiteNo, binaryNo, err, decoded)
					t.Fail()
					continue
				}
				if !suite.Value.Equal(decoded) {
					t.Logf("suite:%d binary:%d got:%v", suiteNo, binaryNo, decoded)
					t.Fail()
				}
			}
		}
	})
}
