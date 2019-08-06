package gosql

type directive struct {
	// This field serves as an indicator that the type is actually a directive
	// as opposed to other types declared by this package. Used by helper
	// functions defined in the internal/ packages.
	_isdir struct{}
}

type (
	// The Column directive has two potential use cases:
	//
	// (1) It can be used in a Where block to produce column specific comparisons
	// for a WHERE clause condition. The type of comparison that can be produced
	// depends on the `sql` tag value supplied to the directive.
	// The expected format of the tag's value is:
	// { column [ comparison-operator [ scalar-operator ] { column | literal } ] }
	//
	// (2) It can be used in an OnConflict block to specify the resulting
	// ON CONFLICT clause's conflict_target as a list of index_column_names.
	// The list should be should be provided in the directive's `sql` tag.
	Column directive

	// This is inteded to be used with Update and Delete commands to
	// explicitly indicate that, if a "where block" is missing, the
	// command should be executed against all the rows of the relation.
	//
	// This acts as a safeguard against unintentionally omitting the
	// "where block" and then generating a query that would delete/update
	// every single rows in a table.
	All directive

	// The Relation directive has two potential use cases:
	//
	// (1) It can be used in a Delete command as the mount for the `rel` tag.
	// This can be useful for Delete commands that have no Return directive,
	// since such commands produce a DELETE query that takes no input other
	// than the optional WHERE clause parameters, nor does it generate any
	// output, and it therefore then becomes unnecessary to provide
	// a proper Go struct representation of the target relation.
	//
	// (2) It can be used as the first directive in a Using or From block
	// to specify the primary target relation for the two clauses produced
	// from those two blocks.
	Relation directive

	LeftJoin  directive
	RightJoin directive
	FullJoin  directive
	CrossJoin directive

	// The Return directive produces a postgres RETURNING clause. The columns
	// to be returned have to be specified in the struct field's tag.
	//
	// The Return directive can be used in Insert, Update, and Delete commands.
	Return directive

	// The Force directive allows for specifying columns that are usually
	// omitted from a query to actually be included in the query by the
	// command in which the directive is used. The columns to be included
	// have to be specified in the struct field's tag.
	//
	// For example if a table has an "id" column whose value is auto-generated
	// by the database, the corresponding Go struct field's column can be
	// flagged as "auto" which will exclude the column from INSERT/UPDATE
	// queries. However, there may be scenarios in an app where the "id"
	// value is already pre-generated and it needs to be INSERTed together
	// with the record, this is where the Force directive can be used to
	// tell the command to include the "id" column in the query.
	Force directive

	// The Default directive produces the DEFAULT marker in place of values
	// of those columns that are listed in the directive's field tag. This can
	// be used to specify those columns of an INSERT/UPDATE command that should
	// have their value set to their default as defined by the database table.
	Default directive

	// The Limit directive can be used inside a Select command to produce
	// a LIMIT clause for the SELECT query. The limit value must be specified
	// in the directive field's `sql` tag.
	Limit directive

	// The Offset directive can be used inside a Select command to produce
	// an OFFSET clause for the SELECT query. The offset value must be specified
	// in the directive field's `sql` tag.
	Offset directive

	// The OrderBy directive can be used inside a Select command to produce
	// an ORDER BY clause for the SELECT query. The list of columns by which
	// to order should be specified in the directive's tag.
	//
	// The expected format for each item in the directive's tag is:
	// [ - ][ qualifier. ]column_name[ :nullsfirst | :nullslast ]
	//
	// The optional preceding "-" produces the DESC sort direction option, if
	// no "-" is provided ASC is produced instead.
	// The optional :nullsfirst and :nullslast produce the NULLS FIRST and
	// NULLS LAST options respectively.
	OrderBy directive

	// The Override directive can be used in an Insert command to produce
	// the OVERRIDING { SYSTEM | USER } VALUE clause.
	Override directive

	// The TextSearch directive can be used in a Filter command to specify
	// the document column that will be used for full-text search.
	TextSearch directive

	// The Index directive can be used in an OnConflict block to specify
	// the resulting ON CONFLICT clause's conflict_target using the name
	// of a unique index. The index name should be should be provided in
	// the directive's `sql` tag.
	Index directive

	// The Constraint directive can be used in an OnConflict block to specify
	// the resulting ON CONFLICT clause's conflict_target using the name of
	// a table constraint. The constraint's name should be should be provided
	// in the directive's `sql` tag.
	Constraint directive

	// The Ignore directive can be used in an OnConflict block to produce
	// the DO NOTHING action of the resulting ON CONFLICT clause.
	Ignore directive

	// The Update directive can be used in an OnConflict block to produce
	// the DO UPDATE SET action of the resulting ON CONFLICT clause. The
	// columns to be updated by the produced action should be listed in
	// the directive's `sql` tag.
	Update directive
)

type AfterScanner interface {
	AfterScan()
}

type ErrorHandler interface {
	HandleError(err error) error
}
