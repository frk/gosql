## convert

| postgres type  | go type           | Valuer                                 | Scanner                                |
|----------------|-------------------|----------------------------------------| ---------------------------------------|
| `bit(1)`       | `bool`            | `BitFromBool`                          | *native*                               |
|                | `uint8`           | *native*                               | *native*                               |
|                | `uint`            | *native*                               | *native*                               |
| `bit(1)[]`     | `[]bool`          | `BitArrayFromBoolSlice`                | `BitArrayToBoolSlice`                  |
|                | `[]uint8`         | `BitArrayFromUint8Slice`               | `BitArrayToUint8Slice`                 |
|                | `[]uint`          | `BitArrayFromUintSlice`                | `BitArrayToUintSlice`                  |
| `bool`         | `bool`            | *native*                               | *native*                               |
| `bool[]`       | `[]bool`          | `BoolArrayFromBoolSlice`               | `BoolArrayToBoolSlice`                 |
| `box`          | `[2][2]float64`   | `BoxFromFloat64Array2Array2`           | `BoxToFloat64Array2Array2`             |
| `box[]`        | `[][2][2]float64` | `BoxArrayFromFloat64Array2Array2Slice` | `BoxArrayToFloat64Array2Array2Slice`   |
| `bpchar(1)`    | `byte`            | `BPCharFromByte`                       | `BPCharToByte`                         |
|                | `rune`		  	 | `BPCharFromRune`                       | `BPCharToRune`                         |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bpchar(1)[]`  | `string`          | `BPCharArrayFromString`                | `BPCharArrayToString`                  |
|                | `[]byte`          | `BPCharArrayFromByteSlice`             | `BPCharArrayToByteSlice`               |
|                | `[]rune`          | `BPCharArrayFromRuneSlice`             | `BPCharArrayToRuneSlice`               |
|                | `[]string`        | `BPCharArrayFromStringSlice`           | `BPCharArrayToStringSlice`             |
| `bytea`        | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bytea[]`      | `[]string`        | `ByteaArrayFromStringSlice`            | `ByteaArrayToStringSlice`              |
|                | `[][]byte`        | `ByteaArrayFromByteSliceSlice`         | `ByteaArrayToByteSliceSlice`           |
| `char(1)`      | `byte`            | `CharFromByte`                         | `CharToByte`                           |
|                | `rune`			 | `CharFromRune`                         | `CharToRune`                           |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `char(1)[]`    | `string`          | `CharArrayFromString`                  | `CharArrayToString`                    |
|                | `[]byte`          | `CharArrayFromByteSlice`               | `CharArrayToByteSlice`                 |
|                | `[]rune`          | `CharArrayFromRuneSlice`               | `CharArrayToRuneSlice`                 |
|                | `[]string`        | `CharArrayFromStringSlice`             | `CharArrayToStringSlice`               |
| `cidr`         | ???               | ???                                    | ???                                    |
| `cidr[]`       | ???               | ???                                    | ???                                    |
| `circle`       | ???               | ???                                    | ???                                    |
| `circle[]`     | ???               | ???                                    | ???                                    |
| `date`         | ???               | ???                                    | ???                                    |
| `date[]`       | ???               | ???                                    | ???                                    |
| `daterange`    | ???               | ???                                    | ???                                    |
| `daterange[]`  | ???               | ???                                    | ???                                    |
| `float4`       | ???               | ???                                    | ???                                    |
| `float4[]`     | ???               | ???                                    | ???                                    |
| `float8`       | ???               | ???                                    | ???                                    |
| `float8[]`     | ???               | ???                                    | ???                                    |
| `inet`         | ???               | ???                                    | ???                                    |
| `inet[]`       | ???               | ???                                    | ???                                    |
| `int2`         | ???               | ???                                    | ???                                    |
| `int2[]`       | ???               | ???                                    | ???                                    |
| `int2vector`   | ???               | ???                                    | ???                                    |
| `int2vector[]` | ???               | ???                                    | ???                                    |
| `int4`         | ???               | ???                                    | ???                                    |
| `int4[]`       | ???               | ???                                    | ???                                    |
| `int4range`    | ???               | ???                                    | ???                                    |
| `int4range[]`  | ???               | ???                                    | ???                                    |
| `int8`         | ???               | ???                                    | ???                                    |
| `int8[]`       | ???               | ???                                    | ???                                    |
| `int8range`    | ???               | ???                                    | ???                                    |
| `int8range[]`  | ???               | ???                                    | ???                                    |
| `interval`     | ???               | ???                                    | ???                                    |
| `interval[]`   | ???               | ???                                    | ???                                    |
| `json`         | `interface{}`     | `JSON`                                 | `JSON`                                 |
|                | `[]byte`          | *native*                               | *native*                               |
|                | `string`          | *native*                               | *native*                               |
| `json[]`       | ???               | ???                                    | ???                                    |
| `jsonb`        | ???               | ???                                    | ???                                    |
| `jsonb[]`      | ???               | ???                                    | ???                                    |
| `line`         | ???               | ???                                    | ???                                    |
| `line[]`       | ???               | ???                                    | ???                                    |
| `lseg`         | ???               | ???                                    | ???                                    |
| `lseg[]`       | ???               | ???                                    | ???                                    |
| `macaddr`      | ???               | ???                                    | ???                                    |
| `macaddr[]`    | ???               | ???                                    | ???                                    |
| `macaddr8`     | ???               | ???                                    | ???                                    |
| `macaddr8[]`   | ???               | ???                                    | ???                                    |
| `money`        | ???               | ???                                    | ???                                    |
| `money[]`      | ???               | ???                                    | ???                                    |
| `numeric`      | ???               | ???                                    | ???                                    |
| `numeric[]`    | ???               | ???                                    | ???                                    |
| `numrange`     | ???               | ???                                    | ???                                    |
| `numrange[]`   | ???               | ???                                    | ???                                    |
| `oidvector`    | ???               | ???                                    | ???                                    |
| `path`         | ???               | ???                                    | ???                                    |
| `path[]`       | ???               | ???                                    | ???                                    |
| `point`        | ???               | ???                                    | ???                                    |
| `point[]`      | ???               | ???                                    | ???                                    |
| `polygon`      | ???               | ???                                    | ???                                    |
| `polygon[]`    | ???               | ???                                    | ???                                    |
| `text`         | ???               | ???                                    | ???                                    |
| `text[]`       | ???               | ???                                    | ???                                    |
| `time`         | ???               | ???                                    | ???                                    |
| `time[]`       | ???               | ???                                    | ???                                    |
| `timestamp`    | ???               | ???                                    | ???                                    |
| `timestamp[]`  | ???               | ???                                    | ???                                    |
| `timestamptz`  | ???               | ???                                    | ???                                    |
| `timestamptz[]`| ???               | ???                                    | ???                                    |
| `timetz`       | ???               | ???                                    | ???                                    |
| `timetz[]`     | ???               | ???                                    | ???                                    |
| `tsquery`      | ???               | ???                                    | ???                                    |
| `tsquery[]`    | ???               | ???                                    | ???                                    |
| `tsrange`      | ???               | ???                                    | ???                                    |
| `tsrange[]`    | ???               | ???                                    | ???                                    |
| `tstzrange`    | ???               | ???                                    | ???                                    |
| `tstzrange[]`  | ???               | ???                                    | ???                                    |
| `tsvector`     | ???               | ???                                    | ???                                    |
| `tsvector[]`   | ???               | ???                                    | ???                                    |
| `uuid`         | ???               | ???                                    | ???                                    |
| `uuid[]`       | ???               | ???                                    | ???                                    |
| `unknown`      | ???               | ???                                    | ???                                    |
| `varbit`       | ???               | ???                                    | ???                                    |
| `varbit[]`     | ???               | ???                                    | ???                                    |
| `varchar`      | ???               | ???                                    | ???                                    |
| `varchar[]`    | ???               | ???                                    | ???                                    |
| `xml`          | ???               | ???                                    | ???                                    |
| `xml[]`        | ???               | ???                                    | ???                                    |
