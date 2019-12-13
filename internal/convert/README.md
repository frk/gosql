## convert

### TODO: 

**basic types**
- [ ] `bit` format: '1'
- [ ] `boolean` format: true
- [ ] `bpchar` format: 'a'
- [ ] `bpchar(3)` format: 'abc'
- [ ] `char` format: 'a'
- [ ] `char(3)` format: 'abc'
- [ ] `cidr` format: '192.168.100.128/25'
- [ ] `date` format: '1999-01-08'
- [ ] `float4` format: 0.1
- [ ] `float8` format: 0.1
- [ ] `inet` format: '192.168.0.1/24'
- [ ] `int2` format: 1
- [ ] `int4` format: 1
- [ ] `int8` format: 1
- [ ] `interval` format: '1 day'
- [ ] `macaddr` format: '08:00:2b:01:02:03'
- [ ] `macaddr8` format: '08:00:2b:01:02:03:04:05'
- [ ] `money` format: '$1.20'
- [ ] `numeric` format: 1.2
- [ ] `text` format: 'foo'
- [ ] `time` format: '04:05:06.789'
- [ ] `timestamp` format: '1999-01-08 04:05:06'
- [ ] `timestamptz` format: 'January 8 04:05:06 1999 PST'
- [ ] `timetz` format: '04:05:06.789-8'
- [ ] `tsquery` format: 'fat & rat'
- [ ] `tsvector` format: 'a fat cat sat on a mat and ate a fat rat'
- [ ] `uuid` format: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
- [ ] `varbit` format: '101'
- [ ] `varbit(1)` format: '1'
- [ ] `varchar` format: 'foo'
- [ ] `varchar(1)` format: 'a'

**aggregate types**
- [x] `bit[]` expected format: `{b[,...]}` where `b` is either `0` or `1`.

| Go Type   | Scanner             | Valuer
| ----------|---------------------|--------
| `[]bool`  | `BitArr2BoolSlice`  |
| `[]uint8` | `BitArr2Uint8Slice` |
| `[]uint`  | `BitArr2UintSlice`  |

- [x] `boolean[]` expected format: `{t[,...]}` where `t` is either `t` or `f`.

| Go Type   | Scanner             | Valuer
| ----------|---------------------|--------
| `[]bool`  | `BoolArr2BoolSlice` |

- [x] `box[]` format: `{(x1,y1),(x2,y2)[;...]}`where `(x1,y1)` and `(x2,y2)` are two opposite corners of a box.

| Go Type           | Scanner                   | Valuer
| ------------------|---------------------------|--------
| `[][2][2]float64` | `BoxArr2Float64a2a2Slice` |

- [ ] `bpchar[]` format: ARRAY['a','b']::bpchar[]
- [ ] `bpchar(3)[]` format: ARRAY['abc','def']::bpchar(3)[]
- [ ] `bytea` format: '\xDEADBEEF'
- [ ] `bytea[]` format: ARRAY['\xDEADBEEF', '\xDEADBEEF']::bytea[]
- [ ] `char[]` format: ARRAY['a','b']::char[]
- [ ] `char(3)[]` format: ARRAY['abc','def']::char(3)[]
- [ ] `cidr[]` format: ARRAY['192.168.100.128/25','2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128']::cidr[]
- [ ] `circle[]` format: ARRAY['<(0,0), 3.5>','((0.5,1), 5)']::circle[]
- [ ] `date[]` format: ARRAY['1999-01-08', 'May 5, 2001']::date[]
- [ ] `daterange[]` format: ARRAY['[1999-01-08, 2001-01-08)', '(1999-01-08, 2001-01-08]']::daterange[]
- [ ] `float4[]` format: `{f[,...]}` where `f` is floating point number.

| Go Type     | Scanner                 | Valuer
| ------------|-------------------------|--------
| `[]float32` | `FloatArr2Float32Slice` |
| `[]float64` | `FloatArr2Float64Slice` |
| `[]int8`    | `FloatArr2Int8Slice`    |
| `[]int16`   | `FloatArr2Int16Slice`   |
| `[]int32`   | `FloatArr2Int32Slice`   |
| `[]int64`   | `FloatArr2Int64Slice`   |
| `[]int`     | `FloatArr2IntSlice`     |
| `[]uint8`   | `FloatArr2Uint8Slice`   |
| `[]uint16`  | `FloatArr2Uint16Slice`  |
| `[]uint32`  | `FloatArr2Uint32Slice`  |
| `[]uint64`  | `FloatArr2Uint64Slice`  |
| `[]uint`    | `FloatArr2UintSlice`    |

- [ ] `float8[]` format: `{f[,...]}` where `f` is floating point number.

| Go Type     | Scanner                 | Valuer
| ------------|-------------------------|--------
| `[]float32` | `FloatArr2Float32Slice` |
| `[]float64` | `FloatArr2Float64Slice` |
| `[]int8`    | `FloatArr2Int8Slice`    |
| `[]int16`   | `FloatArr2Int16Slice`   |
| `[]int32`   | `FloatArr2Int32Slice`   |
| `[]int64`   | `FloatArr2Int64Slice`   |
| `[]int`     | `FloatArr2IntSlice`     |
| `[]uint8`   | `FloatArr2Uint8Slice`   |
| `[]uint16`  | `FloatArr2Uint16Slice`  |
| `[]uint32`  | `FloatArr2Uint32Slice`  |
| `[]uint64`  | `FloatArr2Uint64Slice`  |
| `[]uint`    | `FloatArr2UintSlice`    |

