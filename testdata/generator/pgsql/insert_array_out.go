// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/pgsql"
)

func (q *insertarrayquery) Exec(c gosql.Conn) error {
	const queryString = `INSERT INTO "pgsql_test" (
		"col_bitarr"
		, "col_boolarr"
		, "col_boxarr"
		, "col_bpchararr"
		, "col_byteaarr"
		, "col_chararr"
		, "col_cidrarr"
		, "col_datearr"
		, "col_daterangearr"
		, "col_float4arr"
		, "col_float8arr"
		, "col_inetarr"
		, "col_int2arr"
		, "col_int2vector"
		, "col_int2vectorarr"
		, "col_int4arr"
		, "col_int4rangearr"
		, "col_int8arr"
		, "col_int8rangearr"
		, "col_jsonarr"
		, "col_jsonbarr"
		, "col_linearr"
		, "col_lsegarr"
		, "col_macaddrarr"
		, "col_macaddr8arr"
		, "col_moneyarr"
		, "col_numericarr"
		, "col_numrangearr"
		, "col_patharr"
		, "col_pointarr"
		, "col_polygonarr"
		, "col_textarr"
		, "col_timearr"
		, "col_timestamparr"
		, "col_timestamptzarr"
		, "col_timetzarr"
		, "col_tsqueryarr"
		, "col_tsrangearr"
		, "col_tstzrangearr"
		, "col_tsvectorarr"
		, "col_uuidarr"
		, "col_varbitarr"
		, "col_varchararr"
		, "col_xmlarr"
	) VALUES (
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
		, $12
		, $13
		, $14
		, $15
		, $16
		, $17
		, $18
		, $19
		, $20
		, $21
		, $22
		, $23
		, $24
		, $25
		, $26
		, $27
		, $28
		, $29
		, $30
		, $31
		, $32
		, $33
		, $34
		, $35
		, $36
		, $37
		, $38
		, $39
		, $40
		, $41
		, $42
		, $43
		, $44
	)` // `

	_, err := c.Exec(queryString,
		pgsql.BitArrayFromUint8Slice(q.data.bitarr),
		pgsql.BoolArrayFromBoolSlice(q.data.boolarr),
		pgsql.BoxArrayFromFloat64Array2Array2Slice(q.data.boxarr),
		pgsql.BPCharArrayFromRuneSlice(q.data.bpchararr),
		pgsql.ByteaArrayFromByteSliceSlice(q.data.byteaarr),
		pgsql.BPCharArrayFromByteSlice(q.data.chararr),
		pgsql.CIDRArrayFromIPNetSlice(q.data.cidrarr),
		pgsql.DateArrayFromTimeSlice(q.data.datearr),
		pgsql.DateRangeArrayFromTimeArray2Slice(q.data.daterangearr),
		pgsql.Float4ArrayFromFloat32Slice(q.data.float4arr),
		pgsql.Float8ArrayFromFloat64Slice(q.data.float8arr),
		pgsql.InetArrayFromIPSlice(q.data.inetarr),
		pgsql.Int2ArrayFromInt16Slice(q.data.int2arr),
		pgsql.Int2VectorFromInt16Slice(q.data.int2vector),
		pgsql.Int2VectorArrayFromInt16SliceSlice(q.data.int2vectorarr),
		pgsql.Int4ArrayFromInt32Slice(q.data.int4arr),
		pgsql.Int4RangeArrayFromInt32Array2Slice(q.data.int4rangearr),
		pgsql.Int8ArrayFromInt64Slice(q.data.int8arr),
		pgsql.Int8RangeArrayFromInt64Array2Slice(q.data.int8rangearr),
		pgsql.JSONArrayFromByteSliceSlice(q.data.jsonarr),
		pgsql.JSONArrayFromByteSliceSlice(q.data.jsonbarr),
		pgsql.LineArrayFromFloat64Array3Slice(q.data.linearr),
		pgsql.LsegArrayFromFloat64Array2Array2Slice(q.data.lsegarr),
		pgsql.MACAddrArrayFromHardwareAddrSlice(q.data.macaddrarr),
		pgsql.MACAddr8ArrayFromHardwareAddrSlice(q.data.macaddr8arr),
		pgsql.MoneyArrayFromInt64Slice(q.data.moneyarr),
		pgsql.NumericArrayFromInt64Slice(q.data.numericarr),
		pgsql.NumRangeArrayFromFloat64Array2Slice(q.data.numrangearr),
		pgsql.PathArrayFromFloat64Array2SliceSlice(q.data.patharr),
		pgsql.PointArrayFromFloat64Array2Slice(q.data.pointarr),
		pgsql.PolygonArrayFromFloat64Array2SliceSlice(q.data.polygonarr),
		pgsql.TextArrayFromStringSlice(q.data.textarr),
		pgsql.TimeArrayFromTimeSlice(q.data.timearr),
		pgsql.TimestampArrayFromTimeSlice(q.data.timestamparr),
		pgsql.TimestamptzArrayFromTimeSlice(q.data.timestamptzarr),
		pgsql.TimetzArrayFromTimeSlice(q.data.timetzarr),
		pgsql.TSQueryArrayFromStringSlice(q.data.tsqueryarr),
		pgsql.TsRangeArrayFromTimeArray2Slice(q.data.tsrangearr),
		pgsql.TstzRangeArrayFromTimeArray2Slice(q.data.tstzrangearr),
		pgsql.TSVectorArrayFromStringSliceSlice(q.data.tsvectorarr),
		pgsql.UUIDArrayFromByteArray16Slice(q.data.uuidarr),
		pgsql.VarBitArrayFromBoolSliceSlice(q.data.varbitarr),
		pgsql.VarCharArrayFromStringSlice(q.data.varchararr),
		pgsql.XMLArrayFromByteSliceSlice(q.data.xmlarr),
	)
	return err
}
