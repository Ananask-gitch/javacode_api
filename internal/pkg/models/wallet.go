package models

type Wallet struct {
	UUID   uint64 `gorm:"primaryKey"`
	Amount uint64 `gorm:"not null"`
}
