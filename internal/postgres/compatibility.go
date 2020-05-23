package postgres

import (
	"github.com/frk/gosql/internal/analysis"
	"github.com/frk/gosql/internal/postgres/oid"
)

type compkey struct {
	oid     oid.OID
	typmod1 bool
}

type compentry struct {
	valuer  string
	scanner string
}

type comptable struct {
	literal2oid map[analysis.LiteralType]map[compkey]compentry
	oid2literal map[compkey]map[analysis.LiteralType]compentry
}

var compatibility = comptable{
	literal2oid: init_literal2oid(),
	oid2literal: init_oid2literal(),
}

func (c comptable) getTypeInfoOIDs(typ analysis.TypeInfo) []oid.OID {
	lit := typ.GenericLiteral()
	keys, ok := compatibility.literal2oid[lit]
	if !ok {
		return nil
	}

	var oids []oid.OID
	for key, ce := range keys {
		if ce.valuer == "" && ce.scanner == "" {
			oids = append([]oid.OID{key.oid}, oids...)
		} else {
			oids = append(oids, key.oid)
		}
	}
	return oids
}

func init_literal2oid() map[analysis.LiteralType]map[compkey]compentry {
	oid2literal := init_oid2literal()
	literal2oid := make(map[analysis.LiteralType]map[compkey]compentry)
	for id, litmap := range oid2literal {
		for lit, comp := range litmap {
			if oidmap, ok := literal2oid[lit]; !ok {
				literal2oid[lit] = map[compkey]compentry{id: comp}
			} else {
				oidmap[id] = comp
			}
		}
	}
	return literal2oid
}

