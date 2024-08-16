package web

type ArticleVo struct {
	Id         int64  `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Abstract   string `json:"abstract,omitempty"`
	Content    string `json:"content,omitempty"`
	AuthorId   int64  `json:"authorId,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	Status     uint8  `json:"status,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`

	ReadCnt    int64 `json:"readCnt"`
	LikeCnt    int64 `json:"likeCnt"`
	CollectCnt int64 `json:"collectCnt"`
	Liked      bool  `json:"liked"`
	Collected  bool  `json:"collected"`
}

type PublishReq struct {
	Id      int64
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleEditReq struct {
	Id      int64
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleWithdrawReq struct {
	Id int64
}

type ArticleLikeReq struct {
	Id int64 `json:"id"`
	// true 是点赞，false 是不点赞
	Like bool `json:"like"`
}

type ArticleCollectReq struct {
	Id  int64 `json:"id"`
	Cid int64 `json:"cid"`
}
