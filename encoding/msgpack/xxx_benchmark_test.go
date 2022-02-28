package msgpack_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"math"
	"runtime"
	"testing"
	"time"

	"go.nanasi880.dev/x/encoding/gobutil"
	"go.nanasi880.dev/x/encoding/msgpack"
)

// MessagePack vs JSON and Gob
func Benchmark_VS_Standard(b *testing.B) {
	type benchmarkDataType struct {
		Int8     int8      `json:"0" msgpack:"0"`
		Int16    int16     `json:"1" msgpack:"1"`
		Int32    int32     `json:"2" msgpack:"2"`
		Int64    int64     `json:"3" msgpack:"3"`
		UInt8    uint8     `json:"4" msgpack:"4"`
		UInt16   uint16    `json:"5" msgpack:"5"`
		UInt32   uint32    `json:"6" msgpack:"6"`
		UInt64   uint64    `json:"7" msgpack:"7"`
		String   string    `json:"8" msgpack:"8"`
		DateTime time.Time `json:"9" msgpack:"9"`
	}
	benchmarkData := benchmarkDataType{
		Int8:     math.MaxInt8,
		Int16:    math.MaxInt16,
		Int32:    math.MaxInt32,
		Int64:    math.MaxInt64,
		UInt8:    math.MaxUint8,
		UInt16:   math.MaxUint16,
		UInt32:   math.MaxUint32,
		UInt64:   math.MaxUint64,
		String:   "Benchmark",
		DateTime: time.Date(2006, 1, 2, 15, 4, 5, 700, time.UTC),
	}
	jsonData, err := json.Marshal(benchmarkData)
	if err != nil {
		b.Fatal(err)
	}
	msgpackArrayData, err := msgpack.Marshal(benchmarkData)
	if err != nil {
		b.Fatal(err)
	}
	msgpackMapData, err := msgpack.MarshalStringKey(benchmarkData)
	if err != nil {
		b.Fatal(err)
	}
	gobData, err := gobutil.Marshal(benchmarkData)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Encode", func(b *testing.B) {
		b.Run("JSON", func(b *testing.B) {
			runtime.GC()
			encoder := json.NewEncoder(io.Discard)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := encoder.Encode(&benchmarkData)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("Gob", func(b *testing.B) {
			runtime.GC()
			encoder := gob.NewEncoder(io.Discard)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := encoder.Encode(&benchmarkData)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("MessagePack(Map)", func(b *testing.B) {
			runtime.GC()
			encoder := msgpack.NewEncoder(io.Discard)
			encoder.SetStructKeyType(msgpack.StructKeyTypeString)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := encoder.Encode(&benchmarkData)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("MessagePack(Array)", func(b *testing.B) {
			runtime.GC()
			encoder := msgpack.NewEncoder(io.Discard)
			encoder.SetStructKeyType(msgpack.StructKeyTypeInt)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := encoder.Encode(&benchmarkData)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run("Decode", func(b *testing.B) {
		b.Run("JSON", func(b *testing.B) {
			runtime.GC()
			var (
				reader  = bytes.NewReader(jsonData)
				decoder = json.NewDecoder(reader)
				decoded benchmarkDataType
			)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader.Reset(jsonData)
				err := decoder.Decode(&decoded)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("Gob", func(b *testing.B) {
			runtime.GC()
			var (
				reader  = bytes.NewReader(gobData)
				decoded benchmarkDataType
			)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader.Reset(gobData)
				decoder := gob.NewDecoder(reader)
				err := decoder.Decode(&decoded)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("MessagePack(Map)", func(b *testing.B) {
			runtime.GC()
			var (
				reader  = bytes.NewReader(msgpackMapData)
				decoder = msgpack.NewDecoder(reader)
				decoded benchmarkDataType
			)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader.Reset(msgpackMapData)
				decoder.Reset(reader)
				decoder.SetStructKeyType(msgpack.StructKeyTypeString)
				err := decoder.Decode(&decoded)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("MessagePack(Array)", func(b *testing.B) {
			runtime.GC()
			var (
				reader  = bytes.NewReader(msgpackArrayData)
				decoder = msgpack.NewDecoder(reader)
				decoded benchmarkDataType
			)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader.Reset(msgpackArrayData)
				decoder.Reset(reader)
				decoder.SetStructKeyType(msgpack.StructKeyTypeInt)
				err := decoder.Decode(&decoded)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("MessagePack(Array,[]byte)", func(b *testing.B) {
			runtime.GC()
			var (
				decoder = msgpack.NewDecoderBytes(msgpackArrayData)
				decoded benchmarkDataType
			)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				decoder.ResetBytes(msgpackArrayData)
				decoder.SetStructKeyType(msgpack.StructKeyTypeInt)
				err := decoder.Decode(&decoded)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
