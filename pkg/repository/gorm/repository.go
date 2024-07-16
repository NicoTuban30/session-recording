package gorm

import (
	"cassette/pkg/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ repository.Repository = &Repository{}

type Repository struct {
	db *gorm.DB
}

func NewPostgres(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return New(db)
}

func New(db *gorm.DB) (*Repository, error) {
	if err := db.AutoMigrate(&Session{}); err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Sessions() ([]repository.Session, error) {
	var sessions []Session

	if tx := r.db.Find(&sessions); tx.Error != nil {
		return nil, tx.Error
	}

	return convertSessions(sessions), nil
}

func (r *Repository) Session(id string) (*repository.Session, error) {
	var session Session

	if tx := r.db.First(&session, "id = ?", id); tx.Error != nil {
		return nil, tx.Error
	}

	return convertSession(session), nil
}

func (r *Repository) CreateSession(info *repository.SessionInfo) (*repository.Session, error) {
	if info == nil {
		info = new(repository.SessionInfo)
	}

	session := Session{
		Origin:      info.Origin,
		UserAgent:   info.UserAgent,
		UserEmail:   info.UserEmail,
		QaId:        info.QaId,
		QaSessionId: info.QaSessionId,
	}

	if tx := r.db.Create(&session); tx.Error != nil {
		return nil, tx.Error
	}

	return convertSession(session), nil
}

func (r *Repository) DeleteSession(id string) error {
	if tx := r.db.Delete(&Session{}, "id = ?", id); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *Repository) FindSessionsByQaSessionId(qaSessionId string) ([]repository.Session, error) {
	var sessions []Session

	if tx := r.db.Where("qa_session_id = ?", qaSessionId).Find(&sessions); tx.Error != nil {
		return nil, tx.Error
	}

	return convertSessions(sessions), nil
}

func convertSessions(sessions []Session) []repository.Session {
	result := make([]repository.Session, 0)

	for _, s := range sessions {
		session := convertSession(s)
		result = append(result, *session)
	}

	return result
}

func convertSession(session Session) *repository.Session {
	return &repository.Session{
		ID:          session.ID,
		Created:     session.CreatedAt,
		Updated:     session.UpdatedAt,
		Origin:      session.Origin,
		Address:     session.Address,
		UserAgent:   session.UserAgent,
		UserEmail:   session.UserEmail,
		QaId:        session.QaId,
		QaSessionId: session.QaSessionId,
	}
}