func init_oid2literal() map[compkey]map[analysis.LiteralType]compentry {
	return map[compkey]map[analysis.LiteralType]compentry{
		{oid: oid.Bit, typmod1: true}: {
			analysis.LiteralBool:      {valuer: "BitFromBool"},
			analysis.LiteralUint:      {},
			analysis.LiteralUint8:     {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.BitArr, typmod1: true}: {
			analysis.LiteralBoolSlice:  {valuer: "BitArrayFromBoolSlice", scanner: "BitArrayToBoolSlice"},
			analysis.LiteralUintSlice:  {valuer: "BitArrayFromUintSlice", scanner: "BitArrayToUintSlice"},
			analysis.LiteralUint8Slice: {valuer: "BitArrayFromUint8Slice", scanner: "BitArrayToUint8Slice"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.BPChar, typmod1: true}: {
			analysis.LiteralByte:      {valuer: "BPCharFromByte", scanner: "BPCharToByte"},
			analysis.LiteralRune:      {valuer: "BPCharFromRune", scanner: "BPCharToRune"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.BPCharArr, typmod1: true}: {
			analysis.LiteralRuneSlice:   {valuer: "BPCharArrayFromRuneSlice", scanner: "BPCharArrayToRuneSlice"},
			analysis.LiteralStringSlice: {valuer: "BPCharArrayFromStringSlice", scanner: "BPCharArrayToStringSlice"},
			analysis.LiteralString:      {valuer: "BPCharArrayFromString", scanner: "BPCharArrayToString"},
			analysis.LiteralByteSlice:   {valuer: "BPCharArrayFromByteSlice", scanner: "BPCharArrayToByteSlice"},
		},
		{oid: oid.Char, typmod1: true}: {
			analysis.LiteralByte:      {valuer: "CharFromByte", scanner: "CharToByte"},
			analysis.LiteralRune:      {valuer: "CharFromRune", scanner: "CharToRune"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.CharArr, typmod1: true}: {
			analysis.LiteralRuneSlice:   {valuer: "CharArrayFromRuneSlice", scanner: "CharArrayToRuneSlice"},
			analysis.LiteralStringSlice: {valuer: "CharArrayFromStringSlice", scanner: "CharArrayToStringSlice"},
			analysis.LiteralString:      {valuer: "CharArrayFromString", scanner: "CharArrayToString"},
			analysis.LiteralByteSlice:   {valuer: "CharArrayFromByteSlice", scanner: "CharArrayToByteSlice"},
		},

		////////////////////////////////////////////////////////////////
		{oid: oid.Bool}: {
			analysis.LiteralBool:      {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.BoolArr}: {
			analysis.LiteralBoolSlice: {valuer: "BoolArrayFromBoolSlice", scanner: "BoolArrayToBoolSlice"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Box}: {
			analysis.LiteralFloat64Array2Array2: {valuer: "BoxFromFloat64Array2Array2", scanner: "BoxToFloat64Array2Array2"},
			analysis.LiteralString:              {},
			analysis.LiteralByteSlice:           {},
		},
		{oid: oid.BoxArr}: {
			analysis.LiteralFloat64Array2Array2Slice: {valuer: "BoxArrayFromFloat64Array2Array2Slice", scanner: "BoxArrayToFloat64Array2Array2Slice"},
			analysis.LiteralString:                   {},
			analysis.LiteralByteSlice:                {},
		},
		{oid: oid.Bytea}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.ByteaArr}: {
			analysis.LiteralStringSlice:    {valuer: "ByteaArrayFromStringSlice", scanner: "ByteaArrayToStringSlice"},
			analysis.LiteralByteSliceSlice: {valuer: "ByteaArrayFromByteSliceSlice", scanner: "ByteaArrayToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.CIDR}: {
			analysis.LiteralIPNet:     {valuer: "CIDRFromIPNet", scanner: "CIDRToIPNet"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.CIDRArr}: {
			analysis.LiteralIPNetSlice: {valuer: "CIDRArrayFromIPNetSlice", scanner: "CIDRArrayToIPNetSlice"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.Circle}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.CircleArr}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Date}: {
			analysis.LiteralTime:      {scanner: "DateToTime"},
			analysis.LiteralString:    {scanner: "DateToString"},
			analysis.LiteralByteSlice: {scanner: "DateToByteSlice"},
		},
		{oid: oid.DateArr}: {
			analysis.LiteralTimeSlice: {valuer: "DateArrayFromTimeSlice", scanner: "DateArrayToTimeSlice"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.DateRange}: {
			analysis.LiteralTimeArray2: {valuer: "DateRangeFromTimeArray2", scanner: "DateRangeToTimeArray2"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.DateRangeArr}: {
			analysis.LiteralTimeArray2Slice: {valuer: "DateRangeArrayFromTimeArray2Slice", scanner: "DateRangeArrayToTimeArray2Slice"},
			analysis.LiteralString:          {},
			analysis.LiteralByteSlice:       {},
		},
		{oid: oid.Float4}: {
			analysis.LiteralFloat32:   {},
			analysis.LiteralFloat64:   {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Float4Arr}: {
			analysis.LiteralFloat32Slice: {valuer: "Float4ArrayFromFloat32Slice", scanner: "Float4ArrayToFloat32Slice"},
			analysis.LiteralFloat64Slice: {valuer: "Float4ArrayFromFloat64Slice", scanner: "Float4ArrayToFloat64Slice"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.Float8}: {
			analysis.LiteralFloat32:   {},
			analysis.LiteralFloat64:   {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Float8Arr}: {
			analysis.LiteralFloat32Slice: {valuer: "Float8ArrayFromFloat32Slice", scanner: "Float8ArrayToFloat32Slice"},
			analysis.LiteralFloat64Slice: {valuer: "Float8ArrayFromFloat64Slice", scanner: "Float8ArrayToFloat64Slice"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		/*
			{oid: oid.HStore}: {
				analysis.LiteralStringMap:     {valuer: "HStoreFromStringMap", scanner: "HStoreToStringMap"},
				analysis.LiteralStringPtrMap:  {valuer: "HStoreFromStringPtrMap", scanner: "HStoreToStringPtrMap"},
				analysis.LiteralNullStringMap: {valuer: "HStoreFromNullStringMap", scanner: "HStoreToNullStringMap"},
				analysis.LiteralString:        {},
				analysis.LiteralByteSlice:     {},
			},
			{oid: oid.HStoreArr}: {
				analysis.LiteralStringMapSlice:     {valuer: "HStoreArrayFromStringMapSlice", scanner: "HStoreArrayToStringMapSlice"},
				analysis.LiteralStringPtrMapSlice:  {valuer: "HStoreArrayFromStringPtrMapSlice", scanner: "HStoreArrayToStringPtrMapSlice"},
				analysis.LiteralNullStringMapSlice: {valuer: "HStoreArrayFromNullStringMapSlice", scanner: "HStoreArrayToNullStringMapSlice"},
				analysis.LiteralString:             {},
				analysis.LiteralByteSlice:          {},
			},
		*/
		{oid: oid.Inet}: {
			analysis.LiteralIP:        {valuer: "InetFromIP", scanner: "InetToIP"},
			analysis.LiteralIPNet:     {valuer: "InetFromIPNet", scanner: "InetToIPNet"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.InetArr}: {
			analysis.LiteralIPSlice:    {valuer: "InetArrayFromIPSlice", scanner: "InetArrayFromIPSlice"},
			analysis.LiteralIPNetSlice: {valuer: "InetArrayFromIPNetSlice", scanner: "InetArrayFromIPNetSlice"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.Int2}: {
			analysis.LiteralInt:       {},
			analysis.LiteralInt8:      {},
			analysis.LiteralInt16:     {},
			analysis.LiteralInt32:     {},
			analysis.LiteralInt64:     {},
			analysis.LiteralUint:      {},
			analysis.LiteralUint8:     {},
			analysis.LiteralUint16:    {},
			analysis.LiteralUint32:    {},
			analysis.LiteralUint64:    {},
			analysis.LiteralFloat32:   {},
			analysis.LiteralFloat64:   {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Int2Arr}: {
			analysis.LiteralIntSlice:     {valuer: "Int2ArrayFromIntSlice", scanner: "Int2ArrayToIntSlice"},
			analysis.LiteralInt8Slice:    {valuer: "Int2ArrayFromInt8Slice", scanner: "Int2ArrayToInt8Slice"},
			analysis.LiteralInt16Slice:   {valuer: "Int2ArrayFromInt16Slice", scanner: "Int2ArrayToInt16Slice"},
			analysis.LiteralInt32Slice:   {valuer: "Int2ArrayFromInt32Slice", scanner: "Int2ArrayToInt32Slice"},
			analysis.LiteralInt64Slice:   {valuer: "Int2ArrayFromInt64Slice", scanner: "Int2ArrayToInt64Slice"},
			analysis.LiteralUintSlice:    {valuer: "Int2ArrayFromUintSlice", scanner: "Int2ArrayToUintSlice"},
			analysis.LiteralUint8Slice:   {valuer: "Int2ArrayFromUint8Slice", scanner: "Int2ArrayToUint8Slice"},
			analysis.LiteralUint16Slice:  {valuer: "Int2ArrayFromUint16Slice", scanner: "Int2ArrayToUint16Slice"},
			analysis.LiteralUint32Slice:  {valuer: "Int2ArrayFromUint32Slice", scanner: "Int2ArrayToUint32Slice"},
			analysis.LiteralUint64Slice:  {valuer: "Int2ArrayFromUint64Slice", scanner: "Int2ArrayToUint64Slice"},
			analysis.LiteralFloat32Slice: {valuer: "Int2ArrayFromFloat32Slice", scanner: "Int2ArrayToFloat32Slice"},
			analysis.LiteralFloat64Slice: {valuer: "Int2ArrayFromFloat64Slice", scanner: "Int2ArrayToFloat64Slice"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.Int2Vector}: {
			analysis.LiteralIntSlice:     {valuer: "Int2VectorFromIntSlice", scanner: "Int2VectorToIntSlice"},
			analysis.LiteralInt8Slice:    {valuer: "Int2VectorFromInt8Slice", scanner: "Int2VectorToInt8Slice"},
			analysis.LiteralInt16Slice:   {valuer: "Int2VectorFromInt16Slice", scanner: "Int2VectorToInt16Slice"},
			analysis.LiteralInt32Slice:   {valuer: "Int2VectorFromInt32Slice", scanner: "Int2VectorToInt32Slice"},
			analysis.LiteralInt64Slice:   {valuer: "Int2VectorFromInt64Slice", scanner: "Int2VectorToInt64Slice"},
			analysis.LiteralUintSlice:    {valuer: "Int2VectorFromUintSlice", scanner: "Int2VectorToUintSlice"},
			analysis.LiteralUint8Slice:   {valuer: "Int2VectorFromUint8Slice", scanner: "Int2VectorToUint8Slice"},
			analysis.LiteralUint16Slice:  {valuer: "Int2VectorFromUint16Slice", scanner: "Int2VectorToUint16Slice"},
			analysis.LiteralUint32Slice:  {valuer: "Int2VectorFromUint32Slice", scanner: "Int2VectorToUint32Slice"},
			analysis.LiteralUint64Slice:  {valuer: "Int2VectorFromUint64Slice", scanner: "Int2VectorToUint64Slice"},
			analysis.LiteralFloat32Slice: {valuer: "Int2VectorFromFloat32Slice", scanner: "Int2VectorToFloat32Slice"},
			analysis.LiteralFloat64Slice: {valuer: "Int2VectorFromFloat64Slice", scanner: "Int2VectorToFloat64Slice"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.Int2VectorArr}: {
			analysis.LiteralIntSliceSlice:     {valuer: "Int2VectorArrayFromIntSliceSlice", scanner: "Int2VectorArrayToIntSliceSlice"},
			analysis.LiteralInt8SliceSlice:    {valuer: "Int2VectorArrayFromInt8SliceSlice", scanner: "Int2VectorArrayToInt8SliceSlice"},
			analysis.LiteralInt16SliceSlice:   {valuer: "Int2VectorArrayFromInt16SliceSlice", scanner: "Int2VectorArrayToInt16SliceSlice"},
			analysis.LiteralInt32SliceSlice:   {valuer: "Int2VectorArrayFromInt32SliceSlice", scanner: "Int2VectorArrayToInt32SliceSlice"},
			analysis.LiteralInt64SliceSlice:   {valuer: "Int2VectorArrayFromInt64SliceSlice", scanner: "Int2VectorArrayToInt64SliceSlice"},
			analysis.LiteralUintSliceSlice:    {valuer: "Int2VectorArrayFromUintSliceSlice", scanner: "Int2VectorArrayToUintSliceSlice"},
			analysis.LiteralUint8SliceSlice:   {valuer: "Int2VectorArrayFromUint8SliceSlice", scanner: "Int2VectorArrayToUint8SliceSlice"},
			analysis.LiteralUint16SliceSlice:  {valuer: "Int2VectorArrayFromUint16SliceSlice", scanner: "Int2VectorArrayToUint16SliceSlice"},
			analysis.LiteralUint32SliceSlice:  {valuer: "Int2VectorArrayFromUint32SliceSlice", scanner: "Int2VectorArrayToUint32SliceSlice"},
			analysis.LiteralUint64SliceSlice:  {valuer: "Int2VectorArrayFromUint64SliceSlice", scanner: "Int2VectorArrayToUint64SliceSlice"},
			analysis.LiteralFloat32SliceSlice: {valuer: "Int2VectorArrayFromFloat32SliceSlice", scanner: "Int2VectorArrayToFloat32SliceSlice"},
			analysis.LiteralFloat64SliceSlice: {valuer: "Int2VectorArrayFromFloat64SliceSlice", scanner: "Int2VectorArrayToFloat64SliceSlice"},
			analysis.LiteralString:            {},
			analysis.LiteralByteSlice:         {},
		},
		{oid: oid.Int4}: {
			analysis.LiteralInt:       {},
			analysis.LiteralInt8:      {},
			analysis.LiteralInt16:     {},
			analysis.LiteralInt32:     {},
			analysis.LiteralInt64:     {},
			analysis.LiteralUint:      {},
			analysis.LiteralUint8:     {},
			analysis.LiteralUint16:    {},
			analysis.LiteralUint32:    {},
			analysis.LiteralUint64:    {},
			analysis.LiteralFloat32:   {},
			analysis.LiteralFloat64:   {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Int4Arr}: {
			analysis.LiteralIntSlice:     {valuer: "Int4ArrayFromIntSlice", scanner: "Int4ArrayToIntSlice"},
			analysis.LiteralInt8Slice:    {valuer: "Int4ArrayFromInt8Slice", scanner: "Int4ArrayToInt8Slice"},
			analysis.LiteralInt16Slice:   {valuer: "Int4ArrayFromInt16Slice", scanner: "Int4ArrayToInt16Slice"},
			analysis.LiteralInt32Slice:   {valuer: "Int4ArrayFromInt32Slice", scanner: "Int4ArrayToInt32Slice"},
			analysis.LiteralInt64Slice:   {valuer: "Int4ArrayFromInt64Slice", scanner: "Int4ArrayToInt64Slice"},
			analysis.LiteralUintSlice:    {valuer: "Int4ArrayFromUintSlice", scanner: "Int4ArrayToUintSlice"},
			analysis.LiteralUint8Slice:   {valuer: "Int4ArrayFromUint8Slice", scanner: "Int4ArrayToUint8Slice"},
			analysis.LiteralUint16Slice:  {valuer: "Int4ArrayFromUint16Slice", scanner: "Int4ArrayToUint16Slice"},
			analysis.LiteralUint32Slice:  {valuer: "Int4ArrayFromUint32Slice", scanner: "Int4ArrayToUint32Slice"},
			analysis.LiteralUint64Slice:  {valuer: "Int4ArrayFromUint64Slice", scanner: "Int4ArrayToUint64Slice"},
			analysis.LiteralFloat32Slice: {valuer: "Int4ArrayFromFloat32Slice", scanner: "Int4ArrayToFloat32Slice"},
			analysis.LiteralFloat64Slice: {valuer: "Int4ArrayFromFloat64Slice", scanner: "Int4ArrayToFloat64Slice"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.Int4Range}: {
			analysis.LiteralIntArray2:     {valuer: "Int4RangeFromIntArray2", scanner: "Int4RangeToIntArray2"},
			analysis.LiteralInt8Array2:    {valuer: "Int4RangeFromInt8Array2", scanner: "Int4RangeToInt8Array2"},
			analysis.LiteralInt16Array2:   {valuer: "Int4RangeFromInt16Array2", scanner: "Int4RangeToInt16Array2"},
			analysis.LiteralInt32Array2:   {valuer: "Int4RangeFromInt32Array2", scanner: "Int4RangeToInt32Array2"},
			analysis.LiteralInt64Array2:   {valuer: "Int4RangeFromInt64Array2", scanner: "Int4RangeToInt64Array2"},
			analysis.LiteralUintArray2:    {valuer: "Int4RangeFromUintArray2", scanner: "Int4RangeToUintArray2"},
			analysis.LiteralUint8Array2:   {valuer: "Int4RangeFromUint8Array2", scanner: "Int4RangeToUint8Array2"},
			analysis.LiteralUint16Array2:  {valuer: "Int4RangeFromUint16Array2", scanner: "Int4RangeToUint16Array2"},
			analysis.LiteralUint32Array2:  {valuer: "Int4RangeFromUint32Array2", scanner: "Int4RangeToUint32Array2"},
			analysis.LiteralUint64Array2:  {valuer: "Int4RangeFromUint64Array2", scanner: "Int4RangeToUint64Array2"},
			analysis.LiteralFloat32Array2: {valuer: "Int4RangeFromFloat32Array2", scanner: "Int4RangeToFloat32Array2"},
			analysis.LiteralFloat64Array2: {valuer: "Int4RangeFromFloat64Array2", scanner: "Int4RangeToFloat64Array2"},
			analysis.LiteralString:        {},
			analysis.LiteralByteSlice:     {},
		},
		{oid: oid.Int4RangeArr}: {
			analysis.LiteralIntArray2Slice:     {valuer: "Int4RangeArrayFromIntArray2Slice", scanner: "Int4RangeArrayToIntArray2Slice"},
			analysis.LiteralInt8Array2Slice:    {valuer: "Int4RangeArrayFromInt8Array2Slice", scanner: "Int4RangeArrayToInt8Array2Slice"},
			analysis.LiteralInt16Array2Slice:   {valuer: "Int4RangeArrayFromInt16Array2Slice", scanner: "Int4RangeArrayToInt16Array2Slice"},
			analysis.LiteralInt32Array2Slice:   {valuer: "Int4RangeArrayFromInt32Array2Slice", scanner: "Int4RangeArrayToInt32Array2Slice"},
			analysis.LiteralInt64Array2Slice:   {valuer: "Int4RangeArrayFromInt64Array2Slice", scanner: "Int4RangeArrayToInt64Array2Slice"},
			analysis.LiteralUintArray2Slice:    {valuer: "Int4RangeArrayFromUintArray2Slice", scanner: "Int4RangeArrayToUintArray2Slice"},
			analysis.LiteralUint8Array2Slice:   {valuer: "Int4RangeArrayFromUint8Array2Slice", scanner: "Int4RangeArrayToUint8Array2Slice"},
			analysis.LiteralUint16Array2Slice:  {valuer: "Int4RangeArrayFromUint16Array2Slice", scanner: "Int4RangeArrayToUint16Array2Slice"},
			analysis.LiteralUint32Array2Slice:  {valuer: "Int4RangeArrayFromUint32Array2Slice", scanner: "Int4RangeArrayToUint32Array2Slice"},
			analysis.LiteralUint64Array2Slice:  {valuer: "Int4RangeArrayFromUint64Array2Slice", scanner: "Int4RangeArrayToUint64Array2Slice"},
			analysis.LiteralFloat32Array2Slice: {valuer: "Int4RangeArrayFromFloat32Array2Slice", scanner: "Int4RangeArrayToFloat32Array2Slice"},
			analysis.LiteralFloat64Array2Slice: {valuer: "Int4RangeArrayFromFloat64Array2Slice", scanner: "Int4RangeArrayToFloat64Array2Slice"},
			analysis.LiteralString:             {},
			analysis.LiteralByteSlice:          {},
		},
		{oid: oid.Int8}: {
			analysis.LiteralInt:       {},
			analysis.LiteralInt8:      {},
			analysis.LiteralInt16:     {},
			analysis.LiteralInt32:     {},
			analysis.LiteralInt64:     {},
			analysis.LiteralUint:      {},
			analysis.LiteralUint8:     {},
			analysis.LiteralUint16:    {},
			analysis.LiteralUint32:    {},
			analysis.LiteralUint64:    {},
			analysis.LiteralFloat32:   {},
			analysis.LiteralFloat64:   {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Int8Arr}: {
			analysis.LiteralIntSlice:     {valuer: "Int8ArrayFromIntSlice", scanner: "Int8ArrayToIntSlice"},
			analysis.LiteralInt8Slice:    {valuer: "Int8ArrayFromInt8Slice", scanner: "Int8ArrayToInt8Slice"},
			analysis.LiteralInt16Slice:   {valuer: "Int8ArrayFromInt16Slice", scanner: "Int8ArrayToInt16Slice"},
			analysis.LiteralInt32Slice:   {valuer: "Int8ArrayFromInt32Slice", scanner: "Int8ArrayToInt32Slice"},
			analysis.LiteralInt64Slice:   {valuer: "Int8ArrayFromInt64Slice", scanner: "Int8ArrayToInt64Slice"},
			analysis.LiteralUintSlice:    {valuer: "Int8ArrayFromUintSlice", scanner: "Int8ArrayToUintSlice"},
			analysis.LiteralUint8Slice:   {valuer: "Int8ArrayFromUint8Slice", scanner: "Int8ArrayToUint8Slice"},
			analysis.LiteralUint16Slice:  {valuer: "Int8ArrayFromUint16Slice", scanner: "Int8ArrayToUint16Slice"},
			analysis.LiteralUint32Slice:  {valuer: "Int8ArrayFromUint32Slice", scanner: "Int8ArrayToUint32Slice"},
			analysis.LiteralUint64Slice:  {valuer: "Int8ArrayFromUint64Slice", scanner: "Int8ArrayToUint64Slice"},
			analysis.LiteralFloat32Slice: {valuer: "Int8ArrayFromFloat32Slice", scanner: "Int8ArrayToFloat32Slice"},
			analysis.LiteralFloat64Slice: {valuer: "Int8ArrayFromFloat64Slice", scanner: "Int8ArrayToFloat64Slice"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.Int8Range}: {
			analysis.LiteralIntArray2:     {valuer: "Int8RangeFromIntArray2", scanner: "Int8RangeToIntArray2"},
			analysis.LiteralInt8Array2:    {valuer: "Int8RangeFromInt8Array2", scanner: "Int8RangeToInt8Array2"},
			analysis.LiteralInt16Array2:   {valuer: "Int8RangeFromInt16Array2", scanner: "Int8RangeToInt16Array2"},
			analysis.LiteralInt32Array2:   {valuer: "Int8RangeFromInt32Array2", scanner: "Int8RangeToInt32Array2"},
			analysis.LiteralInt64Array2:   {valuer: "Int8RangeFromInt64Array2", scanner: "Int8RangeToInt64Array2"},
			analysis.LiteralUintArray2:    {valuer: "Int8RangeFromUintArray2", scanner: "Int8RangeToUintArray2"},
			analysis.LiteralUint8Array2:   {valuer: "Int8RangeFromUint8Array2", scanner: "Int8RangeToUint8Array2"},
			analysis.LiteralUint16Array2:  {valuer: "Int8RangeFromUint16Array2", scanner: "Int8RangeToUint16Array2"},
			analysis.LiteralUint32Array2:  {valuer: "Int8RangeFromUint32Array2", scanner: "Int8RangeToUint32Array2"},
			analysis.LiteralUint64Array2:  {valuer: "Int8RangeFromUint64Array2", scanner: "Int8RangeToUint64Array2"},
			analysis.LiteralFloat32Array2: {valuer: "Int8RangeFromFloat32Array2", scanner: "Int8RangeToFloat32Array2"},
			analysis.LiteralFloat64Array2: {valuer: "Int8RangeFromFloat64Array2", scanner: "Int8RangeToFloat64Array2"},
			analysis.LiteralString:        {},
			analysis.LiteralByteSlice:     {},
		},
		{oid: oid.Int8RangeArr}: {
			analysis.LiteralIntArray2Slice:     {valuer: "Int8RangeArrayFromIntArray2Slice", scanner: "Int8RangeArrayToIntArray2Slice"},
			analysis.LiteralInt8Array2Slice:    {valuer: "Int8RangeArrayFromInt8Array2Slice", scanner: "Int8RangeArrayToInt8Array2Slice"},
			analysis.LiteralInt16Array2Slice:   {valuer: "Int8RangeArrayFromInt16Array2Slice", scanner: "Int8RangeArrayToInt16Array2Slice"},
			analysis.LiteralInt32Array2Slice:   {valuer: "Int8RangeArrayFromInt32Array2Slice", scanner: "Int8RangeArrayToInt32Array2Slice"},
			analysis.LiteralInt64Array2Slice:   {valuer: "Int8RangeArrayFromInt64Array2Slice", scanner: "Int8RangeArrayToInt64Array2Slice"},
			analysis.LiteralUintArray2Slice:    {valuer: "Int8RangeArrayFromUintArray2Slice", scanner: "Int8RangeArrayToUintArray2Slice"},
			analysis.LiteralUint8Array2Slice:   {valuer: "Int8RangeArrayFromUint8Array2Slice", scanner: "Int8RangeArrayToUint8Array2Slice"},
			analysis.LiteralUint16Array2Slice:  {valuer: "Int8RangeArrayFromUint16Array2Slice", scanner: "Int8RangeArrayToUint16Array2Slice"},
			analysis.LiteralUint32Array2Slice:  {valuer: "Int8RangeArrayFromUint32Array2Slice", scanner: "Int8RangeArrayToUint32Array2Slice"},
			analysis.LiteralUint64Array2Slice:  {valuer: "Int8RangeArrayFromUint64Array2Slice", scanner: "Int8RangeArrayToUint64Array2Slice"},
			analysis.LiteralFloat32Array2Slice: {valuer: "Int8RangeArrayFromFloat32Array2Slice", scanner: "Int8RangeArrayToFloat32Array2Slice"},
			analysis.LiteralFloat64Array2Slice: {valuer: "Int8RangeArrayFromFloat64Array2Slice", scanner: "Int8RangeArrayToFloat64Array2Slice"},
			analysis.LiteralString:             {},
			analysis.LiteralByteSlice:          {},
		},
		{oid: oid.Interval}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.IntervalArr}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.JSON}: {
			analysis.LiteralEmptyInterface: {valuer: "JSON", scanner: "JSON"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.JSONArr}: {
			analysis.LiteralByteSliceSlice: {valuer: "JSONArrayFromByteSliceSlice", scanner: "JSONArrayToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.JSONB}: {
			analysis.LiteralEmptyInterface: {valuer: "JSON", scanner: "JSON"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.JSONBArr}: {
			analysis.LiteralByteSliceSlice: {valuer: "JSONArrayFromByteSliceSlice", scanner: "JSONArrayToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.Line}: {
			analysis.LiteralFloat64Array3: {valuer: "LineFromFloat64Array3", scanner: "LineToFloat64Array3"},
			analysis.LiteralString:        {},
			analysis.LiteralByteSlice:     {},
		},
		{oid: oid.LineArr}: {
			analysis.LiteralFloat64Array3Slice: {valuer: "LineArrayFromFloat64Array3Slice", scanner: "LineArrayToFloat64Array3Slice"},
			analysis.LiteralString:             {},
			analysis.LiteralByteSlice:          {},
		},
		{oid: oid.LSeg}: {
			analysis.LiteralFloat64Array2Array2: {valuer: "LsegFromFloat64Array2Array2", scanner: "LsegToFloat64Array2Array2"},
			analysis.LiteralString:              {},
			analysis.LiteralByteSlice:           {},
		},
		{oid: oid.LSegArr}: {
			analysis.LiteralFloat64Array2Array2Slice: {valuer: "LsegArrayFromFloat64Array2Array2Slice", scanner: "LsegArrayToFloat64Array2Array2Slice"},
			analysis.LiteralString:                   {},
			analysis.LiteralByteSlice:                {},
		},
		{oid: oid.MACAddr}: {
			analysis.LiteralHardwareAddr: {valuer: "MACAddrFromHardwareAddr", scanner: "MACAddrToHardwareAddr"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.MACAddrArr}: {
			analysis.LiteralHardwareAddrSlice: {valuer: "MACAddrArrayFromHardwareAddrSlice", scanner: "MACAddrArrayToHardwareAddrSlice"},
			analysis.LiteralString:            {},
			analysis.LiteralByteSlice:         {},
		},
		{oid: oid.MACAddr8}: {
			analysis.LiteralHardwareAddr: {valuer: "MACAddr8FromHardwareAddr", scanner: "MACAddr8ToHardwareAddr"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.MACAddr8Arr}: {
			analysis.LiteralHardwareAddrSlice: {valuer: "MACAddr8ArrayFromHardwareAddrSlice", scanner: "MACAddr8ArrayToHardwareAddrSlice"},
			analysis.LiteralString:            {},
			analysis.LiteralByteSlice:         {},
		},
		{oid: oid.Money}: {
			analysis.LiteralInt64:     {valuer: "MoneyFromInt64", scanner: "MoneyToInt64"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.MoneyArr}: {
			analysis.LiteralInt64Slice: {valuer: "MoneyArrayFromInt64Slice", scanner: "MoneyArrayToInt64Slice"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.Numeric}: {
			analysis.LiteralInt:       {},
			analysis.LiteralInt8:      {},
			analysis.LiteralInt16:     {},
			analysis.LiteralInt32:     {},
			analysis.LiteralInt64:     {},
			analysis.LiteralUint:      {},
			analysis.LiteralUint8:     {},
			analysis.LiteralUint16:    {},
			analysis.LiteralUint32:    {},
			analysis.LiteralUint64:    {},
			analysis.LiteralFloat32:   {},
			analysis.LiteralFloat64:   {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.NumericArr}: {
			analysis.LiteralIntSlice:     {valuer: "NumericArrayFromIntSlice", scanner: "NumericArrayToIntSlice"},
			analysis.LiteralInt8Slice:    {valuer: "NumericArrayFromInt8Slice", scanner: "NumericArrayToInt8Slice"},
			analysis.LiteralInt16Slice:   {valuer: "NumericArrayFromInt16Slice", scanner: "NumericArrayToInt16Slice"},
			analysis.LiteralInt32Slice:   {valuer: "NumericArrayFromInt32Slice", scanner: "NumericArrayToInt32Slice"},
			analysis.LiteralInt64Slice:   {valuer: "NumericArrayFromInt64Slice", scanner: "NumericArrayToInt64Slice"},
			analysis.LiteralUintSlice:    {valuer: "NumericArrayFromUintSlice", scanner: "NumericArrayToUintSlice"},
			analysis.LiteralUint8Slice:   {valuer: "NumericArrayFromUint8Slice", scanner: "NumericArrayToUint8Slice"},
			analysis.LiteralUint16Slice:  {valuer: "NumericArrayFromUint16Slice", scanner: "NumericArrayToUint16Slice"},
			analysis.LiteralUint32Slice:  {valuer: "NumericArrayFromUint32Slice", scanner: "NumericArrayToUint32Slice"},
			analysis.LiteralUint64Slice:  {valuer: "NumericArrayFromUint64Slice", scanner: "NumericArrayToUint64Slice"},
			analysis.LiteralFloat32Slice: {valuer: "NumericArrayFromFloat32Slice", scanner: "NumericArrayToFloat32Slice"},
			analysis.LiteralFloat64Slice: {valuer: "NumericArrayFromFloat64Slice", scanner: "NumericArrayToFloat64Slice"},
			analysis.LiteralString:       {},
			analysis.LiteralByteSlice:    {},
		},
		{oid: oid.NumRange}: {
			analysis.LiteralIntArray2:     {valuer: "NumRangeFromIntArray2", scanner: "NumRangeToIntArray2"},
			analysis.LiteralInt8Array2:    {valuer: "NumRangeFromInt8Array2", scanner: "NumRangeToInt8Array2"},
			analysis.LiteralInt16Array2:   {valuer: "NumRangeFromInt16Array2", scanner: "NumRangeToInt16Array2"},
			analysis.LiteralInt32Array2:   {valuer: "NumRangeFromInt32Array2", scanner: "NumRangeToInt32Array2"},
			analysis.LiteralInt64Array2:   {valuer: "NumRangeFromInt64Array2", scanner: "NumRangeToInt64Array2"},
			analysis.LiteralUintArray2:    {valuer: "NumRangeFromUintArray2", scanner: "NumRangeToUintArray2"},
			analysis.LiteralUint8Array2:   {valuer: "NumRangeFromUint8Array2", scanner: "NumRangeToUint8Array2"},
			analysis.LiteralUint16Array2:  {valuer: "NumRangeFromUint16Array2", scanner: "NumRangeToUint16Array2"},
			analysis.LiteralUint32Array2:  {valuer: "NumRangeFromUint32Array2", scanner: "NumRangeToUint32Array2"},
			analysis.LiteralUint64Array2:  {valuer: "NumRangeFromUint64Array2", scanner: "NumRangeToUint64Array2"},
			analysis.LiteralFloat32Array2: {valuer: "NumRangeFromFloat32Array2", scanner: "NumRangeToFloat32Array2"},
			analysis.LiteralFloat64Array2: {valuer: "NumRangeFromFloat64Array2", scanner: "NumRangeToFloat64Array2"},
			analysis.LiteralString:        {},
			analysis.LiteralByteSlice:     {},
		},
		{oid: oid.NumRangeArr}: {
			analysis.LiteralIntArray2Slice:     {valuer: "NumRangeArrayFromIntArray2Slice", scanner: "NumRangeArrayToIntArray2Slice"},
			analysis.LiteralInt8Array2Slice:    {valuer: "NumRangeArrayFromInt8Array2Slice", scanner: "NumRangeArrayToInt8Array2Slice"},
			analysis.LiteralInt16Array2Slice:   {valuer: "NumRangeArrayFromInt16Array2Slice", scanner: "NumRangeArrayToInt16Array2Slice"},
			analysis.LiteralInt32Array2Slice:   {valuer: "NumRangeArrayFromInt32Array2Slice", scanner: "NumRangeArrayToInt32Array2Slice"},
			analysis.LiteralInt64Array2Slice:   {valuer: "NumRangeArrayFromInt64Array2Slice", scanner: "NumRangeArrayToInt64Array2Slice"},
			analysis.LiteralUintArray2Slice:    {valuer: "NumRangeArrayFromUintArray2Slice", scanner: "NumRangeArrayToUintArray2Slice"},
			analysis.LiteralUint8Array2Slice:   {valuer: "NumRangeArrayFromUint8Array2Slice", scanner: "NumRangeArrayToUint8Array2Slice"},
			analysis.LiteralUint16Array2Slice:  {valuer: "NumRangeArrayFromUint16Array2Slice", scanner: "NumRangeArrayToUint16Array2Slice"},
			analysis.LiteralUint32Array2Slice:  {valuer: "NumRangeArrayFromUint32Array2Slice", scanner: "NumRangeArrayToUint32Array2Slice"},
			analysis.LiteralUint64Array2Slice:  {valuer: "NumRangeArrayFromUint64Array2Slice", scanner: "NumRangeArrayToUint64Array2Slice"},
			analysis.LiteralFloat32Array2Slice: {valuer: "NumRangeArrayFromFloat32Array2Slice", scanner: "NumRangeArrayToFloat32Array2Slice"},
			analysis.LiteralFloat64Array2Slice: {valuer: "NumRangeArrayFromFloat64Array2Slice", scanner: "NumRangeArrayToFloat64Array2Slice"},
			analysis.LiteralString:             {},
			analysis.LiteralByteSlice:          {},
		},
		{oid: oid.Path}: {
			analysis.LiteralFloat64Array2Slice: {valuer: "PathFromFloat64Array2Slice", scanner: "PathToFloat64Array2Slice"},
			analysis.LiteralString:             {},
			analysis.LiteralByteSlice:          {},
		},
		{oid: oid.PathArr}: {
			analysis.LiteralFloat64Array2SliceSlice: {valuer: "PathArrayFromFloat64Array2SliceSlice", scanner: "PathArrayToFloat64Array2SliceSlice"},
			analysis.LiteralString:                  {},
			analysis.LiteralByteSlice:               {},
		},
		{oid: oid.Point}: {
			analysis.LiteralFloat64Array2: {valuer: "PointFromFloat64Array2", scanner: "PointToFloat64Array2"},
			analysis.LiteralString:        {},
			analysis.LiteralByteSlice:     {},
		},
		{oid: oid.PointArr}: {
			analysis.LiteralFloat64Array2Slice: {valuer: "PointArrayFromFloat64Array2Slice", scanner: "PointArrayToFloat64Array2Slice"},
			analysis.LiteralString:             {},
			analysis.LiteralByteSlice:          {},
		},
		{oid: oid.Polygon}: {
			analysis.LiteralFloat64Array2Slice: {valuer: "PolygonFromFloat64Array2Slice", scanner: "PolygonToFloat64Array2Slice"},
			analysis.LiteralString:             {},
			analysis.LiteralByteSlice:          {},
		},
		{oid: oid.PolygonArr}: {
			analysis.LiteralFloat64Array2SliceSlice: {valuer: "PolygonArrayFromFloat64Array2SliceSlice", scanner: "PolygonArrayToFloat64Array2SliceSlice"},
			analysis.LiteralString:                  {},
			analysis.LiteralByteSlice:               {},
		},
		{oid: oid.Text}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.TextArr}: {
			analysis.LiteralStringSlice:    {valuer: "TextArrayFromStringSlice", scanner: "TextArrayToStringSlice"},
			analysis.LiteralByteSliceSlice: {valuer: "TextArrayFromByteSliceSlice", scanner: "TextArrayToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.Time}: {
			analysis.LiteralTime:      {},
			analysis.LiteralString:    {scanner: "TimeToString"},
			analysis.LiteralByteSlice: {scanner: "TimeToByteSlice"},
		},
		{oid: oid.TimeArr}: {
			analysis.LiteralTimeSlice: {valuer: "TimeArrayFromTimeSlice", scanner: "TimeArrayToTimeSlice"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Timestamp}: {
			analysis.LiteralTime:      {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.TimestampArr}: {
			analysis.LiteralTimeSlice: {valuer: "TimestampArrayFromTimeSlice", scanner: "TimestampArrayToTimeSlice"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Timestamptz}: {
			analysis.LiteralTime:      {},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.TimestamptzArr}: {
			analysis.LiteralTimeSlice: {valuer: "TimestamptzArrayFromTimeSlice", scanner: "TimestamptzArrayToTimeSlice"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.Timetz}: {
			analysis.LiteralTime:      {},
			analysis.LiteralString:    {scanner: "TimetzToString"},
			analysis.LiteralByteSlice: {scanner: "TimetzToByteSlice"},
		},
		{oid: oid.TimetzArr}: {
			analysis.LiteralTimeSlice: {valuer: "TimetzArrayFromTimeSlice", scanner: "TimetzArrayToTimeSlice"},
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.TSQuery}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.TSQueryArr}: {
			analysis.LiteralStringSlice:    {valuer: "TSQueryArrayFromStringSlice", scanner: "TSQueryArrayToStringSlice"},
			analysis.LiteralByteSliceSlice: {valuer: "TSQueryArrayFromByteSliceSlice", scanner: "TSQueryArrayToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.TsRange}: {
			analysis.LiteralTimeArray2: {valuer: "TsRangeFromTimeArray2", scanner: "TsRangeToTimeArray2"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.TsRangeArr}: {
			analysis.LiteralTimeArray2Slice: {valuer: "TsRangeArrayFromTimeArray2Slice", scanner: "TsRangeArrayToTimeArray2Slice"},
			analysis.LiteralString:          {},
			analysis.LiteralByteSlice:       {},
		},
		{oid: oid.TsTzRange}: {
			analysis.LiteralTimeArray2: {valuer: "TstzRangeFromTimeArray2", scanner: "TstzRangeToTimeArray2"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.TsTzRangeArr}: {
			analysis.LiteralTimeArray2Slice: {valuer: "TstzRangeArrayFromTimeArray2Slice", scanner: "TstzRangeArrayToTimeArray2Slice"},
			analysis.LiteralString:          {},
			analysis.LiteralByteSlice:       {},
		},
		{oid: oid.TSVector}: {
			analysis.LiteralStringSlice:    {valuer: "TSVectorFromStringSlice", scanner: "TSVectorToStringSlice"},
			analysis.LiteralByteSliceSlice: {valuer: "TSVectorFromByteSliceSlice", scanner: "TSVectorToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.TSVectorArr}: {
			analysis.LiteralStringSliceSlice:    {valuer: "TSVectorArrayFromStringSliceSlice", scanner: "TSVectorArrayToStringSliceSlice"},
			analysis.LiteralByteSliceSliceSlice: {valuer: "TSVectorArrayFromByteSliceSliceSlice", scanner: "TSVectorArrayToByteSliceSliceSlice"},
			analysis.LiteralString:              {},
			analysis.LiteralByteSlice:           {},
		},
		{oid: oid.UUID}: {
			analysis.LiteralByteArray16: {valuer: "UUIDFromByteArray16", scanner: "UUIDToByteArray16"},
			analysis.LiteralString:      {},
			analysis.LiteralByteSlice:   {},
		},
		{oid: oid.UUIDArr}: {
			analysis.LiteralByteArray16Slice: {valuer: "UUIDArrayFromByteArray16Slice", scanner: "UUIDArrayToByteArray16Slice"},
			analysis.LiteralStringSlice:      {valuer: "UUIDArrayFromStringSlice", scanner: "UUIDArrayToStringSlice"},
			analysis.LiteralByteSliceSlice:   {valuer: "UUIDArrayFromByteSliceSlice", scanner: "UUIDArrayToByteSliceSlice"},
			analysis.LiteralString:           {},
			analysis.LiteralByteSlice:        {},
		},
		{oid: oid.VarBit}: {
			analysis.LiteralInt64:      {valuer: "VarBitFromInt64", scanner: "VarBitToInt64"},
			analysis.LiteralBoolSlice:  {valuer: "VarBitFromBoolSlice", scanner: "VarBitToBoolSlice"},
			analysis.LiteralUint8Slice: {valuer: "VarBitFromUint8Slice", scanner: "VarBitToUint8Slice"},
			analysis.LiteralString:     {},
			analysis.LiteralByteSlice:  {},
		},
		{oid: oid.VarBitArr}: {
			analysis.LiteralInt64Slice:      {valuer: "VarBitArrayFromInt64Slice", scanner: "VarBitArrayToInt64Slice"},
			analysis.LiteralBoolSliceSlice:  {valuer: "VarBitArrayFromBoolSliceSlice", scanner: "VarBitArrayToBoolSliceSlice"},
			analysis.LiteralUint8SliceSlice: {valuer: "VarBitArrayFromUint8SliceSlice", scanner: "VarBitArrayToUint8SliceSlice"},
			analysis.LiteralStringSlice:     {valuer: "VarBitArrayFromStringSlice", scanner: "VarBitArrayToStringSlice"},
			analysis.LiteralString:          {},
			analysis.LiteralByteSlice:       {},
		},
		{oid: oid.VarChar}: {
			analysis.LiteralString:    {},
			analysis.LiteralByteSlice: {},
		},
		{oid: oid.VarCharArr}: {
			analysis.LiteralStringSlice:    {valuer: "VarCharArrayFromStringSlice", scanner: "VarCharArrayToStringSlice"},
			analysis.LiteralByteSliceSlice: {valuer: "VarCharArrayFromByteSliceSlice", scanner: "VarCharArrayToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.XML}: {
			analysis.LiteralEmptyInterface: {valuer: "XML", scanner: "XML"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
		{oid: oid.XMLArr}: {
			analysis.LiteralByteSliceSlice: {valuer: "XMLArrayFromByteSliceSlice", scanner: "XMLArrayToByteSliceSlice"},
			analysis.LiteralString:         {},
			analysis.LiteralByteSlice:      {},
		},
	}
}
