package csv_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"go.nanasi880.dev/x/encoding/csv"
)

func ExampleMarshal() {

	type csvData struct {
		Name   string
		Age    int
		Memo   string `csv:"Comment"`
		Ignore string `csv:"-"`
	}

	d := []csvData{
		{Name: "Bob", Age: 18, Memo: "my name is Bob", Ignore: "Hi"},
		{Name: "Alice", Age: 18, Memo: "my name is Alice", Ignore: "Hi"},
	}

	encoded, err := csv.Marshal(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(encoded))
	// Output:
	// Name,Age,Comment
	// Bob,18,my name is Bob
	// Alice,18,my name is Alice
}

func ExampleEncoder_Encode() {

	type csvData struct {
		Name   string
		Age    int
		Memo   string `csv:"Comment"`
		Ignore string `csv:"-"`
	}

	d := []csvData{
		{Name: "Bob", Age: 18, Memo: "my name is Bob", Ignore: "Hi"},
		{Name: "Alice", Age: 18, Memo: "my name is Alice", Ignore: "Hi"},
	}

	err := csv.NewEncoder(os.Stdout).Encode(d)
	if err != nil {
		panic(err)
	}
	// Output:
	// Name,Age,Comment
	// Bob,18,my name is Bob
	// Alice,18,my name is Alice
}

func TestEncoder_Encode(t *testing.T) {

	data := []struct {
		Name      string
		UseHeader bool
		Comma     rune
		UseCRLF   bool
		Nil       string
		ToCSV     interface{}
		Want      string
	}{
		{
			Name:      "UseHeader",
			UseHeader: true,
			Comma:     ',',
			UseCRLF:   false,
			Nil:       "null",
			ToCSV: []*struct {
				V1 int       `csv:"Col1"`
				V2 float64   `csv:"Col2"`
				V3 time.Time `csv:"Col3"`
			}{
				{
					V1: 42,
					V2: 1.5,
					V3: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			Want: "Col1,Col2,Col3\n42,1.5,2020-01-01T12:00:00Z\n",
		},
		{
			Name:      "NoUseHeader",
			UseHeader: false,
			Comma:     ',',
			UseCRLF:   false,
			Nil:       "null",
			ToCSV: []*struct {
				V1 int       `csv:"Col1"`
				V2 float64   `csv:"Col2"`
				V3 time.Time `csv:"Col3"`
			}{
				{
					V1: 42,
					V2: 1.5,
					V3: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			Want: "42,1.5,2020-01-01T12:00:00Z\n",
		},
		{
			Name:      "Null",
			UseHeader: false,
			Comma:     ',',
			UseCRLF:   false,
			Nil:       "null",
			ToCSV: struct {
				V1 int  `csv:"Col1"`
				V2 *int `csv:"Col2"`
			}{
				V1: 42,
				V2: nil,
			},
			Want: "42,null\n",
		},
		{
			Name:      "NoUseTag",
			UseHeader: true,
			Comma:     ',',
			UseCRLF:   false,
			Nil:       "",
			ToCSV: struct {
				Col1 int
				Col2 string
			}{
				Col1: 42,
				Col2: "hello",
			},
			Want: "Col1,Col2\n42,hello\n",
		},
		{
			Name:      "AllPrimitiveType",
			UseHeader: false,
			Comma:     ',',
			UseCRLF:   false,
			Nil:       "",
			ToCSV: struct {
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
			Want: "true,-8,-16,-32,-64,-128,8,16,32,64,128,1.5,1.5,(1+2i),(3+4i),hello\n",
		},
	}

	for _, data := range data {
		data := data
		t.Run(data.Name, func(t *testing.T) {

			out := new(strings.Builder)

			enc := csv.NewEncoder(out)
			enc.UseHeader = data.UseHeader
			enc.Comma = data.Comma
			enc.UseCRLF = data.UseCRLF
			enc.Nil = data.Nil

			err := enc.Encode(data.ToCSV)
			if err != nil {
				t.Fatal(err)
			}

			if data.Want != out.String() {
				t.Fatalf("want: %s got: %s", data.Want, out.String())
			}
		})
	}
}
