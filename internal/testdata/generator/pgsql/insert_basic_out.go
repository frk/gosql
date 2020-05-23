// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/pgsql"
)

func (q *insertbasicquery) Exec(c gosql.Conn) error {
	const queryString = `INSERT INTO "pgsql_test" (
		"col_bit"
		, "col_bpchar"
		, "col_char"
		, "col_cidr"
		, "col_inet"
		, "col_macaddr"
		, "col_macaddr8"
		, "col_money"
		, "col_tsvector"
		, "col_uuid"
		, "col_varbit"
	) VALUES (
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
	)` // `

	_, err := c.Exec(queryString,
		pgsql.BitFromBool(q.data.bit),
		pgsql.BPCharFromByte(q.data.bpchar),
		pgsql.BPCharFromRune(q.data.char),
		pgsql.CIDRFromIPNet(q.data.cidr),
		pgsql.InetFromIP(q.data.inet),
		pgsql.MACAddrFromHardwareAddr(q.data.macaddr),
		pgsql.MACAddr8FromHardwareAddr(q.data.macaddr8),
		pgsql.MoneyFromInt64(q.data.money),
		pgsql.TSVectorFromStringSlice(q.data.tsvector),
		pgsql.UUIDFromByteArray16(q.data.uuid),
		pgsql.VarBitFromBoolSlice(q.data.varbit),
	)
	return err
}
