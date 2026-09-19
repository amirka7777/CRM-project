package models

import "time"

type UserRequest struct {
	Email      string
	Password   string
	Created_at time.Time
}
