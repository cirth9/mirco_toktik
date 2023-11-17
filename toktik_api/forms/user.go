package forms

type UserInfo struct {
	Id              int64  `json:"id" gorm:"id,omitempty"`
	Name            string `json:"name" gorm:"name,omitempty"`
	FollowCount     int64  `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow" gorm:"is_follow,omitempty"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}
