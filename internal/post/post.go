package post

type PostFromDB struct {
    ID          int64
	Author      string
	Description string
	TimeStamp   string
}

type UserPost struct {
    ID          int64
	Circle      string
	Description string
	TimeStamp   string
}
