package models

type CustomerActivity struct {
	UserID    uint   `gorm:"primaryKey"`
	CreatedAt int64  `gorm:"primaryKey"`
	Action    string `gorm:"type:varchar(20);index"`
	Data      string `gorm:"type:text"`
}

var (
	CustomAction_ViewProduct   = "VIEW_PRODUCT"
	CustomAction_SearchProduct = "SEARCH_PRODUCT"
)
