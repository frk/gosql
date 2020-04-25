package testdata

import (
	"net"
	"time"
)

type type1 struct {
	bit      bool             `sql:"col_bit"`
	bpchar   byte             `sql:"col_bpchar"`
	char     rune             `sql:"col_char"`
	cidr     net.IPNet        `sql:"col_cidr"`
	inet     net.IP           `sql:"col_inet"`
	macaddr  net.HardwareAddr `sql:"col_macaddr"`
	macaddr8 net.HardwareAddr `sql:"col_macaddr8"`
	money    int64            `sql:"col_money"`
	tsvector []string         `sql:"col_tsvector"`
	uuid     [16]byte         `sql:"col_uuid"`
	varbit   []bool           `sql:"col_varbit"`
}

type type2 struct {
	bitarr    []uint8         `sql:"col_bitarr"`
	boolarr   []bool          `sql:"col_boolarr"`
	boxarr    [][2][2]float64 `sql:"col_boxarr"`
	bpchararr []rune          `sql:"col_bpchararr"`
	//bpchar3arr []string        `sql:"col_bpchar3arr"`
	byteaarr [][]byte `sql:"col_byteaarr"`
	chararr  []byte   `sql:"col_chararr"`
	//char3arr []string    `sql:"col_char3arr"`
	cidrarr []net.IPNet `sql:"col_cidrarr"`
	//circlearr      string          `sql:"col_circlearr"`
	datearr       []time.Time    `sql:"col_datearr"`
	daterangearr  [][2]time.Time `sql:"col_daterangearr"`
	float4arr     []float32      `sql:"col_float4arr"`
	float8arr     []float64      `sql:"col_float8arr"`
	inetarr       []net.IP       `sql:"col_inetarr"`
	int2arr       []int16        `sql:"col_int2arr"`
	int2vector    []int16        `sql:"col_int2vector"`
	int2vectorarr [][]int16      `sql:"col_int2vectorarr"`
	int4arr       []int32        `sql:"col_int4arr"`
	int4rangearr  [][2]int32     `sql:"col_int4rangearr"`
	int8arr       []int64        `sql:"col_int8arr"`
	int8rangearr  [][2]int64     `sql:"col_int8rangearr"`
	//intervalarr    string       `sql:"col_intervalarr"`
	jsonarr        [][]byte           `sql:"col_jsonarr"`
	jsonbarr       [][]byte           `sql:"col_jsonbarr"`
	linearr        [][3]float64       `sql:"col_linearr"`
	lsegarr        [][2][2]float64    `sql:"col_lsegarr"`
	macaddrarr     []net.HardwareAddr `sql:"col_macaddrarr"`
	macaddr8arr    []net.HardwareAddr `sql:"col_macaddr8arr"`
	moneyarr       []int64            `sql:"col_moneyarr"`
	numericarr     []int64            `sql:"col_numericarr"`
	numrangearr    [][2]float64       `sql:"col_numrangearr"`
	patharr        [][][2]float64     `sql:"col_patharr"`
	pointarr       [][2]float64       `sql:"col_pointarr"`
	polygonarr     [][][2]float64     `sql:"col_polygonarr"`
	textarr        []string           `sql:"col_textarr"`
	timearr        []time.Time        `sql:"col_timearr"`
	timestamparr   []time.Time        `sql:"col_timestamparr"`
	timestamptzarr []time.Time        `sql:"col_timestamptzarr"`
	timetzarr      []time.Time        `sql:"col_timetzarr"`
	tsqueryarr     []string           `sql:"col_tsqueryarr"`
	tsrangearr     [][2]time.Time     `sql:"col_tsrangearr"`
	tstzrangearr   [][2]time.Time     `sql:"col_tstzrangearr"`
	tsvectorarr    [][]string         `sql:"col_tsvectorarr"`
	uuidarr        [][16]byte         `sql:"col_uuidarr"`
	varbitarr      [][]bool           `sql:"col_varbitarr"`
	//varbit1arr     string   `sql:"col_varbit1arr"`
	varchararr []string `sql:"col_varchararr"`
	//varchar1arr    string   `sql:"col_varchar1arr"`
	xmlarr [][]byte `sql:"col_xmlarr"`
}
