package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	// 不做组合 做字段
	Author Author
}

// Author 在文章领域是一个值对象
type Author struct {
	Id   int64
	Name string
}
