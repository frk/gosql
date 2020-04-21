package pg2go

import (
	"testing"
)

func TestUUID(t *testing.T) {
	testlist2{{
		valuer:  UUIDFromByteArray16,
		scanner: UUIDToByteArray16,
		data: []testdata{
			{
				input:  uuid16bytes("894c9a8b-bafd-48d7-a705-f0625b52793d"),
				output: uuid16bytes("894c9a8b-bafd-48d7-a705-f0625b52793d")},
			{
				input:  uuid16bytes("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
				output: uuid16bytes("25a2fcf3-ed09-4e95-9617-8bd40e266ca1")},
		},
	}, {
		data: []testdata{
			{
				input:  string("894c9a8b-bafd-48d7-a705-f0625b52793d"),
				output: string("894c9a8b-bafd-48d7-a705-f0625b52793d")},
			{
				input:  string("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
				output: string("25a2fcf3-ed09-4e95-9617-8bd40e266ca1")},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("894c9a8b-bafd-48d7-a705-f0625b52793d"),
				output: []byte("894c9a8b-bafd-48d7-a705-f0625b52793d")},
			{
				input:  []byte("25a2fcf3-ed09-4e95-9617-8bd40e266ca1"),
				output: []byte("25a2fcf3-ed09-4e95-9617-8bd40e266ca1")},
		},
	}}.execute(t, "uuid")
}
