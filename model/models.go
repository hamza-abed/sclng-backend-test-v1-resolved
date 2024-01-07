package model

import "time"

type Licence struct {
	Key  string
	Name string
}

type Repository struct {
	ID        int64
	FullName  string
	Name      string
	CreatedAt time.Time
	Owner     *Owner
	Licence   *Licence
	Languages []*Language
}

/*
*
Many to many RepoLang
*/
type Language struct {
	Name  string
	Bytes int64
}

type Owner struct {
	ID   int64
	Name string
}
