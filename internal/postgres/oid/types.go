package oid

type OID uint32

const (
	Any            OID = 2276
	Bit            OID = 1560
	BitArr         OID = 1561
	Bool           OID = 16
	BoolArr        OID = 1000
	Box            OID = 603
	BoxArr         OID = 1020
	BPChar         OID = 1042
	BPCharArr      OID = 1014
	Bytea          OID = 17
	ByteaArr       OID = 1001
	Char           OID = 18
	CharArr        OID = 1002
	CIDR           OID = 650
	CIDRArr        OID = 651
	Circle         OID = 718
	CircleArr      OID = 719
	Date           OID = 1082
	DateArr        OID = 1182
	DateRange      OID = 3912
	DateRangeArr   OID = 3913
	Float4         OID = 700
	Float4Arr      OID = 1021
	Float8         OID = 701
	Float8Arr      OID = 1022
	Inet           OID = 869
	InetArr        OID = 1041
	Int2           OID = 21
	Int2Arr        OID = 1005
	Int2Vector     OID = 22
	Int2VectorArr  OID = 1006
	Int4           OID = 23
	Int4Arr        OID = 1007
	Int4Range      OID = 3904
	Int4RangeArr   OID = 3905
	Int8           OID = 20
	Int8Arr        OID = 1016
	Int8Range      OID = 3926
	Int8RangeArr   OID = 3927
	Interval       OID = 1186
	IntervalArr    OID = 1187
	JSON           OID = 114
	JSONArr        OID = 199
	JSONB          OID = 3802
	JSONBArr       OID = 3807
	Line           OID = 628
	LineArr        OID = 629
	LSeg           OID = 601
	LSegArr        OID = 1018
	MACAddr        OID = 829
	MACAddrArr     OID = 1040
	MACAddr8       OID = 774
	MACAddr8Arr    OID = 775
	Money          OID = 790
	MoneyArr       OID = 791
	Numeric        OID = 1700
	NumericArr     OID = 1231
	NumRange       OID = 3906
	NumRangeArr    OID = 3907
	OIDVector      OID = 30
	Path           OID = 602
	PathArr        OID = 1019
	Point          OID = 600
	PointArr       OID = 1017
	Polygon        OID = 604
	PolygonArr     OID = 1027
	Text           OID = 25
	TextArr        OID = 1009
	Time           OID = 1083
	TimeArr        OID = 1183
	Timestamp      OID = 1114
	TimestampArr   OID = 1115
	Timestamptz    OID = 1184
	TimestamptzArr OID = 1185
	Timetz         OID = 1266
	TimetzArr      OID = 1270
	TSQuery        OID = 3615
	TSQueryArr     OID = 3645
	TsRange        OID = 3908
	TsRangeArr     OID = 3909
	TsTzRange      OID = 3910
	TsTzRangeArr   OID = 3911
	TSVector       OID = 3614
	TSVectorArr    OID = 3643
	UUID           OID = 2950
	UUIDArr        OID = 2951
	Unknown        OID = 705
	VarBit         OID = 1562
	VarBitArr      OID = 1563
	VarChar        OID = 1043
	VarCharArr     OID = 1015
	XML            OID = 142
	XMLArr         OID = 143

	// HStore    OID = 9999
	// HStoreArr OID = 9998
)

var TypeToArray = map[OID]OID{
	Bit:         BitArr,
	Bool:        BoolArr,
	Box:         BoxArr,
	BPChar:      BPCharArr,
	Bytea:       ByteaArr,
	Char:        CharArr,
	CIDR:        CIDRArr,
	Circle:      CircleArr,
	Date:        DateArr,
	DateRange:   DateRangeArr,
	Float4:      Float4Arr,
	Float8:      Float8Arr,
	Inet:        InetArr,
	Int2:        Int2Arr,
	Int2Vector:  Int2VectorArr,
	Int4:        Int4Arr,
	Int4Range:   Int4RangeArr,
	Int8:        Int8Arr,
	Int8Range:   Int8RangeArr,
	Interval:    IntervalArr,
	JSON:        JSONArr,
	JSONB:       JSONBArr,
	Line:        LineArr,
	LSeg:        LSegArr,
	MACAddr:     MACAddrArr,
	MACAddr8:    MACAddr8Arr,
	Money:       MoneyArr,
	Numeric:     NumericArr,
	NumRange:    NumRangeArr,
	Path:        PathArr,
	Point:       PointArr,
	Polygon:     PolygonArr,
	Text:        TextArr,
	Time:        TimeArr,
	Timestamp:   TimestampArr,
	Timestamptz: TimestamptzArr,
	Timetz:      TimetzArr,
	TSQuery:     TSQueryArr,
	TsRange:     TsRangeArr,
	TsTzRange:   TsTzRangeArr,
	TSVector:    TSVectorArr,
	UUID:        UUIDArr,
	VarBit:      VarBitArr,
	VarChar:     VarCharArr,
	XML:         XMLArr,
	// HStore:      HStoreArr,
}

