package models

import "time"

type Province struct {
	Id        uint
	Name      string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
