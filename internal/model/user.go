package model

type User struct {
	Id   uint64 `db:"id" json:"Id"`
	Name string `db:"name" json:"Name"`
	Age  uint32 `db:"age" json:"Age"`
}
