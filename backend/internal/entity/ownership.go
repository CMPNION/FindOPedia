package entity

import "time"

type Ownership struct {
	ID        int64
	ArticleID int64
	UserID    int64
	Username  string
	ClaimedAt time.Time
}
