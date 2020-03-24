## wip

[![GoDoc](http://godoc.org/github.com/frk/gosql?status.png)](http://godoc.org/github.com/frk/gosql)  [![Coverage](http://gocover.io/_badge/github.com/frk/gosql?nocache=gosql)](http://gocover.io/github.com/frk/gosql)

-----------------------

```go

query := new(InsertUserQuery)
query.User = u
if err := db.ExecQuery(query); err != nil {
	return err
}
_ = query.User.Id

```

-----------------------

```go

q1 := new(SelectParentByIDQuery)
q1.Where.ID = 123
q2 := new(SelectChildrenByParentIDQuery)
q2.Where.ParentID = 123
q2.Limit = 25
if err := db.ExecQuery(q1, q2); err != nil {
	return err
}
_ = q1.Parent
_ = q1.Children

```
