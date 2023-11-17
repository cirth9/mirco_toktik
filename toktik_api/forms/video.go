package forms

type Video struct {
	Id            int64    `json:"id,omitempty"`
	UserInfoId    int64    `json:"-"`
	Author        UserInfo `json:"author,omitempty" gorm:"-"`
	PlayUrl       string   `json:"play_url,omitempty"`
	CoverUrl      string   `json:"cover_url,omitempty"`
	FavoriteCount int64    `json:"favorite_count,omitempty"`
	CommentCount  int64    `json:"comment_count,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty"`
	Title         string   `json:"title,omitempty"`
}

type VideoData struct {
	Content  []byte `json:"content"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}
