package model

import (
	"github.com/jinzhu/gorm"
	"pdk/src/server/lib/db"
	"testing"
)

func init() {
	db.Init("pdk:WKwcyf66fTFKtip4@tcp(192.168.176.128:3306)/pdk?charset=utf8&parseTime=True&loc=Local")
}

type Email struct {
	Email string
	MyUserID uint `gorm:"primary_key"`
}

type Language struct {
	Name string
	MyUserID uint `gorm:"primary_key"`
}

type Address struct {
	Address1 string
	MyUserID uint `gorm:"primary_key"`
}

type MyUser struct {
	gorm.Model
	Name string
	BillingAddress Address
	ShippingAddress Address
	Emails []Email
	Languages []Language
}

func TestMyGORM_Create(t *testing.T) {
	user := MyUser{
		Name:            "jinzhu",
		BillingAddress:  Address{Address1: "Billing Address - Address 2222"},
		ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
		Emails:          []Email{
			{Email: "jinzhu@example.com"},
			{Email: "jinzhu-2@example.com"},
		},
		Languages:       []Language{
			{Name: "ZH111"},
			{Name: "EN"},
		},
	}
	//db.GetGormDB().AutoMigrate(&user,&Address{},&Email{},&Language{})
	//db.GetGormDB().Create(&user)
	db.GetGormDB().Save(&user)
}


