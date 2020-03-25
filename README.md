# WIP

[![GoDoc](http://godoc.org/github.com/frk/gosql?status.png)](http://godoc.org/github.com/frk/gosql)  [![Coverage](http://gocover.io/_badge/github.com/frk/gosql?nocache=gosql)](http://gocover.io/github.com/frk/gosql)


Declare query types

```go

type InsertUserQuery struct {
	User *mypackage.User `rel:"user_table"`
	_ gosql.Return `sql:"id"`
}

type SelectParentByIDQuery struct {
	Parent *myapp.Parent `rel:"parent_table:p"`
	Where struct {
		Id int `sql:"p.id"`
	}
}

type SelectChildrenByParentIDQuery struct {
	Children []*myapp.Child `rel:"children_table:c"`
	Where struct {
		ParentId int `sql:"c.parent_id"`
	}
	Limit int
}

```

-----------------------

Run generator and then use the queries like so:

```go

query := new(InsertUserQuery)
query.User = u
if err := db.ExecQuery(query); err != nil {
	return err
}
_ = query.User.Id

// ....

q1 := new(SelectParentByIDQuery)
q1.Where.ID = 123
q2 := new(SelectChildrenByParentIDQuery)
q2.Where.ParentID = 123
q2.Limit = 25
if err := db.ExecQuery(q1, q2); err != nil {
	return err
}
_ = q1.Parent
_ = q2.Children

```
