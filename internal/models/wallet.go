package models

type Wallet struct {
	UUID   uint64 `gorm:"primaryKey" json:"id"`
	Amount uint64 `gorm:"not null" json:"Amount"`
}
