package domain

type Article struct {
	Title   string
	Content string
	// 不做组合 做字段
	Author Author
}

type Author struct {
	Id   int64
	Name string
}
