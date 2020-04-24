## pgsql

| postgres type  | go type                        | Valuer                                     | Scanner                                  |
|----------------|--------------------------------|--------------------------------------------| -----------------------------------------|
| `bit(1)`       | `bool`                         | `BitFromBool`                              | *native*                                 |
|                | `uint8`                        | *native*                                   | *native*                                 |
|                | `uint`                         | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `bit(1)[]`     | `[]bool`                       | `BitArrayFromBoolSlice`                    | `BitArrayToBoolSlice`                    |
|                | `[]uint8`                      | `BitArrayFromUint8Slice`                   | `BitArrayToUint8Slice`                   |
|                | `[]uint`                       | `BitArrayFromUintSlice`                    | `BitArrayToUintSlice`                    |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `bool`         | `bool`                         | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `bool[]`       | `[]bool`                       | `BoolArrayFromBoolSlice`                   | `BoolArrayToBoolSlice`                   |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `box`          | `[2][2]float64`                | `BoxFromFloat64Array2Array2`               | `BoxToFloat64Array2Array2`               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `box[]`        | `[][2][2]float64`              | `BoxArrayFromFloat64Array2Array2Slice`     | `BoxArrayToFloat64Array2Array2Slice`     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `bpchar(1)`    | `byte`                         | `BPCharFromByte`                           | `BPCharToByte`                           |
|                | `rune`		  	              | `BPCharFromRune`                           | `BPCharToRune`                           |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `bpchar(1)[]`  | `[]rune`                       | `BPCharArrayFromRuneSlice` :warning:       | `BPCharArrayToRuneSlice`                 |
|                | `[]string`                     | `BPCharArrayFromStringSlice` :warning:     | `BPCharArrayToStringSlice`               |
|                | `string`                       | `BPCharArrayFromString` :warning:          | `BPCharArrayToString`                    |
|                | `[]byte`                       | `BPCharArrayFromByteSlice` :warning:       | `BPCharArrayToByteSlice`                 |
| `bytea`        | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `bytea[]`      | `[]string`                     | `ByteaArrayFromStringSlice`                | `ByteaArrayToStringSlice`                |
|                | `[][]byte`                     | `ByteaArrayFromByteSliceSlice`             | `ByteaArrayToByteSliceSlice`             |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `char(1)`      | `byte`                         | `CharFromByte`                             | `CharToByte`                             |
|                | `rune`			              | `CharFromRune`                             | `CharToRune`                             |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `char(1)[]`    | `[]rune`                       | `CharArrayFromRuneSlice` :warning:         | `CharArrayToRuneSlice`                   |
|                | `[]string`                     | `CharArrayFromStringSlice` :warning:       | `CharArrayToStringSlice`                 |
|                | `string`                       | `CharArrayFromString` :warning:            | `CharArrayToString`                      |
|                | `[]byte`                       | `CharArrayFromByteSlice` :warning:         | `CharArrayToByteSlice`                   |
| `cidr`         | `net.IPNet`                    | `CIDRFromIPNet`                            | `CIDRToIPNet`                            |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `cidr[]`       | `[]net.IPNet`                  | `CIDRArrayFromIPNetSlice`                  | `CIDRArrayToIPNetSlice`                  |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `circle`       | ---                            | ---                                        | ---                                      |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `circle[]`     | ---                            | ---                                        | ---                                      |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `date`         | `time.Time`                    | *native*                                   | `DateToTime`                             |
|                | `string`                       | *native*                                   | `DateToString`                           |
|                | `[]byte`                       | *native*                                   | `DateToByteSlice`                        |
| `date[]`       | `[]time.Time`                  | `DateArrayFromTimeSlice`                   | `DateArrayToTimeSlice`                   |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `daterange`    | `[2]time.Time`                 | `DateRangeFromTimeArray2`                  | `DateRangeToTimeArray2`                  |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `daterange[]`  | `[][2]time.Time`               | `DateRangeArrayFromTimeArray2Slice`        | `DateRangeArrayToTimeArray2Slice`        |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `float4`       | `float32`                      | *native*                                   | *native*                                 |
|                | `float64`                      | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `float4[]`     | `[]float32`                    | `Float4ArrayFromFloat32Slice`              | `Float4ArrayToFloat32Slice`              |
|                | `[]float64`                    | `Float4ArrayFromFloat64Slice`              | `Float4ArrayToFloat64Slice`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `float8`       | `float64`                      | *native*                                   | *native*                                 |
|                | `float32`                      | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `float8[]`     | `[]float64`                    | `Float8ArrayFromFloat64Slice`              | `Float8ArrayToFloat64Slice`              |
|                | `[]float32`                    | `Float8ArrayFromFloat32Slice`              | `Float8ArrayToFloat32Slice`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `hstore`       | `map[string]string`            | `HStoreFromStringMap`                      | `HStoreToStringMap`                      |
|                | `map[string]*string`           | `HStoreFromStringPtrMap`                   | `HStoreToStringPtrMap`                   |
|                | `map[string]sql.NullString`    | `HStoreFromNullStringMap`                  | `HStoreToNullStringMap`                  |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `hstore[]`     | `[]map[string]string`          | `HStoreArrayFromStringMapSlice`            | `HStoreArrayToStringMapSlice`            |
|                | `[]map[string]*string`         | `HStoreArrayFromStringPtrMapSlice`         | `HStoreArrayToStringPtrMapSlice`         |
|                | `[]map[string]sql.NullString`  | `HStoreArrayFromNullStringMapSlice`        | `HStoreArrayToNullStringMapSlice`        |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `inet`         | `net.IPNet`                    | `InetFromIPNet`                            | `InetToIPNet`                            |
|                | `net.IP`                       | `InetFromIP`                               | `InetToIP`                               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `inet[]`       | `[]net.IPNet`                  | `InetArrayFromIPNetSlice`                  | `InetArrayToIPNetSlice`                  |
|                | `[]net.IP`                     | `InetArrayFromIPSlice`                     | `InetArrayToIPSlice`                     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int2`         | `int`                          | *native*                                   | *native*                                 |
|                | `int8`                         | *native*                                   | *native*                                 |
|                | `int16`                        | *native*                                   | *native*                                 |
|                | `int32`                        | *native*                                   | *native*                                 |
|                | `int64`                        | *native*                                   | *native*                                 |
|                | `uint`                         | *native*                                   | *native*                                 |
|                | `uint8`                        | *native*                                   | *native*                                 |
|                | `uint16`                       | *native*                                   | *native*                                 |
|                | `uint32`                       | *native*                                   | *native*                                 |
|                | `uint64`                       | *native*                                   | *native*                                 |
|                | `float32`                      | *native*                                   | *native*                                 |
|                | `float64`                      | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int2[]`       | `[]int`                        | `Int2ArrayFromIntSlice`                    | `Int2ArrayToIntSlice`                    |
|                | `[]int8`                       | `Int2ArrayFromInt8Slice`                   | `Int2ArrayToInt8Slice`                   |
|                | `[]int16`                      | `Int2ArrayFromInt16Slice`                  | `Int2ArrayToInt16Slice`                  |
|                | `[]int32`                      | `Int2ArrayFromInt32Slice`                  | `Int2ArrayToInt32Slice`                  |
|                | `[]int64`                      | `Int2ArrayFromInt64Slice`                  | `Int2ArrayToInt64Slice`                  |
|                | `[]uint`                       | `Int2ArrayFromUintSlice`                   | `Int2ArrayToUintSlice`                   |
|                | `[]uint8`                      | `Int2ArrayFromUint8Slice`                  | `Int2ArrayToUint8Slice`                  |
|                | `[]uint16`                     | `Int2ArrayFromUint16Slice`                 | `Int2ArrayToUint16Slice`                 |
|                | `[]uint32`                     | `Int2ArrayFromUint32Slice`                 | `Int2ArrayToUint32Slice`                 |
|                | `[]uint64`                     | `Int2ArrayFromUint64Slice`                 | `Int2ArrayToUint64Slice`                 |
|                | `[]float32`                    | `Int2ArrayFromFloat32Slice`                | `Int2ArrayToFloat32Slice`                |
|                | `[]float64`                    | `Int2ArrayFromFloat64Slice`                | `Int2ArrayToFloat64Slice`                |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int2vector`   | `[]int`                        | `Int2VectorFromIntSlice`                   | `Int2VectorToIntSlice`                   |
|                | `[]int8`                       | `Int2VectorFromInt8Slice`                  | `Int2VectorToInt8Slice`                  |
|                | `[]int16`                      | `Int2VectorFromInt16Slice`                 | `Int2VectorToInt16Slice`                 |
|                | `[]int32`                      | `Int2VectorFromInt32Slice`                 | `Int2VectorToInt32Slice`                 |
|                | `[]int64`                      | `Int2VectorFromInt64Slice`                 | `Int2VectorToInt64Slice`                 |
|                | `[]uint`                       | `Int2VectorFromUintSlice`                  | `Int2VectorToUintSlice`                  |
|                | `[]uint8`                      | `Int2VectorFromUint8Slice`                 | `Int2VectorToUint8Slice`                 |
|                | `[]uint16`                     | `Int2VectorFromUint16Slice`                | `Int2VectorToUint16Slice`                |
|                | `[]uint32`                     | `Int2VectorFromUint32Slice`                | `Int2VectorToUint32Slice`                |
|                | `[]uint64`                     | `Int2VectorFromUint64Slice`                | `Int2VectorToUint64Slice`                |
|                | `[]float32`                    | `Int2VectorFromFloat32Slice`               | `Int2VectorToFloat32Slice`               |
|                | `[]float64`                    | `Int2VectorFromFloat64Slice`               | `Int2VectorToFloat64Slice`               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int2vector[]` | `[][]int`                      | `Int2VectorArrayFromIntSliceSlice`         | `Int2VectorArrayToIntSliceSlice`         |
|                | `[][]int8`                     | `Int2VectorArrayFromInt8SliceSlice`        | `Int2VectorArrayToInt8SliceSlice`        |
|                | `[][]int16`                    | `Int2VectorArrayFromInt16SliceSlice`       | `Int2VectorArrayToInt16SliceSlice`       |
|                | `[][]int32`                    | `Int2VectorArrayFromInt32SliceSlice`       | `Int2VectorArrayToInt32SliceSlice`       |
|                | `[][]int64`                    | `Int2VectorArrayFromInt64SliceSlice`       | `Int2VectorArrayToInt64SliceSlice`       |
|                | `[][]uint`                     | `Int2VectorArrayFromUintSliceSlice`        | `Int2VectorArrayToUintSliceSlice`        |
|                | `[][]uint8`                    | `Int2VectorArrayFromUint8SliceSlice`       | `Int2VectorArrayToUint8SliceSlice`       |
|                | `[][]uint16`                   | `Int2VectorArrayFromUint16SliceSlice`      | `Int2VectorArrayToUint16SliceSlice`      |
|                | `[][]uint32`                   | `Int2VectorArrayFromUint32SliceSlice`      | `Int2VectorArrayToUint32SliceSlice`      |
|                | `[][]uint64`                   | `Int2VectorArrayFromUint64SliceSlice`      | `Int2VectorArrayToUint64SliceSlice`      |
|                | `[][]float32`                  | `Int2VectorArrayFromFloat32SliceSlice`     | `Int2VectorArrayToFloat32SliceSlice`     |
|                | `[][]float64`                  | `Int2VectorArrayFromFloat64SliceSlice`     | `Int2VectorArrayToFloat64SliceSlice`     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int4`         | `int`                          | *native*                                   | *native*                                 |
|                | `int8`                         | *native*                                   | *native*                                 |
|                | `int16`                        | *native*                                   | *native*                                 |
|                | `int32`                        | *native*                                   | *native*                                 |
|                | `int64`                        | *native*                                   | *native*                                 |
|                | `uint`                         | *native*                                   | *native*                                 |
|                | `uint8`                        | *native*                                   | *native*                                 |
|                | `uint16`                       | *native*                                   | *native*                                 |
|                | `uint32`                       | *native*                                   | *native*                                 |
|                | `uint64`                       | *native*                                   | *native*                                 |
|                | `float32`                      | *native*                                   | *native*                                 |
|                | `float64`                      | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int4[]`       | `[]int`                        | `Int4ArrayFromIntSlice`                    | `Int4ArrayToIntSlice`                    |
|                | `[]int8`                       | `Int4ArrayFromInt8Slice`                   | `Int4ArrayToInt8Slice`                   |
|                | `[]int16`                      | `Int4ArrayFromInt16Slice`                  | `Int4ArrayToInt16Slice`                  |
|                | `[]int32`                      | `Int4ArrayFromInt32Slice`                  | `Int4ArrayToInt32Slice`                  |
|                | `[]int64`                      | `Int4ArrayFromInt64Slice`                  | `Int4ArrayToInt64Slice`                  |
|                | `[]uint`                       | `Int4ArrayFromUintSlice`                   | `Int4ArrayToUintSlice`                   |
|                | `[]uint8`                      | `Int4ArrayFromUint8Slice`                  | `Int4ArrayToUint8Slice`                  |
|                | `[]uint16`                     | `Int4ArrayFromUint16Slice`                 | `Int4ArrayToUint16Slice`                 |
|                | `[]uint32`                     | `Int4ArrayFromUint32Slice`                 | `Int4ArrayToUint32Slice`                 |
|                | `[]uint64`                     | `Int4ArrayFromUint64Slice`                 | `Int4ArrayToUint64Slice`                 |
|                | `[]float32`                    | `Int4ArrayFromFloat32Slice`                | `Int4ArrayToFloat32Slice`                |
|                | `[]float64`                    | `Int4ArrayFromFloat64Slice`                | `Int4ArrayToFloat64Slice`                |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int4range`    | `[2]int`                       | `Int4RangeFromIntArray2`                   | `Int4RangeToIntArray2`                   |
|                | `[2]int8`                      | `Int4RangeFromInt8Array2`                  | `Int4RangeToInt8Array2`                  |
|                | `[2]int16`                     | `Int4RangeFromInt16Array2`                 | `Int4RangeToInt16Array2`                 |
|                | `[2]int32`                     | `Int4RangeFromInt32Array2`                 | `Int4RangeToInt32Array2`                 |
|                | `[2]int64`                     | `Int4RangeFromInt64Array2`                 | `Int4RangeToInt64Array2`                 |
|                | `[2]uint`                      | `Int4RangeFromUintArray2`                  | `Int4RangeToUintArray2`                  |
|                | `[2]uint8`                     | `Int4RangeFromUint8Array2`                 | `Int4RangeToUint8Array2`                 |
|                | `[2]uint16`                    | `Int4RangeFromUint16Array2`                | `Int4RangeToUint16Array2`                |
|                | `[2]uint32`                    | `Int4RangeFromUint32Array2`                | `Int4RangeToUint32Array2`                |
|                | `[2]uint64`                    | `Int4RangeFromUint64Array2`                | `Int4RangeToUint64Array2`                |
|                | `[2]float32`                   | `Int4RangeFromFloat32Array2`               | `Int4RangeToFloat32Array2`               |
|                | `[2]float64`                   | `Int4RangeFromFloat64Array2`               | `Int4RangeToFloat64Array2`               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int4range[]`  | `[][2]int`                     | `Int4RangeArrayFromIntArray2Slice`         | `Int4RangeArrayToIntArray2Slice`         |    
|                | `[][2]int8`                    | `Int4RangeArrayFromInt8Array2Slice`        | `Int4RangeArrayToInt8Array2Slice`        |
|                | `[][2]int16`                   | `Int4RangeArrayFromInt16Array2Slice`       | `Int4RangeArrayToInt16Array2Slice`       |
|                | `[][2]int32`                   | `Int4RangeArrayFromInt32Array2Slice`       | `Int4RangeArrayToInt32Array2Slice`       |
|                | `[][2]int64`                   | `Int4RangeArrayFromInt64Array2Slice`       | `Int4RangeArrayToInt64Array2Slice`       |
|                | `[][2]uint`                    | `Int4RangeArrayFromUintArray2Slice`        | `Int4RangeArrayToUintArray2Slice`        |
|                | `[][2]uint8`                   | `Int4RangeArrayFromUint8Array2Slice`       | `Int4RangeArrayToUint8Array2Slice`       |
|                | `[][2]uint16`                  | `Int4RangeArrayFromUint16Array2Slice`      | `Int4RangeArrayToUint16Array2Slice`      |
|                | `[][2]uint32`                  | `Int4RangeArrayFromUint32Array2Slice`      | `Int4RangeArrayToUint32Array2Slice`      |
|                | `[][2]uint64`                  | `Int4RangeArrayFromUint64Array2Slice`      | `Int4RangeArrayToUint64Array2Slice`      |
|                | `[][2]float32`                 | `Int4RangeArrayFromFloat32Array2Slice`     | `Int4RangeArrayToFloat32Array2Slice`     |
|                | `[][2]float64`                 | `Int4RangeArrayFromFloat64Array2Slice`     | `Int4RangeArrayToFloat64Array2Slice`     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int8`         | `int`                          | *native*                                   | *native*                                 |
|                | `int8`                         | *native*                                   | *native*                                 |
|                | `int16`                        | *native*                                   | *native*                                 |
|                | `int32`                        | *native*                                   | *native*                                 |
|                | `int64`                        | *native*                                   | *native*                                 |
|                | `uint`                         | *native*                                   | *native*                                 |
|                | `uint8`                        | *native*                                   | *native*                                 |
|                | `uint16`                       | *native*                                   | *native*                                 |
|                | `uint32`                       | *native*                                   | *native*                                 |
|                | `uint64`                       | *native*                                   | *native*                                 |
|                | `float32`                      | *native*                                   | *native*                                 |
|                | `float64`                      | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int8[]`       | `[]int`                        | `Int8ArrayFromIntSlice`                    | `Int8ArrayToIntSlice`                    |
|                | `[]int8`                       | `Int8ArrayFromInt8Slice`                   | `Int8ArrayToInt8Slice`                   |
|                | `[]int16`                      | `Int8ArrayFromInt16Slice`                  | `Int8ArrayToInt16Slice`                  |
|                | `[]int32`                      | `Int8ArrayFromInt32Slice`                  | `Int8ArrayToInt32Slice`                  |
|                | `[]int64`                      | `Int8ArrayFromInt64Slice`                  | `Int8ArrayToInt64Slice`                  |
|                | `[]uint`                       | `Int8ArrayFromUintSlice`                   | `Int8ArrayToUintSlice`                   |
|                | `[]uint8`                      | `Int8ArrayFromUint8Slice`                  | `Int8ArrayToUint8Slice`                  |
|                | `[]uint16`                     | `Int8ArrayFromUint16Slice`                 | `Int8ArrayToUint16Slice`                 |
|                | `[]uint32`                     | `Int8ArrayFromUint32Slice`                 | `Int8ArrayToUint32Slice`                 |
|                | `[]uint64`                     | `Int8ArrayFromUint64Slice`                 | `Int8ArrayToUint64Slice`                 |
|                | `[]float32`                    | `Int8ArrayFromFloat32Slice`                | `Int8ArrayToFloat32Slice`                |
|                | `[]float64`                    | `Int8ArrayFromFloat64Slice`                | `Int8ArrayToFloat64Slice`                |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int8range`    | `[2]int`                       | `Int8RangeFromIntArray2`                   | `Int8RangeToIntArray2`                   |
|                | `[2]int8`                      | `Int8RangeFromInt8Array2`                  | `Int8RangeToInt8Array2`                  |
|                | `[2]int16`                     | `Int8RangeFromInt16Array2`                 | `Int8RangeToInt16Array2`                 |
|                | `[2]int32`                     | `Int8RangeFromInt32Array2`                 | `Int8RangeToInt32Array2`                 |
|                | `[2]int64`                     | `Int8RangeFromInt64Array2`                 | `Int8RangeToInt64Array2`                 |
|                | `[2]uint`                      | `Int8RangeFromUintArray2`                  | `Int8RangeToUintArray2`                  |
|                | `[2]uint8`                     | `Int8RangeFromUint8Array2`                 | `Int8RangeToUint8Array2`                 |
|                | `[2]uint16`                    | `Int8RangeFromUint16Array2`                | `Int8RangeToUint16Array2`                |
|                | `[2]uint32`                    | `Int8RangeFromUint32Array2`                | `Int8RangeToUint32Array2`                |
|                | `[2]uint64`                    | `Int8RangeFromUint64Array2`                | `Int8RangeToUint64Array2`                |
|                | `[2]float32`                   | `Int8RangeFromFloat32Array2`               | `Int8RangeToFloat32Array2`               |
|                | `[2]float64`                   | `Int8RangeFromFloat64Array2`               | `Int8RangeToFloat64Array2`               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `int8range[]`  | `[][2]int`                     | `Int8RangeArrayFromIntArray2Slice`         | `Int8RangeArrayToIntArray2Slice`         |    
|                | `[][2]int8`                    | `Int8RangeArrayFromInt8Array2Slice`        | `Int8RangeArrayToInt8Array2Slice`        |
|                | `[][2]int16`                   | `Int8RangeArrayFromInt16Array2Slice`       | `Int8RangeArrayToInt16Array2Slice`       |
|                | `[][2]int32`                   | `Int8RangeArrayFromInt32Array2Slice`       | `Int8RangeArrayToInt32Array2Slice`       |
|                | `[][2]int64`                   | `Int8RangeArrayFromInt64Array2Slice`       | `Int8RangeArrayToInt64Array2Slice`       |
|                | `[][2]uint`                    | `Int8RangeArrayFromUintArray2Slice`        | `Int8RangeArrayToUintArray2Slice`        |
|                | `[][2]uint8`                   | `Int8RangeArrayFromUint8Array2Slice`       | `Int8RangeArrayToUint8Array2Slice`       |
|                | `[][2]uint16`                  | `Int8RangeArrayFromUint16Array2Slice`      | `Int8RangeArrayToUint16Array2Slice`      |
|                | `[][2]uint32`                  | `Int8RangeArrayFromUint32Array2Slice`      | `Int8RangeArrayToUint32Array2Slice`      |
|                | `[][2]uint64`                  | `Int8RangeArrayFromUint64Array2Slice`      | `Int8RangeArrayToUint64Array2Slice`      |
|                | `[][2]float32`                 | `Int8RangeArrayFromFloat32Array2Slice`     | `Int8RangeArrayToFloat32Array2Slice`     |
|                | `[][2]float64`                 | `Int8RangeArrayFromFloat64Array2Slice`     | `Int8RangeArrayToFloat64Array2Slice`     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `interval`     | **TODO**                       | **TODO**                                   | **TODO**                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `interval[]`   | **TODO**                       | **TODO**                                   | **TODO**                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `json`         | `interface{}`                  | `JSON`                                     | `JSON`                                   |
|                | `[]byte`                       | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
| `json[]`       | `[][]byte`                     | `JSONArrayFromByteSliceSlice`              | `JSONArrayToByteSliceSlice`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `jsonb`        | `interface{}`                  | `JSON`                                     | `JSON`                                   |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `jsonb[]`      | `[][]byte`                     | `JSONArrayFromByteSliceSlice`              | `JSONArrayToByteSliceSlice`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `line`         | `[3]float64`                   | `LineFromFloat64Array3`                    | `LineToFloat64Array3`                    |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `line[]`       | `[][3]float64`                 | `LineArrayFromFloat64Array3Slice`          | `LineArrayToFloat64Array3Slice`          |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `lseg`         | `[2][2]float64`                | `LsegFromFloat64Array2Array2`              | `LsegToFloat64Array2Array2`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `lseg[]`       | `[][2][2]float64`              | `LsegArrayFromFloat64Array2Array2Slice`    | `LsegArrayToFloat64Array2Array2Slice`    |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `macaddr`      | `net.HardwareAddr`             | `MACAddrFromHardwareAddr`                  | `MACAddrToHardwareAddr`                  |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `macaddr[]`    | `[]net.HardwareAddr`           | `MACAddrArrayFromHardwareAddrSlice`        | `MACAddrArrayToHardwareAddrSlice`        |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `macaddr8`     | `net.HardwareAddr`             | `MACAddr8FromHardwareAddr`                 | `MACAddr8ToHardwareAddr`                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `macaddr8[]`   | `[]net.HardwareAddr`           | `MACAddr8ArrayFromHardwareAddrSlice`       | `MACAddr8ArrayToHardwareAddrSlice`       |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `money`        | `int64`                        | `MoneyFromInt64`                           | `MoneyToInt64`                           |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `money[]`      | `[]int64`                      | `MoneyArrayFromInt64Slice`                 | `MoneyArrayToInt64Slice`                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `numeric`      | `int`                          | *native*                                   | *native*                                 |
|                | `int8`                         | *native*                                   | *native*                                 |
|                | `int16`                        | *native*                                   | *native*                                 |
|                | `int32`                        | *native*                                   | *native*                                 |
|                | `int64`                        | *native*                                   | *native*                                 |
|                | `uint`                         | *native*                                   | *native*                                 |
|                | `uint8`                        | *native*                                   | *native*                                 |
|                | `uint16`                       | *native*                                   | *native*                                 |
|                | `uint32`                       | *native*                                   | *native*                                 |
|                | `uint64`                       | *native*                                   | *native*                                 |
|                | `float32`                      | *native*                                   | *native*                                 |
|                | `float64`                      | *native*                                   | *native*                                 |
|                | `big.Int`                      | **TODO**                                   | **TODO**                                 |
|                | `big.Float`                    | **TODO**                                   | **TODO**                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `numeric[]`    | `[]int`                        | `NumericArrayFromIntSlice`                 | `NumericArrayToIntSlice`                 |
|                | `[]int8`                       | `NumericArrayFromInt8Slice`                | `NumericArrayToInt8Slice`                |
|                | `[]int16`                      | `NumericArrayFromInt16Slice`               | `NumericArrayToInt16Slice`               |
|                | `[]int32`                      | `NumericArrayFromInt32Slice`               | `NumericArrayToInt32Slice`               |
|                | `[]int64`                      | `NumericArrayFromInt64Slice`               | `NumericArrayToInt64Slice`               |
|                | `[]uint`                       | `NumericArrayFromUintSlice`                | `NumericArrayToUintSlice`                |
|                | `[]uint8`                      | `NumericArrayFromUint8Slice`               | `NumericArrayToUint8Slice`               |
|                | `[]uint16`                     | `NumericArrayFromUint16Slice`              | `NumericArrayToUint16Slice`              |
|                | `[]uint32`                     | `NumericArrayFromUint32Slice`              | `NumericArrayToUint32Slice`              |
|                | `[]uint64`                     | `NumericArrayFromUint64Slice`              | `NumericArrayToUint64Slice`              |
|                | `[]float32`                    | `NumericArrayFromFloat32Slice`             | `NumericArrayToFloat32Slice`             |
|                | `[]float64`                    | `NumericArrayFromFloat64Slice`             | `NumericArrayToFloat64Slice`             |
|                | `[]big.Int`                    | **TODO**                                   | **TODO**                                 |
|                | `[]big.Float`                  | **TODO**                                   | **TODO**                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `numrange`     | `[2]int`                       | `NumRangeFromIntArray2`                    | `NumRangeToIntArray2`                    |
|                | `[2]int8`                      | `NumRangeFromInt8Array2`                   | `NumRangeToInt8Array2`                   |
|                | `[2]int16`                     | `NumRangeFromInt16Array2`                  | `NumRangeToInt16Array2`                  |
|                | `[2]int32`                     | `NumRangeFromInt32Array2`                  | `NumRangeToInt32Array2`                  |
|                | `[2]int64`                     | `NumRangeFromInt64Array2`                  | `NumRangeToInt64Array2`                  |
|                | `[2]uint`                      | `NumRangeFromUintArray2`                   | `NumRangeToUintArray2`                   |
|                | `[2]uint8`                     | `NumRangeFromUint8Array2`                  | `NumRangeToUint8Array2`                  |
|                | `[2]uint16`                    | `NumRangeFromUint16Array2`                 | `NumRangeToUint16Array2`                 |
|                | `[2]uint32`                    | `NumRangeFromUint32Array2`                 | `NumRangeToUint32Array2`                 |
|                | `[2]uint64`                    | `NumRangeFromUint64Array2`                 | `NumRangeToUint64Array2`                 |
|                | `[2]float32`                   | `NumRangeFromFloat32Array2`                | `NumRangeToFloat32Array2`                |
|                | `[2]float64`                   | `NumRangeFromFloat64Array2`                | `NumRangeToFloat64Array2`                |
|                | `[2]big.Int`                   | **TODO**                                   | **TODO**                                 |
|                | `[2]big.Float`                 | **TODO**                                   | **TODO**                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `numrange[]`   | `[][2]int`                     | `NumRangeArrayFromIntArray2Slice`          | `NumRangeArrayToIntArray2Slice`          |
|                | `[][2]int8`                    | `NumRangeArrayFromInt8Array2Slice`         | `NumRangeArrayToInt8Array2Slice`         |
|                | `[][2]int16`                   | `NumRangeArrayFromInt16Array2Slice`        | `NumRangeArrayToInt16Array2Slice`        |
|                | `[][2]int32`                   | `NumRangeArrayFromInt32Array2Slice`        | `NumRangeArrayToInt32Array2Slice`        |
|                | `[][2]int64`                   | `NumRangeArrayFromInt64Array2Slice`        | `NumRangeArrayToInt64Array2Slice`        |
|                | `[][2]uint`                    | `NumRangeArrayFromUintArray2Slice`         | `NumRangeArrayToUintArray2Slice`         |
|                | `[][2]uint8`                   | `NumRangeArrayFromUint8Array2Slice`        | `NumRangeArrayToUint8Array2Slice`        |
|                | `[][2]uint16`                  | `NumRangeArrayFromUint16Array2Slice`       | `NumRangeArrayToUint16Array2Slice`       |
|                | `[][2]uint32`                  | `NumRangeArrayFromUint32Array2Slice`       | `NumRangeArrayToUint32Array2Slice`       |
|                | `[][2]uint64`                  | `NumRangeArrayFromUint64Array2Slice`       | `NumRangeArrayToUint64Array2Slice`       |
|                | `[][2]float32`                 | `NumRangeArrayFromFloat32Array2Slice`      | `NumRangeArrayToFloat32Array2Slice`      |
|                | `[][2]float64`                 | `NumRangeArrayFromFloat64Array2Slice`      | `NumRangeArrayToFloat64Array2Slice`      |
|                | `[][2]big.Int`                 | **TODO**                                   | **TODO**                                 |
|                | `[][2]big.Float`               | **TODO**                                   | **TODO**                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `path`         | `[][2]float64`                 | `PathFromFloat64Array2Slice`               | `PathToFloat64Array2Slice`               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `path[]`       | `[][][2]float64`               | `PathArrayFromFloat64Array2SliceSlice`     | `PathArrayToFloat64Array2SliceSlice`     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `point`        | `[2]float64`                   | `PointFromFloat64Array2`                   | `PointToFloat64Array2`                   |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `point[]`      | `[][2]float64`                 | `PointArrayFromFloat64Array2Slice`         | `PointArrayToFloat64Array2Slice`         |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `polygon`      | `[][2]float64`                 | `PolygonFromFloat64Array2Slice`            | `PolygonToFloat64Array2Slice`            |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `polygon[]`    | `[][][2]float64`               | `PolygonArrayFromFloat64Array2SliceSlice`  | `PolygonArrayToFloat64Array2SliceSlice`  |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `text`         | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `text[]`       | `[]string`                     | `TextArrayFromStringSlice`                 | `TextArrayToStringSlice`                 |
|                | `[][]byte`                     | `TextArrayFromByteSliceSlice`              | `TextArrayToByteSliceSlice`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `time`         | `time.Time`                    | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | `TimeToString`                           |
|                | `[]byte`                       | *native*                                   | `TimeToByteSlice`                        |
| `time[]`       | `[]time.Time`                  | `TimeArrayFromTimeSlice`                   | `TimeArrayToTimeSlice`                   |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `timestamp`    | `time.Time`                    | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `timestamp[]`  | `[]time.Time`                  | `TimestampArrayFromTimeSlice`              | `TimestampArrayToTimeSlice`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `timestamptz`  | `time.Time`                    | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `timestamptz[]`| `[]time.Time`                  | `TimestamptzArrayFromTimeSlice`            | `TimestamptzArrayToTimeSlice`            |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `timetz`       | `time.Time`                    | *native*                                   | *native*                                 |
|                | `string`                       | *native*                                   | `TimetzToString`                         |
|                | `[]byte`                       | *native*                                   | `TimetzToByteSlice`                      |
| `timetz[]`     | `[]time.Time`                  | `TimetzArrayFromTimeSlice`                 | `TimetzArrayToTimeSlice`                 |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tsquery`      | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tsquery[]`    | `[]string`                     | `TSQueryArrayFromStringSlice`              | `TSQueryArrayToStringSlice`              |
|                | `[][]byte`                     | `TSQueryArrayFromByteSliceSlice`           | `TSQueryArrayToByteSliceSlice`           |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tsrange`      | `[2]time.Time`                 | `TsRangeFromTimeArray2`                    | `TsRangeToTimeArray2`                    |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tsrange[]`    | `[][2]time.Time`               | `TsRangeArrayFromTimeArray2Slice`          | `TsRangeArrayToTimeArray2Slice`          |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tstzrange`    | `[2]time.Time`                 | `TstzRangeFromTimeArray2`                  | `TstzRangeToTimeArray2`                  |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tstzrange[]`  | `[][2]time.Time`               | `TstzRangeArrayFromTimeArray2Slice`        | `TstzRangeArrayToTimeArray2Slice`        |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tsvector`     | `[]string`                     | `TSVectorFromStringSlice`                  | `TSVectorToStringSlice`                  |
|                | `[][]byte`                     | `TSVectorFromByteSliceSlice`               | `TSVectorToByteSliceSlice`               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `tsvector[]`   | `[][]string`                   | `TSVectorArrayFromStringSliceSlice`        | `TSVectorArrayToStringSliceSlice`        |
|                | `[][][]byte`                   | `TSVectorArrayFromByteSliceSliceSlice`     | `TSVectorArrayToByteSliceSliceSlice`     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `uuid`         | `[16]byte`                     | `UUIDFromByteArray16`                      | `UUIDToByteArray16`                      |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `uuid[]`       | `[][16]byte`                   | `UUIDArrayFromByteArray16Slice`            | `UUIDArrayToByteArray16Slice`            |
|                | `[]string`                     | `UUIDArrayFromStringSlice`                 | `UUIDArrayToStringSlice`                 |
|                | `[][]byte`                     | `UUIDArrayFromByteSliceSlice`              | `UUIDArrayToByteSliceSlice`              |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `varbit`       | `int64`                        | `VarBitFromInt64`                          | `VarBitToInt64`                          |
|                | `[]bool`                       | `VarBitFromBoolSlice`                      | `VarBitToBoolSlice`                      |
|                | `[]uint8`                      | `VarBitFromUint8Slice`                     | `VarBitToUint8Slice`                     |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `varbit[]`     | `[][]bool`                     | `VarBitArrayFromBoolSliceSlice`            | `VarBitArrayToBoolSliceSlice`            |
|                | `[][]uint8`                    | `VarBitArrayFromUint8SliceSlice`           | `VarBitArrayToUint8SliceSlice`           |
|                | `[]string`                     | `VarBitArrayFromStringSlice`               | `VarBitArrayToStringSlice`               |
|                | `[]int64`                      | `VarBitArrayFromInt64Slice`                | `VarBitArrayToInt64Slice`                |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `varchar`      | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `varchar[]`    | `[]string`                     | `VarCharArrayFromStringSlice`              | `VarCharArrayToStringSlice`              |
|                | `[][]byte`                     | `VarCharArrayFromByteSliceSlice`           | `VarCharArrayToByteSliceSlice`           |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `xml`          | `interface{}`                  | `XML`                                      | `XML`                                    |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
| `xml[]`        | `[][]byte`                     | `XMLArrayFromByteSliceSlice`               | `XMLArrayToByteSliceSlice`               |
|                | `string`                       | *native*                                   | *native*                                 |
|                | `[]byte`                       | *native*                                   | *native*                                 |
