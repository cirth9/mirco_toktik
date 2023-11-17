package models

/*
message Message {
  required int64 id = 1; // 消息id
  required int64 to_user_id = 2; // 该消息接收者的id
  required int64 from_user_id =3; // 该消息发送者的id
  required string content = 4; // 消息内容
  optional string create_time = 5; // 消息创建时间
}
*/

type Message struct {
	Id         int64  `json:"id" gorm:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id" gorm:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id" gorm:"from_user_id,omitempty"`
	Content    string `json:"content" gorm:"content,omitempty"`
	CreateTime string `json:"create_time" gorm:"create_time,omitempty"`
}

func (m *Message) TableName() string {
	return "message"
}
