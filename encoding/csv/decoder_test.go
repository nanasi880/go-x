package csv_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"go.nanasi880.dev/x/encoding/csv"
)

func ExampleUnmarshalString() {
	const csvString = `Name,Age,Comment
Bob,18,my name is Bob
Alice,18,my name is Alice`

	type csvData struct {
		Name   string
		Age    int
		Memo   string `csv:"Comment"`
		Ignore int    `csv:"-"`
	}

	var d []csvData
	if err := csv.UnmarshalString(csvString, &d); err != nil {
		panic(err)
	}

	fmt.Println(len(d))
	fmt.Println(d[0])
	fmt.Println(d[1])
	// Output:
	// 2
	// {Bob 18 my name is Bob 0}
	// {Alice 18 my name is Alice 0}
}

func TestDecoder_Decode(t *testing.T) {

	data := []struct {
		Name      string
		UseHeader bool
		Nil       string
		CSV       interface{}
		Want      []string
	}{
		{
			Name:      "AllPrimitiveType",
			UseHeader: true,
			Nil:       "",
			CSV: []*struct {
				Bool       bool
				Int8       int8
				Int16      int16
				Int32      int32
				Int64      int64
				Int        int
				Uint8      uint8
				Uint16     uint16
				Uint32     uint32
				Uint64     uint64
				Uint       uint
				Float32    float32
				Float64    float64
				Complex64  complex64
				Complex128 complex128
				String     string
			}{
				{
					Bool:       true,
					Int8:       -8,
					Int16:      -16,
					Int32:      -32,
					Int64:      -64,
					Int:        -128,
					Uint8:      8,
					Uint16:     16,
					Uint32:     32,
					Uint64:     64,
					Uint:       128,
					Float32:    1.5,
					Float64:    1.5,
					Complex64:  complex(1, 2),
					Complex128: complex(3, 4),
					String:     "hello",
				},
			},
			Want: []string{"&{true -8 -16 -32 -64 -128 8 16 32 64 128 1.5 1.5 (1+2i) (3+4i) hello}"},
		},
	}

	for _, data := range data {
		data := data
		t.Run(data.Name, func(t *testing.T) {

			buf := new(bytes.Buffer)
			enc := csv.NewEncoder(buf)

			enc.Nil = data.Nil
			enc.UseHeader = data.UseHeader

			if err := enc.Encode(data.CSV); err != nil {
				t.Fatal(err)
			}

			dec := csv.NewDecoder(buf)
			dec.Nil = data.Nil
			dec.UseHeader = data.UseHeader

			out := reflect.New(reflect.TypeOf(data.CSV))
			if err := dec.Decode(out.Interface()); err != nil {
				t.Fatal(err)
			}

			{
				out := out.Elem()
				for i := 0; i < out.Len(); i++ {
					s := fmt.Sprint(out.Index(i).Interface())
					if s != data.Want[i] {
						t.Fatal(i, s)
					}
				}
			}
		})
	}
}
