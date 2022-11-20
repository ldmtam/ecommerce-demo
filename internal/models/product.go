package models

type Product struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(100);index:,class:FULLTEXT,option:WITH PARSER ngram" json:"name"`
	Price     uint   `json:"price"`
	CreatedAt int64  `json:"createdAt"`
}
