## convert

| postgres type  | go type                        | Valuer                                 | Scanner                                |
|----------------|--------------------------------|----------------------------------------| ---------------------------------------|
| `bit(1)`       | `bool`                         | `BitFromBool`                          | *native*                               |
|                | `uint8`                        | *native*                               | *native*                               |
|                | `uint`                         | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `bit(1)[]`     | `[]bool`                       | `BitArrayFromBoolSlice`                | `BitArrayToBoolSlice`                  |
|                | `[]uint8`                      | `BitArrayFromUint8Slice`               | `BitArrayToUint8Slice`                 |
|                | `[]uint`                       | `BitArrayFromUintSlice`                | `BitArrayToUintSlice`                  |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `bool`         | `bool`                         | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `bool[]`       | `[]bool`                       | `BoolArrayFromBoolSlice`               | `BoolArrayToBoolSlice`                 |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `box`          | `[2][2]float64`                | `BoxFromFloat64Array2Array2`           | `BoxToFloat64Array2Array2`             |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `box[]`        | `[][2][2]float64`              | `BoxArrayFromFloat64Array2Array2Slice` | `BoxArrayToFloat64Array2Array2Slice`   |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `bpchar(1)`    | `byte`                         | `BPCharFromByte`                       | `BPCharToByte`                         |
|                | `rune`		  	              | `BPCharFromRune`                       | `BPCharToRune`                         |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `bpchar(1)[]`  | `[]rune`                       | `BPCharArrayFromRuneSlice` :warning:   | `BPCharArrayToRuneSlice`               |
|                | `[]string`                     | `BPCharArrayFromStringSlice` :warning: | `BPCharArrayToStringSlice`             |
|                | `string`                       | `BPCharArrayFromString` :warning:      | `BPCharArrayToString`                  |
|                | `[]byte`                       | `BPCharArrayFromByteSlice` :warning:   | `BPCharArrayToByteSlice`               |
| `bytea`        | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `bytea[]`      | `[]string`                     | `ByteaArrayFromStringSlice`            | `ByteaArrayToStringSlice`              |
|                | `[][]byte`                     | `ByteaArrayFromByteSliceSlice`         | `ByteaArrayToByteSliceSlice`           |
| `char(1)`      | `byte`                         | `CharFromByte`                         | `CharToByte`                           |
|                | `rune`			              | `CharFromRune`                         | `CharToRune`                           |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `char(1)[]`    | `[]rune`                       | `CharArrayFromRuneSlice` :warning:     | `CharArrayToRuneSlice`                 |
|                | `[]string`                     | `CharArrayFromStringSlice` :warning:   | `CharArrayToStringSlice`               |
|                | `string`                       | `CharArrayFromString` :warning:        | `CharArrayToString`                    |
|                | `[]byte`                       | `CharArrayFromByteSlice` :warning:     | `CharArrayToByteSlice`                 |
| `cidr`         | `net.IPNet`                    | `CIDRFromIPNet`                        | `CIDRToIPNet`                          |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `cidr[]`       | `[]net.IPNet`                  | `CIDRArrayFromIPNetSlice`              | `CIDRArrayToIPNetSlice`                |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `circle`       | ---                            | ---                                    | ---                                    |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `circle[]`     | ---                            | ---                                    | ---                                    |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `date`         | `time.Time`                    | *native*                               | `DateToTime`                           |
|                | `string`                       | *native*                               | `DateToString`                         |
|                | `[]byte`                       | *native*                               | `DateToByteSlice`                      |
| `date[]`       | `[]time.Time`                  | `DateArrayFromTimeSlice`               | `DateArrayToTimeSlice`                 |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `daterange`    | `[2]time.Time`                 | `DateRangeFromTimeArray2`              | `DateRangeToTimeArray2`                |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `daterange[]`  | `[][2]time.Time`               | `DateRangeArrayFromTimeArray2Slice`    | `DateRangeArrayToTimeArray2Slice`      |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `float4`       | `float32`                      | *native*                               | *native*                               |
|                | `float64`                      | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `float4[]`     | `[]float32`                    | `Float4ArrayFromFloat32Slice`          | `Float4ArrayToFloat32Slice`            |
|                | `[]float64`                    | `Float4ArrayFromFloat64Slice`          | `Float4ArrayToFloat64Slice`            |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `float8`       | `float64`                      | *native*                               | *native*                               |
|                | `float32`                      | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `float8[]`     | `[]float64`                    | `Float8ArrayFromFloat64Slice`          | `Float8ArrayToFloat64Slice`            |
|                | `[]float32`                    | `Float8ArrayFromFloat32Slice`          | `Float8ArrayToFloat32Slice`            |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `hstore`       | `map[string]string`            | `HStoreFromStringMap`                  | `HStoreToStringMap`                    |
|                | `map[string]*string`           | `HStoreFromStringPtrMap`               | `HStoreToStringPtrMap`                 |
|                | `map[string]sql.NullString`    | `HStoreFromNullStringMap`              | `HStoreToNullStringMap`                |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `hstore[]`     | `[]map[string]string`          | `HStoreArrayFromStringMapSlice`        | `HStoreArrayToStringMapSlice`          |
|                | `[]map[string]*string`         | `HStoreArrayFromStringPtrMapSlice`     | `HStoreArrayToStringPtrMapSlice`       |
|                | `[]map[string]sql.NullString`  | `HStoreArrayFromNullStringMapSlice`    | `HStoreArrayToNullStringMapSlice`      |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `inet`         | `net.IPNet`                    | `InetFromIPNet`                        | `InetToIPNet`                          |
|                | `net.IP`                       | `InetFromIP`                           | `InetToIP`                             |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `inet[]`       | `[]net.IPNet`                  | `InetArrayFromIPNetSlice`              | `InetArrayToIPNetSlice`                |
|                | `[]net.IP`                     | `InetArrayFromIPSlice`                 | `InetArrayToIPSlice`                   |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int2`         | `int`                          | *native*                               | *native*                               |
|                | `int8`                         | *native*                               | *native*                               |
|                | `int16`                        | *native*                               | *native*                               |
|                | `int32`                        | *native*                               | *native*                               |
|                | `int64`                        | *native*                               | *native*                               |
|                | `uint`                         | *native*                               | *native*                               |
|                | `uint8`                        | *native*                               | *native*                               |
|                | `uint16`                       | *native*                               | *native*                               |
|                | `uint32`                       | *native*                               | *native*                               |
|                | `uint64`                       | *native*                               | *native*                               |
|                | `float32`                      | *native*                               | *native*                               |
|                | `float64`                      | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int2[]`       | `[]int`                        | `Int2ArrayFromIntSlice`                | `Int2ArrayToIntSlice`                  |
|                | `[]int8`                       | `Int2ArrayFromInt8Slice`               | `Int2ArrayToInt8Slice`                 |
|                | `[]int16`                      | `Int2ArrayFromInt16Slice`              | `Int2ArrayToInt16Slice`                |
|                | `[]int32`                      | `Int2ArrayFromInt32Slice`              | `Int2ArrayToInt32Slice`                |
|                | `[]int64`                      | `Int2ArrayFromInt64Slice`              | `Int2ArrayToInt64Slice`                |
|                | `[]uint`                       | `Int2ArrayFromUintSlice`               | `Int2ArrayToUintSlice`                 |
|                | `[]uint8`                      | `Int2ArrayFromUint8Slice`              | `Int2ArrayToUint8Slice`                |
|                | `[]uint16`                     | `Int2ArrayFromUint16Slice`             | `Int2ArrayToUint16Slice`               |
|                | `[]uint32`                     | `Int2ArrayFromUint32Slice`             | `Int2ArrayToUint32Slice`               |
|                | `[]uint64`                     | `Int2ArrayFromUint64Slice`             | `Int2ArrayToUint64Slice`               |
|                | `[]float32`                    | `Int2ArrayFromFloat32Slice`            | `Int2ArrayToFloat32Slice`              |
|                | `[]float64`                    | `Int2ArrayFromFloat64Slice`            | `Int2ArrayToFloat64Slice`              |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int2vector`   | `[]int`                        | `Int2VectorFromIntSlice`               | `Int2VectorToIntSlice`                 |
|                | `[]int8`                       | `Int2VectorFromInt8Slice`              | `Int2VectorToInt8Slice`                |
|                | `[]int16`                      | `Int2VectorFromInt16Slice`             | `Int2VectorToInt16Slice`               |
|                | `[]int32`                      | `Int2VectorFromInt32Slice`             | `Int2VectorToInt32Slice`               |
|                | `[]int64`                      | `Int2VectorFromInt64Slice`             | `Int2VectorToInt64Slice`               |
|                | `[]uint`                       | `Int2VectorFromUintSlice`              | `Int2VectorToUintSlice`                |
|                | `[]uint8`                      | `Int2VectorFromUint8Slice`             | `Int2VectorToUint8Slice`               |
|                | `[]uint16`                     | `Int2VectorFromUint16Slice`            | `Int2VectorToUint16Slice`              |
|                | `[]uint32`                     | `Int2VectorFromUint32Slice`            | `Int2VectorToUint32Slice`              |
|                | `[]uint64`                     | `Int2VectorFromUint64Slice`            | `Int2VectorToUint64Slice`              |
|                | `[]float32`                    | `Int2VectorFromFloat32Slice`           | `Int2VectorToFloat32Slice`             |
|                | `[]float64`                    | `Int2VectorFromFloat64Slice`           | `Int2VectorToFloat64Slice`             |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int2vector[]` | `[][]int`                      | `Int2VectorArrayFromIntSlice`          | `Int2VectorArrayToIntSlice`            |
|                | `[][]int8`                     | `Int2VectorArrayFromInt8Slice`         | `Int2VectorArrayToInt8Slice`           |
|                | `[][]int16`                    | `Int2VectorArrayFromInt16Slice`        | `Int2VectorArrayToInt16Slice`          |
|                | `[][]int32`                    | `Int2VectorArrayFromInt32Slice`        | `Int2VectorArrayToInt32Slice`          |
|                | `[][]int64`                    | `Int2VectorArrayFromInt64Slice`        | `Int2VectorArrayToInt64Slice`          |
|                | `[][]uint`                     | `Int2VectorArrayFromUintSlice`         | `Int2VectorArrayToUintSlice`           |
|                | `[][]uint8`                    | `Int2VectorArrayFromUint8Slice`        | `Int2VectorArrayToUint8Slice`          |
|                | `[][]uint16`                   | `Int2VectorArrayFromUint16Slice`       | `Int2VectorArrayToUint16Slice`         |
|                | `[][]uint32`                   | `Int2VectorArrayFromUint32Slice`       | `Int2VectorArrayToUint32Slice`         |
|                | `[][]uint64`                   | `Int2VectorArrayFromUint64Slice`       | `Int2VectorArrayToUint64Slice`         |
|                | `[][]float32`                  | `Int2VectorArrayFromFloat32Slice`      | `Int2VectorArrayToFloat32Slice`        |
|                | `[][]float64`                  | `Int2VectorArrayFromFloat64Slice`      | `Int2VectorArrayToFloat64Slice`        |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int4`         | `int`                          | *native*                               | *native*                               |
|                | `int8`                         | *native*                               | *native*                               |
|                | `int16`                        | *native*                               | *native*                               |
|                | `int32`                        | *native*                               | *native*                               |
|                | `int64`                        | *native*                               | *native*                               |
|                | `uint`                         | *native*                               | *native*                               |
|                | `uint8`                        | *native*                               | *native*                               |
|                | `uint16`                       | *native*                               | *native*                               |
|                | `uint32`                       | *native*                               | *native*                               |
|                | `uint64`                       | *native*                               | *native*                               |
|                | `float32`                      | *native*                               | *native*                               |
|                | `float64`                      | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int4[]`       | `[]int`                        | `Int4ArrayFromIntSlice`                | `Int4ArrayToIntSlice`                  |
|                | `[]int8`                       | `Int4ArrayFromInt8Slice`               | `Int4ArrayToInt8Slice`                 |
|                | `[]int16`                      | `Int4ArrayFromInt16Slice`              | `Int4ArrayToInt16Slice`                |
|                | `[]int32`                      | `Int4ArrayFromInt32Slice`              | `Int4ArrayToInt32Slice`                |
|                | `[]int64`                      | `Int4ArrayFromInt64Slice`              | `Int4ArrayToInt64Slice`                |
|                | `[]uint`                       | `Int4ArrayFromUintSlice`               | `Int4ArrayToUintSlice`                 |
|                | `[]uint8`                      | `Int4ArrayFromUint8Slice`              | `Int4ArrayToUint8Slice`                |
|                | `[]uint16`                     | `Int4ArrayFromUint16Slice`             | `Int4ArrayToUint16Slice`               |
|                | `[]uint32`                     | `Int4ArrayFromUint32Slice`             | `Int4ArrayToUint32Slice`               |
|                | `[]uint64`                     | `Int4ArrayFromUint64Slice`             | `Int4ArrayToUint64Slice`               |
|                | `[]float32`                    | `Int4ArrayFromFloat32Slice`            | `Int4ArrayToFloat32Slice`              |
|                | `[]float64`                    | `Int4ArrayFromFloat64Slice`            | `Int4ArrayToFloat64Slice`              |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int4range`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `int4range[]`  | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `int8`         | `int`                          | *native*                               | *native*                               |
|                | `int8`                         | *native*                               | *native*                               |
|                | `int16`                        | *native*                               | *native*                               |
|                | `int32`                        | *native*                               | *native*                               |
|                | `int64`                        | *native*                               | *native*                               |
|                | `uint`                         | *native*                               | *native*                               |
|                | `uint8`                        | *native*                               | *native*                               |
|                | `uint16`                       | *native*                               | *native*                               |
|                | `uint32`                       | *native*                               | *native*                               |
|                | `uint64`                       | *native*                               | *native*                               |
|                | `float32`                      | *native*                               | *native*                               |
|                | `float64`                      | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int8[]`       | `[]int`                        | `Int8ArrayFromIntSlice`                | `Int8ArrayToIntSlice`                  |
|                | `[]int8`                       | `Int8ArrayFromInt8Slice`               | `Int8ArrayToInt8Slice`                 |
|                | `[]int16`                      | `Int8ArrayFromInt16Slice`              | `Int8ArrayToInt16Slice`                |
|                | `[]int32`                      | `Int8ArrayFromInt32Slice`              | `Int8ArrayToInt32Slice`                |
|                | `[]int64`                      | `Int8ArrayFromInt64Slice`              | `Int8ArrayToInt64Slice`                |
|                | `[]uint`                       | `Int8ArrayFromUintSlice`               | `Int8ArrayToUintSlice`                 |
|                | `[]uint8`                      | `Int8ArrayFromUint8Slice`              | `Int8ArrayToUint8Slice`                |
|                | `[]uint16`                     | `Int8ArrayFromUint16Slice`             | `Int8ArrayToUint16Slice`               |
|                | `[]uint32`                     | `Int8ArrayFromUint32Slice`             | `Int8ArrayToUint32Slice`               |
|                | `[]uint64`                     | `Int8ArrayFromUint64Slice`             | `Int8ArrayToUint64Slice`               |
|                | `[]float32`                    | `Int8ArrayFromFloat32Slice`            | `Int8ArrayToFloat32Slice`              |
|                | `[]float64`                    | `Int8ArrayFromFloat64Slice`            | `Int8ArrayToFloat64Slice`              |
|                | `string`                       | *native*                               | *native*                               |
|                | `[]byte`                       | *native*                               | *native*                               |
| `int8range`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `int8range[]`  | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `interval`     | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `interval[]`   | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `json`         | `interface{}`                  | `JSON`                                 | `JSON`                                 |
|                | `[]byte`                       | *native*                               | *native*                               |
|                | `string`                       | *native*                               | *native*                               |
| `json[]`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `jsonb`        | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `jsonb[]`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `line`         | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `line[]`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `lseg`         | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `lseg[]`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `macaddr`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `macaddr[]`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `macaddr8`     | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `macaddr8[]`   | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `money`        | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `money[]`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `numeric`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `numeric[]`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `numrange`     | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `numrange[]`   | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `oidvector`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `path`         | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `path[]`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `point`        | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `point[]`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `polygon`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `polygon[]`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `text`         | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `text[]`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `time`         | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `time[]`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `timestamp`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `timestamp[]`  | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `timestamptz`  | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `timestamptz[]`| ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `timetz`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `timetz[]`     | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tsquery`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tsquery[]`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tsrange`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tsrange[]`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tstzrange`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tstzrange[]`  | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tsvector`     | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `tsvector[]`   | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `uuid`         | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `uuid[]`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `unknown`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `varbit`       | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `varbit[]`     | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `varchar`      | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `varchar[]`    | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `xml`          | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
| `xml[]`        | ???                            | ???                                    | ???                                    |
|                | `string`                       | ???                                    | ???                                    |
|                | `[]byte`                       | ???                                    | ???                                    |