- [ ] `hstore[]` format: ARRAY['a=>1,b=>2', 'c=>3,d=>4']::hstore[]
- [ ] `inet[]` format: ARRAY['192.168.0.1/24', '128.0.0.0/16']::inet[]
- [ ] `int2[]` format: ARRAY[1,2]::int2[]
- [ ] `int2vector` format: '1 2'
- [ ] `int2vector[]` format: ARRAY['1 2', '3 4']::int2vector[]
- [ ] `int4[]` format: ARRAY[1,2]::int4[]
- [ ] `int4range[]` format: ARRAY['[0,9)', '(0,9]']::int4range[]
- [ ] `int8[]` format: ARRAY[1,2]::int8[]
- [ ] `int8range[]` format: ARRAY['[0,9)', '(0,9]']::int8range[]
- [ ] `interval[]` format: ARRAY['1 day','5 years 4 months 34 minutes ago']::interval[]
- [ ] `json[]` format: ARRAY['{"foo":["bar", "baz", 123]}', '["foo", 123]']::json[]
- [ ] `jsonb[]` format: ARRAY['{"foo":["bar", "baz", 123]}', '["foo", 123]']::jsonb[]
- [ ] `line[]` format: ARRAY['{1,2,3}', '{4,5,6}']::line[]
- [ ] `lseg[]` format: ARRAY['[(1,2), (3,4)]', '[(1.5,2.5), (3.5,4.5)]']::lseg[]
- [ ] `macaddr[]` format: ARRAY['08:00:2b:01:02:03', '08002b010203']::macaddr[]
- [ ] `macaddr8[]` format: ARRAY['08:00:2b:01:02:03:04:05', '08002b0102030405']::macaddr8[]
- [ ] `money[]` format: ARRAY['$1.20', '$0.99']::money[]
- [ ] `numeric[]` format: ARRAY[1.2,3.4]::numeric[]
- [ ] `numrange[]` format: ARRAY['[1.2,3.4)', '(1.2,3.4]']::numrange[]
- [ ] `path[]` format: ARRAY['[(1,1),(2,2),(3,3)]', '[(1.5,1.5),(2.5,2.5),(3.5,3.5)]']::path[]
- [ ] `point[]` format: ARRAY['(1,1)', '(2,2)']::point[]
- [ ] `polygon[]` format: ARRAY['((1,1),(2,2),(3,3))', '((1.5,1.5),(2.5,2.5),(3.5,3.5))']::polygon[]
- [ ] `text[]` format: ARRAY['foo', 'bar']::text[]
- [ ] `time[]` format: ARRAY['04:05:06.789', '040506']::time[]
- [ ] `timestamp[]` format: ARRAY['1999-01-08 04:05:06', '2004-10-19 10:23:54']::timestamp[]
- [ ] `timestamptz[]` format: ARRAY['January 8 04:05:06 1999 PST','2004-10-19 10:23:54+02']::timestamptz[]
- [ ] `timetz[]` format: ARRAY['04:05:06.789-8','2003-04-12 04:05:06 America/New_York']::timetz[]
- [ ] `tsquery[]` format: ARRAY['fat & rat', 'fat & rat & ! cat']::tsquery[]
- [ ] `tsrange[]` format: ARRAY['[1999-01-08 04:05:06, 2004-10-19 10:23:54)', '(1999-01-08 04:05:06, 2004-10-19 10:23:54]']::tsrange[]
- [ ] `tstzrange[]` format: ARRAY['[January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02)','(January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02]']::tstzrange[]
- [ ] `tsvector[]` format: ARRAY['a fat cat sat on a mat','and ate a fat rat']::tsvector[]
- [ ] `uuid[]` format: ARRAY['a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11','a0eebc999c0b4ef8bb6d6bb9bd380a11']::uuid[]
- [ ] `varbit[]` format: '{101, 00}'::varbit[]
- [ ] `varbit(1)[]` format: '{1, 0}'::varbit(1)[]
- [ ] `varchar[]` format: ARRAY['foo', 'bar']::varchar[]
- [ ] `varchar(1)[]` format: ARRAY['a', 'b']::varchar(1)[]
- [ ] `xml[]` format: ARRAY['<foo>bar</foo>','<bar>foo</bar>']::xml[]

**composite types**
- [x] `box` format: `(x1,y1),(x2,y2)` where `(x1,y1)` and `(x2,y2)` are two opposite corners of the box.

| Go Type         | Scanner           | Valuer
| ----------------|-------------------|--------
| `[2][2]float64` | `Box2Float64a2a2` |

- [ ] `circle` format: `<(0,0), 3.5>`
- [ ] `hstore` format: `a=>1,b=>2`
- [ ] `line` format: `{1,2,3}`
- [ ] `lseg`  format:`[(1,2), (3,4)]`
- [ ] `path` format: `[(1,1),(2,2),(3,3)]`
- [ ] `point` format: `(1,1)`
- [ ] `polygon` format: `((1,1),(2,2),(3,3))`

**range types**
- [ ] `daterange` format: '[1999-01-08, 2001-01-08)'
- [ ] `int4range` format: '[0,9)'
- [ ] `int8range` format: '[0,9)'
- [ ] `numrange` format: '[1.2,3.4)'
- [ ] `tsrange` format: '[1999-01-08 04:05:06, 2004-10-19 10:23:54)'
- [ ] `tstzrange` format: '[January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02)'

**"format" types**
- [ ] `json` format: '{"foo":["bar", "baz", 123]}'
- [ ] `jsonb` format: '{"foo":["bar", "baz", 123]}'
- [ ] `xml` format: '<foo>bar</foo>'
