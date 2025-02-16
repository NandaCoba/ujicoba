package models

import "time"

type Karyawan struct {
	Id         int       `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Position   string    `json:"position"`
	Department string    `json:"departement"`
	Salary     int       `json:"salary"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
