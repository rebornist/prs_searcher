package models

type Search struct {
	Code        string `bson:"code"`
	Keyword     string `bson:"keyword"`
	PrivateInfo string `bson:"private_info"`
	Content     string `bson:"content"`
	Remarks     string `bson:"remarks"`
}
