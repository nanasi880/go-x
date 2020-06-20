package csv_test

import (
	"strings"
	"testing"
	"time"

	"go.nanasi880.dev/x/encoding/csv"
)

func TestEncoder_Encode(t *testing.T) {

	data := []struct {
		Name      string
		UseHeader bool
		Comma     rune
		UseCRLF   bool
		ToCSV     interface{}
		Want      string
	}{
		{
			Name:      "UseHeader",
			UseHeader: true,
			Comma:     ',',
			UseCRLF:   false,
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
	}

	for _, data := range data {
		data := data
		t.Run(data.Name, func(t *testing.T) {

			out := new(strings.Builder)

			enc := csv.NewEncoder(out)
			enc.UseHeader = data.UseHeader
			enc.Comma = data.Comma
			enc.UseCRLF = data.UseCRLF

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
