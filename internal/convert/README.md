## convert

| postgres type  | go type           | Valuer                                 | Scanner                                |
|----------------|-------------------|----------------------------------------| ---------------------------------------|
| `bit(1)`       | `bool`            | `BitFromBool`                          | **native**                             |
|                | `uint8`           | **native**                             | **native**                             |
|                | `uint`            | **native**                             | **native**                             |
| `_bit(1)`      | `[]bool`          | `BitArrayFromBoolSlice`                | `BitArrayToBoolSlice`                  |
|                | `[]uint8`         | `BitArrayFromUint8Slice`               | `BitArrayToUint8Slice`                 |
|                | `[]uint`          | `BitArrayFromUintSlice`                | `BitArrayToUintSlice`                  |
| `bool`         | `bool`            | **native**                             | **native**                             |
| `_bool`        | `[]bool`          | `BoolArrayFromBoolSlice`               | `BoolArrayToBoolSlice`                 |
| `box`          | `[2][2]float64`   | `BoxFromFloat64Array2Array2`           | `BoxToFloat64Array2Array2`             |
| `_box`         | `[][2][2]float64` | `BoxArrayFromFloat64Array2Array2Slice` | `BoxArrayToFloat64Array2Array2Slice`   |
| `bpchar(1)`    | `byte`            | `BPCharFromByte`                       | `BPCharToByte`                         |
|                | `rune`			 | `BPCharFromRune`                       | `BPCharToRune`                         |
|                | `string`          | **native**                             | **native**                             |
|                | `[]byte`          | **native**                             | **native**                             |
| `_bpchar(1)`   | ???               | ???                                    | ???                                    |
| `bytea`        | ???               | ???                                    | ???                                    |
| `_bytea`       | ???               | ???                                    | ???                                    |
| `char(1)`      | `byte`            | `CharFromByte`                         | `CharToByte`                           |
|                | `rune`			 | `CharFromRune`                         | `CharToRune`                           |
|                | `string`          | **native**                             | **native**                             |
|                | `[]byte`          | **native**                             | **native**                             |
| `_char(1)`     | ???               | ???                                    | ???                                    |
| `cidr`         | ???               | ???                                    | ???                                    |
| `_cidr`        | ???               | ???                                    | ???                                    |
| `circle`       | ???               | ???                                    | ???                                    |
| `_circle`      | ???               | ???                                    | ???                                    |
| `date`         | ???               | ???                                    | ???                                    |
| `_date`        | ???               | ???                                    | ???                                    |
| `daterange`    | ???               | ???                                    | ???                                    |
| `_daterange`   | ???               | ???                                    | ???                                    |
| `float4`       | ???               | ???                                    | ???                                    |
| `_float4`      | ???               | ???                                    | ???                                    |
| `float8`       | ???               | ???                                    | ???                                    |
| `_float8`      | ???               | ???                                    | ???                                    |
| `inet`         | ???               | ???                                    | ???                                    |
| `_inet`        | ???               | ???                                    | ???                                    |
| `int2`         | ???               | ???                                    | ???                                    |
| `_int2`        | ???               | ???                                    | ???                                    |
| `int2vector`   | ???               | ???                                    | ???                                    |
| `_int2vector`  | ???               | ???                                    | ???                                    |
| `int4`         | ???               | ???                                    | ???                                    |
| `_int4`        | ???               | ???                                    | ???                                    |
| `int4range`    | ???               | ???                                    | ???                                    |
| `_int4range`   | ???               | ???                                    | ???                                    |
| `int8`         | ???               | ???                                    | ???                                    |
| `_int8`        | ???               | ???                                    | ???                                    |
| `int8range`    | ???               | ???                                    | ???                                    |
| `_int8range`   | ???               | ???                                    | ???                                    |
| `interval`     | ???               | ???                                    | ???                                    |
| `_interval`    | ???               | ???                                    | ???                                    |
| `json`         | `interface{}`     | `JSON`                                 | `JSON`                                 |
| `json`         | `[]byte`          | **native**                             | **native**                             |
| `json`         | `string`          | **native**                             | **native**                             |
| `_json`        | ???               | ???                                    | ???                                    |
| `jsonb`        | ???               | ???                                    | ???                                    |
| `_jsonb`       | ???               | ???                                    | ???                                    |
| `line`         | ???               | ???                                    | ???                                    |
| `_line`        | ???               | ???                                    | ???                                    |
| `lseg`         | ???               | ???                                    | ???                                    |
| `_lseg`        | ???               | ???                                    | ???                                    |
| `macaddr`      | ???               | ???                                    | ???                                    |
| `_macaddr`     | ???               | ???                                    | ???                                    |
| `macaddr8`     | ???               | ???                                    | ???                                    |
| `_macaddr8`    | ???               | ???                                    | ???                                    |
| `money`        | ???               | ???                                    | ???                                    |
| `_money`       | ???               | ???                                    | ???                                    |
| `numeric`      | ???               | ???                                    | ???                                    |
| `_numeric`     | ???               | ???                                    | ???                                    |
| `numrange`     | ???               | ???                                    | ???                                    |
| `_numrange`    | ???               | ???                                    | ???                                    |
| `oidvector`    | ???               | ???                                    | ???                                    |
| `path`         | ???               | ???                                    | ???                                    |
| `_path`        | ???               | ???                                    | ???                                    |
| `point`        | ???               | ???                                    | ???                                    |
| `_point`       | ???               | ???                                    | ???                                    |
| `polygon`      | ???               | ???                                    | ???                                    |
| `_polygon`     | ???               | ???                                    | ???                                    |
| `text`         | ???               | ???                                    | ???                                    |
| `_text`        | ???               | ???                                    | ???                                    |
| `time`         | ???               | ???                                    | ???                                    |
| `_time`        | ???               | ???                                    | ???                                    |
| `timestamp`    | ???               | ???                                    | ???                                    |
| `_timestamp`   | ???               | ???                                    | ???                                    |
| `timestamptz`  | ???               | ???                                    | ???                                    |
| `_timestamptz` | ???               | ???                                    | ???                                    |
| `timetz`       | ???               | ???                                    | ???                                    |
| `_timetz`      | ???               | ???                                    | ???                                    |
| `tsquery`      | ???               | ???                                    | ???                                    |
| `_tsquery`     | ???               | ???                                    | ???                                    |
| `tsrange`      | ???               | ???                                    | ???                                    |
| `_tsrange`     | ???               | ???                                    | ???                                    |
| `tstzrange`    | ???               | ???                                    | ???                                    |
| `_tstzrange`   | ???               | ???                                    | ???                                    |
| `tsvector`     | ???               | ???                                    | ???                                    |
| `_tsvector`    | ???               | ???                                    | ???                                    |
| `uuid`         | ???               | ???                                    | ???                                    |
| `_uuid`        | ???               | ???                                    | ???                                    |
| `unknown`      | ???               | ???                                    | ???                                    |
| `varbit`       | ???               | ???                                    | ???                                    |
| `_varbit`      | ???               | ???                                    | ???                                    |
| `varchar`      | ???               | ???                                    | ???                                    |
| `_varchar`     | ???               | ???                                    | ???                                    |
| `xml`          | ???               | ???                                    | ???                                    |
| `_xml`         | ???               | ???                                    | ???                                    |
