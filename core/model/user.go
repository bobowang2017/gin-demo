package model

import "gin-demo/infra/model"

type User struct {
	model.BaseModel
	UserId    string `gorm:"default:null" json:"userId"`
	Username  string `gorm:"default:null" json:"username"`
	Avatar    string `gorm:"default:null" json:"avatar"`
	Sex       int    `gorm:"default:null" json:"sex"`
	Birthday  string `gorm:"default:null" json:"birthday"`
	Phone     string `gorm:"default:null" json:"phone"`
	Password  string `gorm:"default:null" json:"password"`
	Email     string `gorm:"default:null" json:"email"`
	Salt      string `gorm:"default:null" json:"salt"`
	Hometown  string `gorm:"default:null" json:"hometown"`
	Career    string `gorm:"default:null" json:"career"`
	ImId      string `gorm:"default:null" json:"imId"`
	VipLevel  *int   `gorm:"not null" json:"vipLevel"`
	Label     string `gorm:"default:null" json:"label"`
	Signature string `gorm:"default:null" json:"signature"`
	Status    *int   `gorm:"not null" json:"status"`
	Admin     *int   `gorm:"not null" json:"admin"`
}

type UserApiKey struct {
	model.BaseModel
	ApiKey   string `gorm:"default:null" json:"apiKey"`
	Platform string `gorm:"default:null" json:"platform"`
	UserId   string `gorm:"default:null" json:"userId"`
}

type UserAlbum struct {
	model.BaseModel
	UserId   string `gorm:"default:null" json:"userId"`
	AlbumSrc string `gorm:"default:null" json:"albumSrc"`
}

type UserLatestLogin struct {
	model.BaseModel
	UserId string `gorm:"default:null" json:"userId"`
	Ip     string `gorm:"default:null" json:"ip"`
}

type UserSignup struct {
	model.BaseModel
	UserId string `gorm:"default:null" json:"userId"`
	Year   int    `gorm:"default:null" json:"year"`
	Month  int    `gorm:"default:null" json:"month"`
	Day    int    `gorm:"default:null" json:"day"`
}

type UserVisit struct {
	model.BaseModel
	UserId      string `gorm:"default:null" json:"userId"`
	VisitUserId string `gorm:"default:null" json:"visitUserId"`
}

type VipPoint struct {
	model.BaseModel
	UserId    string `gorm:"default:null" json:"userId"`
	Value     int    `gorm:"not null" json:"value"`
	PointType string `gorm:"default:null" json:"pointType"`
	Status    *int   `gorm:"not null" json:"status"`
}

type VipPointRecord struct {
	model.BaseModel
	UserId string `gorm:"default:null" json:"userId"`
	Type   string `gorm:"default:null" json:"type"`
	Amount int    `gorm:"not null" json:"amount"`
}
