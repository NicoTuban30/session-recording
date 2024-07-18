package gorm

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID string `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Origin  string
	Address string

	UserAgent  string
	UserEmail  string
	QaId       string
	QaSessionId string
	AgoraStreamUrl string
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	return
}
