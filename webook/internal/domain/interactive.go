package domain

// Interactive 这个是总体交互的计数
type Interactive struct {
	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`
	// 这个是当下这个资源，你有没有点赞或者收集
	// 你也可以考虑把这两个字段分离出去，作为一个单独的结构体
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}
