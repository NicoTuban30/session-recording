package repository

import (
	"time"
)

type Session struct {
	ID string `json:"id"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Origin  string `json:"origin"`
	Address string `json:"address"`

	UserAgent      string `json:"userAgent"`
	UserEmail      string `json:"userEmail"`
	QaId           string `json:"qaId"`
	QaSessionId    string `json:"qaSessionId"`
	AgoraStreamUrl string `json:"agoraStreamUrl"`
}

type Repository interface {
	Sessions() ([]Session, error)
	Session(id string) (*Session, error)

	CreateSession(info *SessionInfo) (*Session, error)
	UpdateSessionAgoraStreamURL(id string, agoraStreamUrl string) error
	DeleteSession(id string) error

	FindSessionsByQaSessionId(qaSessionId string) ([]Session, error)
}



type SessionInfo struct {
	Origin  string
	Address string

	UserAgent  string
	UserEmail  string
	QaId       string
	QaSessionId string
	AgoraStreamUrl string
}
