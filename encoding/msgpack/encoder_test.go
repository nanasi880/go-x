package msgpack_test

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"go.nanasi880.dev/x/encoding/msgpack"
)

func TestMarshal(t *testing.T) {

	data := &struct {
		I int         `msgpack:"0"`
		S string      `msgpack:"2"`
		T time.Time   `msgpack:"3"`
		N interface{} `msgpack:"4"`
	}{
		I: 42,
		S: "Hello",
		T: time.Date(2020, 12, 31, 12, 34, 56, 999, time.UTC),
		N: nil,
	}

	b, err := msgpack.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	_ = b
}

func TestMarshal_Cyclic(t *testing.T) {

	type data struct {
		P *data `msgpack:"0"`
	}

	root := new(data)
	root.P = &data{
		P: &data{
			P: root,
		},
	}

	_, err := msgpack.Marshal(root)
	if err == nil {
		t.Fatal()
	}
}

func TestEncoder_EncodeTestSuite(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Nil {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Bool", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Bool {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Uint8", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint8 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Uint16", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint16 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Uint32", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint32 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Uint64", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Uint64 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Int8", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int8 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Int16", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int16 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Int32", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int32 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Int64", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Int64 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Float32", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Float32 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Float64", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Float64 {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("String", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.String {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Binary", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Binary {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Array", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Array {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Map", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Map {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
	t.Run("Time", func(t *testing.T) {
		for suiteNo, suite := range msgPackTestSuite.Time {
			encoded, err := msgpack.Marshal(suite.Value)
			if err != nil {
				t.Fatal(err)
			}
			for _, pattern := range suite.MsgPack {
				if bytes.Compare(encoded, pattern) == 0 {
					goto PASS
				}
			}
			t.Fatalf("suite:%d hex:%s", suiteNo, hex.EncodeToString(encoded))
		PASS:
		}
	})
}

func TestEncoder_CustomStructTag(t *testing.T) {

	jsonKey := struct {
		Value1 string `json:"value_1"`
	}{
		Value1: "Hello",
	}

	buf := new(bytes.Buffer)
	enc := msgpack.NewEncoder(buf).
		SetStructKeyType(msgpack.StructKeyTypeString).
		SetStructTagName("json")

	err := enc.Encode(jsonKey)
	if err != nil {
		t.Fatal(err)
	}

	msgpackKey := struct {
		Value1 string `msgpack:"value_1"`
	}{}
	err = msgpack.UnmarshalStringKey(buf.Bytes(), &msgpackKey)
	if err != nil {
		t.Fatal(err)
	}
	if jsonKey.Value1 != msgpackKey.Value1 {
		t.Fatal(msgpackKey)
	}
}
