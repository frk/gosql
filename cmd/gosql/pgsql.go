package main

type pgsqlTypeKey struct {
	oid     pgoid
	typmod1 bool
	noscale bool
}

type pgsqlTypeEntry struct {
	valuer  string
	scanner string
}

var pgsqlTypeTable = map[pgsqlTypeKey]map[goTypeId]pgsqlTypeEntry{
	{oid: pgtyp_bit, typmod1: true}: {
		goTypeBool:      {valuer: "BitFromBool"},
		goTypeUint:      {},
		goTypeUint8:     {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_bitarr, typmod1: true}: {
		goTypeBoolSlice:  {valuer: "BitArrayFromBoolSlice", scanner: "BitArrayToBoolSlice"},
		goTypeUintSlice:  {valuer: "BitArrayFromUintSlice", scanner: "BitArrayToUintSlice"},
		goTypeUint8Slice: {valuer: "BitArrayFromUint8Slice", scanner: "BitArrayToUint8Slice"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_bool}: {
		goTypeBool:      {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_boolarr}: {
		goTypeBoolSlice: {valuer: "BoolArrayFromBoolSlice", scanner: "BoolArrayToBoolSlice"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_box}: {
		goTypeFloat64Array2Array2: {valuer: "BoxFromFloat64Array2Array2", scanner: "BoxToFloat64Array2Array2"},
		goTypeString:              {},
		goTypeByteSlice:           {},
	},
	{oid: pgtyp_boxarr}: {
		goTypeFloat64Array2Array2Slice: {valuer: "BoxArrayFromFloat64Array2Array2Slice", scanner: "BoxArrayToFloat64Array2Array2Slice"},
		goTypeString:                   {},
		goTypeByteSlice:                {},
	},
	{oid: pgtyp_bpchar, typmod1: true}: {
		goTypeByte:      {valuer: "BPCharFromByte", scanner: "BPCharToByte"},
		goTypeRune:      {valuer: "BPCharFromRune", scanner: "BPCharToRune"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_bpchararr, typmod1: true}: {
		goTypeRuneSlice:   {valuer: "BPCharArrayFromRuneSlice", scanner: "BPCharArrayToRuneSlice"},
		goTypeStringSlice: {valuer: "BPCharArrayFromStringSlice", scanner: "BPCharArrayToStringSlice"},
		goTypeString:      {valuer: "BPCharArrayFromString", scanner: "BPCharArrayToString"},
		goTypeByteSlice:   {valuer: "BPCharArrayFromByteSlice", scanner: "BPCharArrayToByteSlice"},
	},
	{oid: pgtyp_bytea}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_byteaarr}: {
		goTypeStringSlice:    {valuer: "ByteaArrayFromStringSlice", scanner: "ByteaArrayToStringSlice"},
		goTypeByteSliceSlice: {valuer: "ByteaArrayFromByteSliceSlice", scanner: "ByteaArrayToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_char, typmod1: true}: {
		goTypeByte:      {valuer: "CharFromByte", scanner: "CharToByte"},
		goTypeRune:      {valuer: "CharFromRune", scanner: "CharToRune"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_chararr, typmod1: true}: {
		goTypeRuneSlice:   {valuer: "CharArrayFromRuneSlice", scanner: "CharArrayToRuneSlice"},
		goTypeStringSlice: {valuer: "CharArrayFromStringSlice", scanner: "CharArrayToStringSlice"},
		goTypeString:      {valuer: "CharArrayFromString", scanner: "CharArrayToString"},
		goTypeByteSlice:   {valuer: "CharArrayFromByteSlice", scanner: "CharArrayToByteSlice"},
	},
	{oid: pgtyp_cidr}: {
		goTypeIPNet:     {valuer: "CIDRFromIPNet", scanner: "CIDRToIPNet"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_cidrarr}: {
		goTypeIPNetSlice: {valuer: "CIDRArrayFromIPNetSlice", scanner: "CIDRArrayToIPNetSlice"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_circle}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_circlearr}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_date}: {
		goTypeTime:      {scanner: "DateToTime"},
		goTypeString:    {scanner: "DateToString"},
		goTypeByteSlice: {scanner: "DateToByteSlice"},
	},
	{oid: pgtyp_datearr}: {
		goTypeTimeSlice: {valuer: "DateArrayFromTimeSlice", scanner: "DateArrayToTimeSlice"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_daterange}: {
		goTypeTimeArray2: {valuer: "DateRangeFromTimeArray2", scanner: "DateRangeToTimeArray2"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_daterangearr}: {
		goTypeTimeArray2Slice: {valuer: "DateRangeArrayFromTimeArray2Slice", scanner: "DateRangeArrayToTimeArray2Slice"},
		goTypeString:          {},
		goTypeByteSlice:       {},
	},
	{oid: pgtyp_float4}: {
		goTypeFloat32:   {},
		goTypeFloat64:   {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_float4arr}: {
		goTypeFloat32Slice: {valuer: "Float4ArrayFromFloat32Slice", scanner: "Float4ArrayToFloat32Slice"},
		goTypeFloat64Slice: {valuer: "Float4ArrayFromFloat64Slice", scanner: "Float4ArrayToFloat64Slice"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_float8}: {
		goTypeFloat32:   {},
		goTypeFloat64:   {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_float8arr}: {
		goTypeFloat32Slice: {valuer: "Float8ArrayFromFloat32Slice", scanner: "Float8ArrayToFloat32Slice"},
		goTypeFloat64Slice: {valuer: "Float8ArrayFromFloat64Slice", scanner: "Float8ArrayToFloat64Slice"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_hstore}: {
		goTypeStringMap:     {valuer: "HStoreFromStringMap", scanner: "HStoreToStringMap"},
		goTypeStringPtrMap:  {valuer: "HStoreFromStringPtrMap", scanner: "HStoreToStringPtrMap"},
		goTypeNullStringMap: {valuer: "HStoreFromNullStringMap", scanner: "HStoreToNullStringMap"},
		goTypeString:        {},
		goTypeByteSlice:     {},
	},
	{oid: pgtyp_hstorearr}: {
		goTypeStringMapSlice:     {valuer: "HStoreArrayFromStringMapSlice", scanner: "HStoreArrayToStringMapSlice"},
		goTypeStringPtrMapSlice:  {valuer: "HStoreArrayFromStringPtrMapSlice", scanner: "HStoreArrayToStringPtrMapSlice"},
		goTypeNullStringMapSlice: {valuer: "HStoreArrayFromNullStringMapSlice", scanner: "HStoreArrayToNullStringMapSlice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_inet}: {
		goTypeIP:        {valuer: "InetFromIP", scanner: "InetToIP"},
		goTypeIPNet:     {valuer: "InetFromIPNet", scanner: "InetToIPNet"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_inetarr}: {
		goTypeIPSlice:    {valuer: "InetArrayFromIPSlice", scanner: "InetArrayFromIPSlice"},
		goTypeIPNetSlice: {valuer: "InetArrayFromIPNetSlice", scanner: "InetArrayFromIPNetSlice"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_int2}: {
		goTypeInt:       {},
		goTypeInt8:      {},
		goTypeInt16:     {},
		goTypeInt32:     {},
		goTypeInt64:     {},
		goTypeUint:      {},
		goTypeUint8:     {},
		goTypeUint16:    {},
		goTypeUint32:    {},
		goTypeUint64:    {},
		goTypeFloat32:   {},
		goTypeFloat64:   {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_int2arr}: {
		goTypeIntSlice:     {valuer: "Int2ArrayFromIntSlice", scanner: "Int2ArrayToIntSlice"},
		goTypeInt8Slice:    {valuer: "Int2ArrayFromInt8Slice", scanner: "Int2ArrayToInt8Slice"},
		goTypeInt16Slice:   {valuer: "Int2ArrayFromInt16Slice", scanner: "Int2ArrayToInt16Slice"},
		goTypeInt32Slice:   {valuer: "Int2ArrayFromInt32Slice", scanner: "Int2ArrayToInt32Slice"},
		goTypeInt64Slice:   {valuer: "Int2ArrayFromInt64Slice", scanner: "Int2ArrayToInt64Slice"},
		goTypeUintSlice:    {valuer: "Int2ArrayFromUintSlice", scanner: "Int2ArrayToUintSlice"},
		goTypeUint8Slice:   {valuer: "Int2ArrayFromUint8Slice", scanner: "Int2ArrayToUint8Slice"},
		goTypeUint16Slice:  {valuer: "Int2ArrayFromUint16Slice", scanner: "Int2ArrayToUint16Slice"},
		goTypeUint32Slice:  {valuer: "Int2ArrayFromUint32Slice", scanner: "Int2ArrayToUint32Slice"},
		goTypeUint64Slice:  {valuer: "Int2ArrayFromUint64Slice", scanner: "Int2ArrayToUint64Slice"},
		goTypeFloat32Slice: {valuer: "Int2ArrayFromFloat32Slice", scanner: "Int2ArrayToFloat32Slice"},
		goTypeFloat64Slice: {valuer: "Int2ArrayFromFloat64Slice", scanner: "Int2ArrayToFloat64Slice"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_int2vector}: {
		goTypeIntSlice:     {valuer: "Int2VectorFromIntSlice", scanner: "Int2VectorToIntSlice"},
		goTypeInt8Slice:    {valuer: "Int2VectorFromInt8Slice", scanner: "Int2VectorToInt8Slice"},
		goTypeInt16Slice:   {valuer: "Int2VectorFromInt16Slice", scanner: "Int2VectorToInt16Slice"},
		goTypeInt32Slice:   {valuer: "Int2VectorFromInt32Slice", scanner: "Int2VectorToInt32Slice"},
		goTypeInt64Slice:   {valuer: "Int2VectorFromInt64Slice", scanner: "Int2VectorToInt64Slice"},
		goTypeUintSlice:    {valuer: "Int2VectorFromUintSlice", scanner: "Int2VectorToUintSlice"},
		goTypeUint8Slice:   {valuer: "Int2VectorFromUint8Slice", scanner: "Int2VectorToUint8Slice"},
		goTypeUint16Slice:  {valuer: "Int2VectorFromUint16Slice", scanner: "Int2VectorToUint16Slice"},
		goTypeUint32Slice:  {valuer: "Int2VectorFromUint32Slice", scanner: "Int2VectorToUint32Slice"},
		goTypeUint64Slice:  {valuer: "Int2VectorFromUint64Slice", scanner: "Int2VectorToUint64Slice"},
		goTypeFloat32Slice: {valuer: "Int2VectorFromFloat32Slice", scanner: "Int2VectorToFloat32Slice"},
		goTypeFloat64Slice: {valuer: "Int2VectorFromFloat64Slice", scanner: "Int2VectorToFloat64Slice"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_int2vectorarr}: {
		goTypeIntSliceSlice:     {valuer: "Int2VectorArrayFromIntSliceSlice", scanner: "Int2VectorArrayToIntSliceSlice"},
		goTypeInt8SliceSlice:    {valuer: "Int2VectorArrayFromInt8SliceSlice", scanner: "Int2VectorArrayToInt8SliceSlice"},
		goTypeInt16SliceSlice:   {valuer: "Int2VectorArrayFromInt16SliceSlice", scanner: "Int2VectorArrayToInt16SliceSlice"},
		goTypeInt32SliceSlice:   {valuer: "Int2VectorArrayFromInt32SliceSlice", scanner: "Int2VectorArrayToInt32SliceSlice"},
		goTypeInt64SliceSlice:   {valuer: "Int2VectorArrayFromInt64SliceSlice", scanner: "Int2VectorArrayToInt64SliceSlice"},
		goTypeUintSliceSlice:    {valuer: "Int2VectorArrayFromUintSliceSlice", scanner: "Int2VectorArrayToUintSliceSlice"},
		goTypeUint8SliceSlice:   {valuer: "Int2VectorArrayFromUint8SliceSlice", scanner: "Int2VectorArrayToUint8SliceSlice"},
		goTypeUint16SliceSlice:  {valuer: "Int2VectorArrayFromUint16SliceSlice", scanner: "Int2VectorArrayToUint16SliceSlice"},
		goTypeUint32SliceSlice:  {valuer: "Int2VectorArrayFromUint32SliceSlice", scanner: "Int2VectorArrayToUint32SliceSlice"},
		goTypeUint64SliceSlice:  {valuer: "Int2VectorArrayFromUint64SliceSlice", scanner: "Int2VectorArrayToUint64SliceSlice"},
		goTypeFloat32SliceSlice: {valuer: "Int2VectorArrayFromFloat32SliceSlice", scanner: "Int2VectorArrayToFloat32SliceSlice"},
		goTypeFloat64SliceSlice: {valuer: "Int2VectorArrayFromFloat64SliceSlice", scanner: "Int2VectorArrayToFloat64SliceSlice"},
		goTypeString:            {},
		goTypeByteSlice:         {},
	},
	{oid: pgtyp_int4}: {
		goTypeInt:       {},
		goTypeInt8:      {},
		goTypeInt16:     {},
		goTypeInt32:     {},
		goTypeInt64:     {},
		goTypeUint:      {},
		goTypeUint8:     {},
		goTypeUint16:    {},
		goTypeUint32:    {},
		goTypeUint64:    {},
		goTypeFloat32:   {},
		goTypeFloat64:   {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_int4arr}: {
		goTypeIntSlice:     {valuer: "Int4ArrayFromIntSlice", scanner: "Int4ArrayToIntSlice"},
		goTypeInt8Slice:    {valuer: "Int4ArrayFromInt8Slice", scanner: "Int4ArrayToInt8Slice"},
		goTypeInt16Slice:   {valuer: "Int4ArrayFromInt16Slice", scanner: "Int4ArrayToInt16Slice"},
		goTypeInt32Slice:   {valuer: "Int4ArrayFromInt32Slice", scanner: "Int4ArrayToInt32Slice"},
		goTypeInt64Slice:   {valuer: "Int4ArrayFromInt64Slice", scanner: "Int4ArrayToInt64Slice"},
		goTypeUintSlice:    {valuer: "Int4ArrayFromUintSlice", scanner: "Int4ArrayToUintSlice"},
		goTypeUint8Slice:   {valuer: "Int4ArrayFromUint8Slice", scanner: "Int4ArrayToUint8Slice"},
		goTypeUint16Slice:  {valuer: "Int4ArrayFromUint16Slice", scanner: "Int4ArrayToUint16Slice"},
		goTypeUint32Slice:  {valuer: "Int4ArrayFromUint32Slice", scanner: "Int4ArrayToUint32Slice"},
		goTypeUint64Slice:  {valuer: "Int4ArrayFromUint64Slice", scanner: "Int4ArrayToUint64Slice"},
		goTypeFloat32Slice: {valuer: "Int4ArrayFromFloat32Slice", scanner: "Int4ArrayToFloat32Slice"},
		goTypeFloat64Slice: {valuer: "Int4ArrayFromFloat64Slice", scanner: "Int4ArrayToFloat64Slice"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_int4range}: {
		goTypeIntArray2:     {valuer: "Int4RangeFromIntArray2", scanner: "Int4RangeToIntArray2"},
		goTypeInt8Array2:    {valuer: "Int4RangeFromInt8Array2", scanner: "Int4RangeToInt8Array2"},
		goTypeInt16Array2:   {valuer: "Int4RangeFromInt16Array2", scanner: "Int4RangeToInt16Array2"},
		goTypeInt32Array2:   {valuer: "Int4RangeFromInt32Array2", scanner: "Int4RangeToInt32Array2"},
		goTypeInt64Array2:   {valuer: "Int4RangeFromInt64Array2", scanner: "Int4RangeToInt64Array2"},
		goTypeUintArray2:    {valuer: "Int4RangeFromUintArray2", scanner: "Int4RangeToUintArray2"},
		goTypeUint8Array2:   {valuer: "Int4RangeFromUint8Array2", scanner: "Int4RangeToUint8Array2"},
		goTypeUint16Array2:  {valuer: "Int4RangeFromUint16Array2", scanner: "Int4RangeToUint16Array2"},
		goTypeUint32Array2:  {valuer: "Int4RangeFromUint32Array2", scanner: "Int4RangeToUint32Array2"},
		goTypeUint64Array2:  {valuer: "Int4RangeFromUint64Array2", scanner: "Int4RangeToUint64Array2"},
		goTypeFloat32Array2: {valuer: "Int4RangeFromFloat32Array2", scanner: "Int4RangeToFloat32Array2"},
		goTypeFloat64Array2: {valuer: "Int4RangeFromFloat64Array2", scanner: "Int4RangeToFloat64Array2"},
		goTypeString:        {},
		goTypeByteSlice:     {},
	},
	{oid: pgtyp_int4rangearr}: {
		goTypeIntArray2Slice:     {valuer: "Int4RangeArrayFromIntArray2Slice", scanner: "Int4RangeArrayToIntArray2Slice"},
		goTypeInt8Array2Slice:    {valuer: "Int4RangeArrayFromInt8Array2Slice", scanner: "Int4RangeArrayToInt8Array2Slice"},
		goTypeInt16Array2Slice:   {valuer: "Int4RangeArrayFromInt16Array2Slice", scanner: "Int4RangeArrayToInt16Array2Slice"},
		goTypeInt32Array2Slice:   {valuer: "Int4RangeArrayFromInt32Array2Slice", scanner: "Int4RangeArrayToInt32Array2Slice"},
		goTypeInt64Array2Slice:   {valuer: "Int4RangeArrayFromInt64Array2Slice", scanner: "Int4RangeArrayToInt64Array2Slice"},
		goTypeUintArray2Slice:    {valuer: "Int4RangeArrayFromUintArray2Slice", scanner: "Int4RangeArrayToUintArray2Slice"},
		goTypeUint8Array2Slice:   {valuer: "Int4RangeArrayFromUint8Array2Slice", scanner: "Int4RangeArrayToUint8Array2Slice"},
		goTypeUint16Array2Slice:  {valuer: "Int4RangeArrayFromUint16Array2Slice", scanner: "Int4RangeArrayToUint16Array2Slice"},
		goTypeUint32Array2Slice:  {valuer: "Int4RangeArrayFromUint32Array2Slice", scanner: "Int4RangeArrayToUint32Array2Slice"},
		goTypeUint64Array2Slice:  {valuer: "Int4RangeArrayFromUint64Array2Slice", scanner: "Int4RangeArrayToUint64Array2Slice"},
		goTypeFloat32Array2Slice: {valuer: "Int4RangeArrayFromFloat32Array2Slice", scanner: "Int4RangeArrayToFloat32Array2Slice"},
		goTypeFloat64Array2Slice: {valuer: "Int4RangeArrayFromFloat64Array2Slice", scanner: "Int4RangeArrayToFloat64Array2Slice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_int8}: {
		goTypeInt:       {},
		goTypeInt8:      {},
		goTypeInt16:     {},
		goTypeInt32:     {},
		goTypeInt64:     {},
		goTypeUint:      {},
		goTypeUint8:     {},
		goTypeUint16:    {},
		goTypeUint32:    {},
		goTypeUint64:    {},
		goTypeFloat32:   {},
		goTypeFloat64:   {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_int8arr}: {
		goTypeIntSlice:     {valuer: "Int8ArrayFromIntSlice", scanner: "Int8ArrayToIntSlice"},
		goTypeInt8Slice:    {valuer: "Int8ArrayFromInt8Slice", scanner: "Int8ArrayToInt8Slice"},
		goTypeInt16Slice:   {valuer: "Int8ArrayFromInt16Slice", scanner: "Int8ArrayToInt16Slice"},
		goTypeInt32Slice:   {valuer: "Int8ArrayFromInt32Slice", scanner: "Int8ArrayToInt32Slice"},
		goTypeInt64Slice:   {valuer: "Int8ArrayFromInt64Slice", scanner: "Int8ArrayToInt64Slice"},
		goTypeUintSlice:    {valuer: "Int8ArrayFromUintSlice", scanner: "Int8ArrayToUintSlice"},
		goTypeUint8Slice:   {valuer: "Int8ArrayFromUint8Slice", scanner: "Int8ArrayToUint8Slice"},
		goTypeUint16Slice:  {valuer: "Int8ArrayFromUint16Slice", scanner: "Int8ArrayToUint16Slice"},
		goTypeUint32Slice:  {valuer: "Int8ArrayFromUint32Slice", scanner: "Int8ArrayToUint32Slice"},
		goTypeUint64Slice:  {valuer: "Int8ArrayFromUint64Slice", scanner: "Int8ArrayToUint64Slice"},
		goTypeFloat32Slice: {valuer: "Int8ArrayFromFloat32Slice", scanner: "Int8ArrayToFloat32Slice"},
		goTypeFloat64Slice: {valuer: "Int8ArrayFromFloat64Slice", scanner: "Int8ArrayToFloat64Slice"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_int8range}: {
		goTypeIntArray2:     {valuer: "Int8RangeFromIntArray2", scanner: "Int8RangeToIntArray2"},
		goTypeInt8Array2:    {valuer: "Int8RangeFromInt8Array2", scanner: "Int8RangeToInt8Array2"},
		goTypeInt16Array2:   {valuer: "Int8RangeFromInt16Array2", scanner: "Int8RangeToInt16Array2"},
		goTypeInt32Array2:   {valuer: "Int8RangeFromInt32Array2", scanner: "Int8RangeToInt32Array2"},
		goTypeInt64Array2:   {valuer: "Int8RangeFromInt64Array2", scanner: "Int8RangeToInt64Array2"},
		goTypeUintArray2:    {valuer: "Int8RangeFromUintArray2", scanner: "Int8RangeToUintArray2"},
		goTypeUint8Array2:   {valuer: "Int8RangeFromUint8Array2", scanner: "Int8RangeToUint8Array2"},
		goTypeUint16Array2:  {valuer: "Int8RangeFromUint16Array2", scanner: "Int8RangeToUint16Array2"},
		goTypeUint32Array2:  {valuer: "Int8RangeFromUint32Array2", scanner: "Int8RangeToUint32Array2"},
		goTypeUint64Array2:  {valuer: "Int8RangeFromUint64Array2", scanner: "Int8RangeToUint64Array2"},
		goTypeFloat32Array2: {valuer: "Int8RangeFromFloat32Array2", scanner: "Int8RangeToFloat32Array2"},
		goTypeFloat64Array2: {valuer: "Int8RangeFromFloat64Array2", scanner: "Int8RangeToFloat64Array2"},
		goTypeString:        {},
		goTypeByteSlice:     {},
	},
	{oid: pgtyp_int8rangearr}: {
		goTypeIntArray2Slice:     {valuer: "Int8RangeArrayFromIntArray2Slice", scanner: "Int8RangeArrayToIntArray2Slice"},
		goTypeInt8Array2Slice:    {valuer: "Int8RangeArrayFromInt8Array2Slice", scanner: "Int8RangeArrayToInt8Array2Slice"},
		goTypeInt16Array2Slice:   {valuer: "Int8RangeArrayFromInt16Array2Slice", scanner: "Int8RangeArrayToInt16Array2Slice"},
		goTypeInt32Array2Slice:   {valuer: "Int8RangeArrayFromInt32Array2Slice", scanner: "Int8RangeArrayToInt32Array2Slice"},
		goTypeInt64Array2Slice:   {valuer: "Int8RangeArrayFromInt64Array2Slice", scanner: "Int8RangeArrayToInt64Array2Slice"},
		goTypeUintArray2Slice:    {valuer: "Int8RangeArrayFromUintArray2Slice", scanner: "Int8RangeArrayToUintArray2Slice"},
		goTypeUint8Array2Slice:   {valuer: "Int8RangeArrayFromUint8Array2Slice", scanner: "Int8RangeArrayToUint8Array2Slice"},
		goTypeUint16Array2Slice:  {valuer: "Int8RangeArrayFromUint16Array2Slice", scanner: "Int8RangeArrayToUint16Array2Slice"},
		goTypeUint32Array2Slice:  {valuer: "Int8RangeArrayFromUint32Array2Slice", scanner: "Int8RangeArrayToUint32Array2Slice"},
		goTypeUint64Array2Slice:  {valuer: "Int8RangeArrayFromUint64Array2Slice", scanner: "Int8RangeArrayToUint64Array2Slice"},
		goTypeFloat32Array2Slice: {valuer: "Int8RangeArrayFromFloat32Array2Slice", scanner: "Int8RangeArrayToFloat32Array2Slice"},
		goTypeFloat64Array2Slice: {valuer: "Int8RangeArrayFromFloat64Array2Slice", scanner: "Int8RangeArrayToFloat64Array2Slice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_interval}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_intervalarr}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_json}: {
		goTypeEmptyInterface: {valuer: "JSON", scanner: "JSON"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_jsonarr}: {
		goTypeByteSliceSlice: {valuer: "JSONArrayFromByteSliceSlice", scanner: "JSONArrayToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_jsonb}: {
		goTypeEmptyInterface: {valuer: "JSON", scanner: "JSON"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_jsonbarr}: {
		goTypeByteSliceSlice: {valuer: "JSONArrayFromByteSliceSlice", scanner: "JSONArrayToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_line}: {
		goTypeFloat64Array3: {valuer: "LineFromFloat64Array3", scanner: "LineToFloat64Array3"},
		goTypeString:        {},
		goTypeByteSlice:     {},
	},
	{oid: pgtyp_linearr}: {
		goTypeFloat64Array3Slice: {valuer: "LineArrayFromFloat64Array3Slice", scanner: "LineArrayToFloat64Array3Slice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_lseg}: {
		goTypeFloat64Array2Array2: {valuer: "LsegFromFloat64Array2Array2", scanner: "LsegToFloat64Array2Array2"},
		goTypeString:              {},
		goTypeByteSlice:           {},
	},
	{oid: pgtyp_lsegarr}: {
		goTypeFloat64Array2Array2Slice: {valuer: "LsegArrayFromFloat64Array2Array2Slice", scanner: "LsegArrayToFloat64Array2Array2Slice"},
		goTypeString:                   {},
		goTypeByteSlice:                {},
	},
	{oid: pgtyp_macaddr}: {
		goTypeHardwareAddr: {valuer: "MACAddrFromHardwareAddr", scanner: "MACAddrToHardwareAddr"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_macaddrarr}: {
		goTypeHardwareAddrSlice: {valuer: "MACAddrArrayFromHardwareAddrSlice", scanner: "MACAddrArrayToHardwareAddrSlice"},
		goTypeString:            {},
		goTypeByteSlice:         {},
	},
	{oid: pgtyp_macaddr8}: {
		goTypeHardwareAddr: {valuer: "MACAddr8FromHardwareAddr", scanner: "MACAddr8ToHardwareAddr"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_macaddr8arr}: {
		goTypeHardwareAddrSlice: {valuer: "MACAddr8ArrayFromHardwareAddrSlice", scanner: "MACAddr8ArrayToHardwareAddrSlice"},
		goTypeString:            {},
		goTypeByteSlice:         {},
	},
	{oid: pgtyp_money}: {
		goTypeInt64:     {valuer: "MoneyFromInt64", scanner: "MoneyToInt64"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_moneyarr}: {
		goTypeInt64Slice: {valuer: "MoneyArrayFromInt64Slice", scanner: "MoneyArrayToInt64Slice"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_numeric}: {
		goTypeInt:       {},
		goTypeInt8:      {},
		goTypeInt16:     {},
		goTypeInt32:     {},
		goTypeInt64:     {},
		goTypeUint:      {},
		goTypeUint8:     {},
		goTypeUint16:    {},
		goTypeUint32:    {},
		goTypeUint64:    {},
		goTypeFloat32:   {},
		goTypeFloat64:   {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_numericarr}: {
		goTypeIntSlice:     {valuer: "NumericArrayFromIntSlice", scanner: "NumericArrayToIntSlice"},
		goTypeInt8Slice:    {valuer: "NumericArrayFromInt8Slice", scanner: "NumericArrayToInt8Slice"},
		goTypeInt16Slice:   {valuer: "NumericArrayFromInt16Slice", scanner: "NumericArrayToInt16Slice"},
		goTypeInt32Slice:   {valuer: "NumericArrayFromInt32Slice", scanner: "NumericArrayToInt32Slice"},
		goTypeInt64Slice:   {valuer: "NumericArrayFromInt64Slice", scanner: "NumericArrayToInt64Slice"},
		goTypeUintSlice:    {valuer: "NumericArrayFromUintSlice", scanner: "NumericArrayToUintSlice"},
		goTypeUint8Slice:   {valuer: "NumericArrayFromUint8Slice", scanner: "NumericArrayToUint8Slice"},
		goTypeUint16Slice:  {valuer: "NumericArrayFromUint16Slice", scanner: "NumericArrayToUint16Slice"},
		goTypeUint32Slice:  {valuer: "NumericArrayFromUint32Slice", scanner: "NumericArrayToUint32Slice"},
		goTypeUint64Slice:  {valuer: "NumericArrayFromUint64Slice", scanner: "NumericArrayToUint64Slice"},
		goTypeFloat32Slice: {valuer: "NumericArrayFromFloat32Slice", scanner: "NumericArrayToFloat32Slice"},
		goTypeFloat64Slice: {valuer: "NumericArrayFromFloat64Slice", scanner: "NumericArrayToFloat64Slice"},
		goTypeString:       {},
		goTypeByteSlice:    {},
	},
	{oid: pgtyp_numrange}: {
		goTypeIntArray2:     {valuer: "NumRangeFromIntArray2", scanner: "NumRangeToIntArray2"},
		goTypeInt8Array2:    {valuer: "NumRangeFromInt8Array2", scanner: "NumRangeToInt8Array2"},
		goTypeInt16Array2:   {valuer: "NumRangeFromInt16Array2", scanner: "NumRangeToInt16Array2"},
		goTypeInt32Array2:   {valuer: "NumRangeFromInt32Array2", scanner: "NumRangeToInt32Array2"},
		goTypeInt64Array2:   {valuer: "NumRangeFromInt64Array2", scanner: "NumRangeToInt64Array2"},
		goTypeUintArray2:    {valuer: "NumRangeFromUintArray2", scanner: "NumRangeToUintArray2"},
		goTypeUint8Array2:   {valuer: "NumRangeFromUint8Array2", scanner: "NumRangeToUint8Array2"},
		goTypeUint16Array2:  {valuer: "NumRangeFromUint16Array2", scanner: "NumRangeToUint16Array2"},
		goTypeUint32Array2:  {valuer: "NumRangeFromUint32Array2", scanner: "NumRangeToUint32Array2"},
		goTypeUint64Array2:  {valuer: "NumRangeFromUint64Array2", scanner: "NumRangeToUint64Array2"},
		goTypeFloat32Array2: {valuer: "NumRangeFromFloat32Array2", scanner: "NumRangeToFloat32Array2"},
		goTypeFloat64Array2: {valuer: "NumRangeFromFloat64Array2", scanner: "NumRangeToFloat64Array2"},
		goTypeString:        {},
		goTypeByteSlice:     {},
	},
	{oid: pgtyp_numrangearr}: {
		goTypeIntArray2Slice:     {valuer: "NumRangeArrayFromIntArray2Slice", scanner: "NumRangeArrayToIntArray2Slice"},
		goTypeInt8Array2Slice:    {valuer: "NumRangeArrayFromInt8Array2Slice", scanner: "NumRangeArrayToInt8Array2Slice"},
		goTypeInt16Array2Slice:   {valuer: "NumRangeArrayFromInt16Array2Slice", scanner: "NumRangeArrayToInt16Array2Slice"},
		goTypeInt32Array2Slice:   {valuer: "NumRangeArrayFromInt32Array2Slice", scanner: "NumRangeArrayToInt32Array2Slice"},
		goTypeInt64Array2Slice:   {valuer: "NumRangeArrayFromInt64Array2Slice", scanner: "NumRangeArrayToInt64Array2Slice"},
		goTypeUintArray2Slice:    {valuer: "NumRangeArrayFromUintArray2Slice", scanner: "NumRangeArrayToUintArray2Slice"},
		goTypeUint8Array2Slice:   {valuer: "NumRangeArrayFromUint8Array2Slice", scanner: "NumRangeArrayToUint8Array2Slice"},
		goTypeUint16Array2Slice:  {valuer: "NumRangeArrayFromUint16Array2Slice", scanner: "NumRangeArrayToUint16Array2Slice"},
		goTypeUint32Array2Slice:  {valuer: "NumRangeArrayFromUint32Array2Slice", scanner: "NumRangeArrayToUint32Array2Slice"},
		goTypeUint64Array2Slice:  {valuer: "NumRangeArrayFromUint64Array2Slice", scanner: "NumRangeArrayToUint64Array2Slice"},
		goTypeFloat32Array2Slice: {valuer: "NumRangeArrayFromFloat32Array2Slice", scanner: "NumRangeArrayToFloat32Array2Slice"},
		goTypeFloat64Array2Slice: {valuer: "NumRangeArrayFromFloat64Array2Slice", scanner: "NumRangeArrayToFloat64Array2Slice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_path}: {
		goTypeFloat64Array2Slice: {valuer: "PathFromFloat64Array2Slice", scanner: "PathToFloat64Array2Slice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_patharr}: {
		goTypeFloat64Array2SliceSlice: {valuer: "PathArrayFromFloat64Array2SliceSlice", scanner: "PathArrayToFloat64Array2SliceSlice"},
		goTypeString:                  {},
		goTypeByteSlice:               {},
	},
	{oid: pgtyp_point}: {
		goTypeFloat64Array2: {valuer: "PointFromFloat64Array2", scanner: "PointToFloat64Array2"},
		goTypeString:        {},
		goTypeByteSlice:     {},
	},
	{oid: pgtyp_pointarr}: {
		goTypeFloat64Array2Slice: {valuer: "PointArrayFromFloat64Array2Slice", scanner: "PointArrayToFloat64Array2Slice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_polygon}: {
		goTypeFloat64Array2Slice: {valuer: "PolygonFromFloat64Array2Slice", scanner: "PolygonToFloat64Array2Slice"},
		goTypeString:             {},
		goTypeByteSlice:          {},
	},
	{oid: pgtyp_polygonarr}: {
		goTypeFloat64Array2SliceSlice: {valuer: "PolygonArrayFromFloat64Array2SliceSlice", scanner: "PolygonArrayToFloat64Array2SliceSlice"},
		goTypeString:                  {},
		goTypeByteSlice:               {},
	},
	{oid: pgtyp_text}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_textarr}: {
		goTypeStringSlice:    {valuer: "TextArrayFromStringSlice", scanner: "TextArrayToStringSlice"},
		goTypeByteSliceSlice: {valuer: "TextArrayFromByteSliceSlice", scanner: "TextArrayToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_time}: {
		goTypeTime:      {},
		goTypeString:    {scanner: "TimeToString"},
		goTypeByteSlice: {scanner: "TimeToByteSlice"},
	},
	{oid: pgtyp_timearr}: {
		goTypeTimeSlice: {valuer: "TimeArrayFromTimeSlice", scanner: "TimeArrayToTimeSlice"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_timestamp}: {
		goTypeTime:      {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_timestamparr}: {
		goTypeTimeSlice: {valuer: "TimestampArrayFromTimeSlice", scanner: "TimestampArrayToTimeSlice"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_timestamptz}: {
		goTypeTime:      {},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_timestamptzarr}: {
		goTypeTimeSlice: {valuer: "TimestamptzArrayFromTimeSlice", scanner: "TimestamptzArrayToTimeSlice"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_timetz}: {
		goTypeTime:      {},
		goTypeString:    {scanner: "TimetzToString"},
		goTypeByteSlice: {scanner: "TimetzToByteSlice"},
	},
	{oid: pgtyp_timetzarr}: {
		goTypeTimeSlice: {valuer: "TimetzArrayFromTimeSlice", scanner: "TimetzArrayToTimeSlice"},
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_tsquery}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_tsqueryarr}: {
		goTypeStringSlice:    {valuer: "TSQueryArrayFromStringSlice", scanner: "TSQueryArrayToStringSlice"},
		goTypeByteSliceSlice: {valuer: "TSQueryArrayFromByteSliceSlice", scanner: "TSQueryArrayToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_tsrange}: {
		goTypeTimeArray2: {valuer: "TsRangeFromTimeArray2", scanner: "TsRangeToTimeArray2"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_tsrangearr}: {
		goTypeTimeArray2Slice: {valuer: "TsRangeArrayFromTimeArray2Slice", scanner: "TsRangeArrayToTimeArray2Slice"},
		goTypeString:          {},
		goTypeByteSlice:       {},
	},
	{oid: pgtyp_tstzrange}: {
		goTypeTimeArray2: {valuer: "TstzRangeFromTimeArray2", scanner: "TstzRangeToTimeArray2"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_tstzrangearr}: {
		goTypeTimeArray2Slice: {valuer: "TstzRangeArrayFromTimeArray2Slice", scanner: "TstzRangeArrayToTimeArray2Slice"},
		goTypeString:          {},
		goTypeByteSlice:       {},
	},
	{oid: pgtyp_tsvector}: {
		goTypeStringSlice:    {valuer: "TSVectorFromStringSlice", scanner: "TSVectorToStringSlice"},
		goTypeByteSliceSlice: {valuer: "TSVectorFromByteSliceSlice", scanner: "TSVectorToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_tsvectorarr}: {
		goTypeStringSliceSlice:    {valuer: "TSVectorArrayFromStringSliceSlice", scanner: "TSVectorArrayToStringSliceSlice"},
		goTypeByteSliceSliceSlice: {valuer: "TSVectorArrayFromByteSliceSliceSlice", scanner: "TSVectorArrayToByteSliceSliceSlice"},
		goTypeString:              {},
		goTypeByteSlice:           {},
	},
	{oid: pgtyp_uuid}: {
		goTypeByteArray16: {valuer: "UUIDFromByteArray16", scanner: "UUIDToByteArray16"},
		goTypeString:      {},
		goTypeByteSlice:   {},
	},
	{oid: pgtyp_uuidarr}: {
		goTypeByteArray16Slice: {valuer: "UUIDArrayFromByteArray16Slice", scanner: "UUIDArrayToByteArray16Slice"},
		goTypeStringSlice:      {valuer: "UUIDArrayFromStringSlice", scanner: "UUIDArrayToStringSlice"},
		goTypeByteSliceSlice:   {valuer: "UUIDArrayFromByteSliceSlice", scanner: "UUIDArrayToByteSliceSlice"},
		goTypeString:           {},
		goTypeByteSlice:        {},
	},
	{oid: pgtyp_varbit}: {
		goTypeInt64:      {valuer: "VarBitFromInt64", scanner: "VarBitToInt64"},
		goTypeBoolSlice:  {valuer: "VarBitFromBoolSlice", scanner: "VarBitToBoolSlice"},
		goTypeUint8Slice: {valuer: "VarBitFromUint8Slice", scanner: "VarBitToUint8Slice"},
		goTypeString:     {},
		goTypeByteSlice:  {},
	},
	{oid: pgtyp_varbitarr}: {
		goTypeInt64Slice:      {valuer: "VarBitArrayFromInt64Slice", scanner: "VarBitArrayToInt64Slice"},
		goTypeBoolSliceSlice:  {valuer: "VarBitArrayFromBoolSliceSlice", scanner: "VarBitArrayToBoolSliceSlice"},
		goTypeUint8SliceSlice: {valuer: "VarBitArrayFromUint8SliceSlice", scanner: "VarBitArrayToUint8SliceSlice"},
		goTypeStringSlice:     {valuer: "VarBitArrayFromStringSlice", scanner: "VarBitArrayToStringSlice"},
		goTypeString:          {},
		goTypeByteSlice:       {},
	},
	{oid: pgtyp_varchar}: {
		goTypeString:    {},
		goTypeByteSlice: {},
	},
	{oid: pgtyp_varchararr}: {
		goTypeStringSlice:    {valuer: "VarCharArrayFromStringSlice", scanner: "VarCharArrayToStringSlice"},
		goTypeByteSliceSlice: {valuer: "VarCharArrayFromByteSliceSlice", scanner: "VarCharArrayToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_xml}: {
		goTypeEmptyInterface: {valuer: "XML", scanner: "XML"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
	{oid: pgtyp_xmlarr}: {
		goTypeByteSliceSlice: {valuer: "XMLArrayFromByteSliceSlice", scanner: "XMLArrayToByteSliceSlice"},
		goTypeString:         {},
		goTypeByteSlice:      {},
	},
}