var ArrayToElem = map[OID]OID{
	BitArr:         Bit,
	BoolArr:        Bool,
	BoxArr:         Box,
	BPCharArr:      BPChar,
	ByteaArr:       Bytea,
	CharArr:        Char,
	CIDRArr:        CIDR,
	CircleArr:      Circle,
	DateArr:        Date,
	DateRangeArr:   DateRange,
	Float4Arr:      Float4,
	Float8Arr:      Float8,
	InetArr:        Inet,
	Int2Arr:        Int2,
	Int2VectorArr:  Int2Vector,
	Int4Arr:        Int4,
	Int4RangeArr:   Int4Range,
	Int8Arr:        Int8,
	Int8RangeArr:   Int8Range,
	IntervalArr:    Interval,
	JSONArr:        JSON,
	JSONBArr:       JSONB,
	LineArr:        Line,
	LSegArr:        LSeg,
	MACAddrArr:     MACAddr,
	MACAddr8Arr:    MACAddr8,
	MoneyArr:       Money,
	NumericArr:     Numeric,
	NumRangeArr:    NumRange,
	PathArr:        Path,
	PointArr:       Point,
	PolygonArr:     Polygon,
	TextArr:        Text,
	TimeArr:        Time,
	TimestampArr:   Timestamp,
	TimestamptzArr: Timestamptz,
	TimetzArr:      Timetz,
	TSQueryArr:     TSQuery,
	TsRangeArr:     TsRange,
	TsTzRangeArr:   TsTzRange,
	TSVectorArr:    TSVector,
	UUIDArr:        UUID,
	VarBitArr:      VarBit,
	VarCharArr:     VarChar,
	XMLArr:         XML,
	// HStoreArr:      HStore,
}

var TypeToZeroValue = map[OID]string{
	// Any:        "",
	Bit:           `'0'`,
	BitArr:        `'{}'`,
	Bool:          `false`,
	BoolArr:       `'{}'`,
	Box:           `'(0,0),(0,0)'`,
	BoxArr:        `'{}'`,
	BPChar:        `''`,
	BPCharArr:     `'{}'`,
	Bytea:         `'\x'`,
	ByteaArr:      `'{}'`,
	Char:          `''`,
	CharArr:       `'{}'`,
	CIDR:          `'0.0.0.0/0'`,
	CIDRArr:       `'{}'`,
	Circle:        `'<(0,0),0>'`,
	CircleArr:     `'{}'`,
	Date:          `'0001-01-01'`,
	DateArr:       `'{}'`,
	DateRange:     `'(,)'`,
	DateRangeArr:  `'{}'`,
	Float4:        `0`,
	Float4Arr:     `'{}'`,
	Float8:        `0`,
	Float8Arr:     `'{}'`,
	Inet:          `'0.0.0.0'`,
	InetArr:       `'{}'`,
	Int2:          `0`,
	Int2Arr:       `'{}'`,
	Int2Vector:    `''`,
	Int2VectorArr: `'{}'`,
	Int4:          `0`,
	Int4Arr:       `'{}'`,
	Int4Range:     `'(,)'`,
	Int4RangeArr:  `'{}'`,
	Int8:          `0`,
	Int8Arr:       `'{}'`,
	Int8Range:     `'(,)'`,
	Int8RangeArr:  `'{}'`,
	Interval:      `'00:00:00'`,
	IntervalArr:   `'{}'`,
	JSON:          `'null'`,
	JSONArr:       `'{}'`,
	JSONB:         `'null'`,
	JSONBArr:      `'{}'`,
	// Line:          `''`,
	LineArr:        `'{}'`,
	LSeg:           `[(0,0),(0,0)]`,
	LSegArr:        `'{}'`,
	MACAddr:        `'00:00:00:00:00:00'`,
	MACAddrArr:     `'{}'`,
	MACAddr8:       `'00:00:00:00:00:00:00:00'`,
	MACAddr8Arr:    `'{}'`,
	Money:          `0`,
	MoneyArr:       `'{}'`,
	Numeric:        `0`,
	NumericArr:     `'{}'`,
	NumRange:       `'(,)'`,
	NumRangeArr:    `'{}'`,
	OIDVector:      `''`,
	Path:           `'((0,0))'`,
	PathArr:        `'{}'`,
	Point:          `'(0,0)'`,
	PointArr:       `'{}'`,
	Polygon:        `'((0,0))'`,
	PolygonArr:     `'{}'`,
	Text:           `''`,
	TextArr:        `'{}'`,
	Time:           `'00:00:00'`,
	TimeArr:        `'{}'`,
	Timestamp:      `'0001-01-01 00:00:00'`,
	TimestampArr:   `'{}'`,
	Timestamptz:    `'0001-01-01 00:00:00'`,
	TimestamptzArr: `'{}'`,
	Timetz:         `'00:00:00+00'`,
	TimetzArr:      `'{}'`,
	TSQuery:        `''`,
	TSQueryArr:     `'{}'`,
	TsRange:        `'(,)'`,
	TsRangeArr:     `'{}'`,
	TsTzRange:      `'(,)'`,
	TsTzRangeArr:   `'{}'`,
	TSVector:       `''`,
	TSVectorArr:    `'{}'`,
	// Unknown:        "",
	UUID:       `'00000000-0000-0000-0000-000000000000'`,
	UUIDArr:    `'{}'`,
	VarBit:     `''`,
	VarBitArr:  `'{}'`,
	VarChar:    `''`,
	VarCharArr: `'{}'`,
	XML:        `''`,
	XMLArr:     `'{}'`,
}
