package service

import (
	"database/sql"
)

type ChartService struct {
	db *sql.DB
}
