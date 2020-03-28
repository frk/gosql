## convert

| postgres type  | go type           | Valuer                                 | Scanner                                |
|----------------|-------------------|----------------------------------------| ---------------------------------------|
| `bit(1)`       | `bool`            | `BitFromBool`                          | *native*                               |
|                | `uint8`           | *native*                               | *native*                               |
|                | `uint`            | *native*                               | *native*                               |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bit(1)[]`     | `[]bool`          | `BitArrayFromBoolSlice`                | `BitArrayToBoolSlice`                  |
|                | `[]uint8`         | `BitArrayFromUint8Slice`               | `BitArrayToUint8Slice`                 |
|                | `[]uint`          | `BitArrayFromUintSlice`                | `BitArrayToUintSlice`                  |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bool`         | `bool`            | *native*                               | *native*                               |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bool[]`       | `[]bool`          | `BoolArrayFromBoolSlice`               | `BoolArrayToBoolSlice`                 |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `box`          | `[2][2]float64`   | `BoxFromFloat64Array2Array2`           | `BoxToFloat64Array2Array2`             |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `box[]`        | `[][2][2]float64` | `BoxArrayFromFloat64Array2Array2Slice` | `BoxArrayToFloat64Array2Array2Slice`   |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bpchar(1)`    | `byte`            | `BPCharFromByte`                       | `BPCharToByte`                         |
|                | `rune`		  	 | `BPCharFromRune`                       | `BPCharToRune`                         |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bpchar(1)[]`  | `[]rune`          | `BPCharArrayFromRuneSlice` :warning:   | `BPCharArrayToRuneSlice`               |
|                | `[]string`        | `BPCharArrayFromStringSlice` :warning: | `BPCharArrayToStringSlice`             |
|                | `string`          | `BPCharArrayFromString` :warning:      | `BPCharArrayToString`                  |
|                | `[]byte`          | `BPCharArrayFromByteSlice` :warning:   | `BPCharArrayToByteSlice`               |
| `bytea`        | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `bytea[]`      | `[]string`        | `ByteaArrayFromStringSlice`            | `ByteaArrayToStringSlice`              |
|                | `[][]byte`        | `ByteaArrayFromByteSliceSlice`         | `ByteaArrayToByteSliceSlice`           |
| `char(1)`      | `byte`            | `CharFromByte`                         | `CharToByte`                           |
|                | `rune`			 | `CharFromRune`                         | `CharToRune`                           |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `char(1)[]`    | `[]rune`          | `CharArrayFromRuneSlice` :warning:     | `CharArrayToRuneSlice`                 |
|                | `[]string`        | `CharArrayFromStringSlice` :warning:   | `CharArrayToStringSlice`               |
|                | `string`          | `CharArrayFromString` :warning:        | `CharArrayToString`                    |
|                | `[]byte`          | `CharArrayFromByteSlice` :warning:     | `CharArrayToByteSlice`                 |
| `cidr`         | `net.IPNet`       | `CIDRFromIPNet`                        | `CIDRToIPNet`                          |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `cidr[]`       | `[]net.IPNet`     | `CIDRArrayFromIPNetSlice`              | `CIDRArrayToIPNetSlice`                |
|                | `string`          | *native*                               | *native*                               |
|                | `[]byte`          | *native*                               | *native*                               |
| `circle`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `circle[]`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `date`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `date[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `daterange`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `daterange[]`  | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `float4`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `float4[]`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `float8`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `float8[]`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `inet`         | `net.IPNet`       | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `inet[]`       | `[]net.IPNet`     | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int2`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int2[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int2vector`   | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int2vector[]` | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int4`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int4[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int4range`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int4range[]`  | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int8`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int8[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int8range`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `int8range[]`  | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `interval`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `interval[]`   | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `json`         | `interface{}`     | `JSON`                                 | `JSON`                                 |
|                | `[]byte`          | *native*                               | *native*                               |
|                | `string`          | *native*                               | *native*                               |
| `json[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `jsonb`        | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `jsonb[]`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `line`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `line[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `lseg`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `lseg[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `macaddr`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `macaddr[]`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `macaddr8`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `macaddr8[]`   | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `money`        | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `money[]`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `numeric`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `numeric[]`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `numrange`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `numrange[]`   | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `oidvector`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `path`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `path[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `point`        | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `point[]`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `polygon`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `polygon[]`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `text`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `text[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `time`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `time[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `timestamp`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `timestamp[]`  | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `timestamptz`  | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `timestamptz[]`| ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `timetz`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `timetz[]`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tsquery`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tsquery[]`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tsrange`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tsrange[]`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tstzrange`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tstzrange[]`  | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tsvector`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `tsvector[]`   | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `uuid`         | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `uuid[]`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `unknown`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `varbit`       | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `varbit[]`     | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `varchar`      | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `varchar[]`    | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `xml`          | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
| `xml[]`        | ???               | ???                                    | ???                                    |
|                | `string`          | ???                                    | ???                                    |
|                | `[]byte`          | ???                                    | ???                                    |
