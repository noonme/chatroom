package model

type User struct {
	//为了序列化和反序列化成功，我们必须保证
	//用户信息的json字符串的key和结构体的字段对于的tag名字必须一致
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
