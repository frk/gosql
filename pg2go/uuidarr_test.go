package pg2go

import (
	"testing"
)

func TestUUIDArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }
	U := uuid16bytes

	testlist2{{
		valuer:  UUIDArrayFromByteArray16Slice,
		scanner: UUIDArrayToByteArray16Slice,
		data: []testdata{
			{input: [][16]byte(nil), output: [][16]byte(nil)},
			{input: [][16]byte{}, output: [][16]byte{}},
			{
				input:  [][16]byte{U("894c9a8b-bafd-48d7-a705-f0625b52793d")},
				output: [][16]byte{U("894c9a8b-bafd-48d7-a705-f0625b52793d")}},
			{
				input: [][16]byte{
					U("894c9a8b-bafd-48d7-a705-f0625b52793d"),
					U("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
				},
				output: [][16]byte{
					U("894c9a8b-bafd-48d7-a705-f0625b52793d"),
					U("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
				}},
		},
	}, {
		valuer:  UUIDArrayFromStringSlice,
		scanner: UUIDArrayToStringSlice,
		data: []testdata{
			{input: []string(nil), output: []string(nil)},
			{input: []string{}, output: []string{}},
			{
				input:  []string{"894c9a8b-bafd-48d7-a705-f0625b52793d"},
				output: []string{"894c9a8b-bafd-48d7-a705-f0625b52793d"}},
			{
				input: []string{
					"894c9a8b-bafd-48d7-a705-f0625b52793d",
					"25a2fcf3-ed09-4e95-9617-8bd40e266ca1",
				},
				output: []string{
					"894c9a8b-bafd-48d7-a705-f0625b52793d",
					"25a2fcf3-ed09-4e95-9617-8bd40e266ca1",
				}},
		},
	}, {
		valuer:  UUIDArrayFromByteSliceSlice,
		scanner: UUIDArrayToByteSliceSlice,
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{input: [][]byte{}, output: [][]byte{}},
			{
				input:  [][]byte{B("894c9a8b-bafd-48d7-a705-f0625b52793d")},
				output: [][]byte{B("894c9a8b-bafd-48d7-a705-f0625b52793d")}},
			{
				input: [][]byte{
					B("894c9a8b-bafd-48d7-a705-f0625b52793d"),
					B("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
				},
				output: [][]byte{
					B("894c9a8b-bafd-48d7-a705-f0625b52793d"),
					B("25a2fcf3-ed09-4e95-9617-8bd40e266ca1")}},
			{
				input: [][]byte{
					B("894c9a8b-bafd-48d7-a705-f0625b52793d"),
					nil,
					B("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
					nil,
				},
				output: [][]byte{
					B("894c9a8b-bafd-48d7-a705-f0625b52793d"),
					nil,
					B("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
					nil}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{894c9a8b-bafd-48d7-a705-f0625b52793d,25a2fcf3-ed09-4e95-9617-8bd40e266ca1}`),
				output: string(`{894c9a8b-bafd-48d7-a705-f0625b52793d,25a2fcf3-ed09-4e95-9617-8bd40e266ca1}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{894c9a8b-bafd-48d7-a705-f0625b52793d,25a2fcf3-ed09-4e95-9617-8bd40e266ca1}`),
				output: []byte(`{894c9a8b-bafd-48d7-a705-f0625b52793d,25a2fcf3-ed09-4e95-9617-8bd40e266ca1}`)},
		},
	}}.execute(t, "uuidarr")
}
