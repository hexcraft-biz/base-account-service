package models

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/hexcraft-biz/model"
	"github.com/jmoiron/sqlx"
)

const (
	PW_SALT_BYTES = 16
)

//================================================================
// Data Struct
//================================================================
type EntityUser struct {
	*model.Prototype `dive:""`
	Identity         string `db:"identity"`
	Password         string `db:"password"`
	Salt             string `db:"salt"`
	Status           string `db:"status"`
}

func (u *EntityUser) GetAbsUser() (*AbsUser, error) {
	return &AbsUser{
		ID:        *u.ID,
		Identity:  u.Identity,
		Password:  u.Password,
		Salt:      u.Salt,
		Status:    u.Status,
		CreatedAt: u.Ctime.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.Mtime.Format("2006-01-02 15:04:05"),
	}, nil
}

type AbsUser struct {
	ID        uuid.UUID `json:"id"`
	Identity  string    `json:"identity"`
	Password  string    `json:"_"`
	Salt      string    `json:"_"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

//================================================================
// Engine
//================================================================
type UsersTableEngine struct {
	*model.Engine
}

func NewUsersTableEngine(db *sqlx.DB) *UsersTableEngine {
	return &UsersTableEngine{
		Engine: model.NewEngine(db, "users"),
	}
}

func (e *UsersTableEngine) Insert(identity string, password string, status string) (*EntityUser, error) {
	saltBytes := make([]byte, PW_SALT_BYTES)
	if _, err := io.ReadFull(rand.Reader, saltBytes); err != nil {
		log.Fatal(err)
	}
	salt := string(saltBytes)

	hash := sha512.Sum512([]byte(password + salt))

	u := &EntityUser{
		Prototype: model.NewPrototype(),
		Identity:  identity,
		Password:  string(hash[:]),
		Salt:      salt,
		Status:    status,
	}

	_, err := e.Engine.Insert(u)
	return u, err
}

func (e *UsersTableEngine) GetByID(id string) (*EntityUser, error) {
	row := EntityUser{}
	q := `SELECT * FROM ` + e.TblName + ` WHERE id = UUID_TO_BIN(?);`
	if err := e.Engine.Get(&row, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &row, nil
}

func (e *UsersTableEngine) GetByIdentity(identity string) (*EntityUser, error) {
	row := EntityUser{}
	q := `SELECT * FROM ` + e.TblName + ` WHERE identity = ?;`
	if err := e.Engine.Get(&row, q, identity); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &row, nil
}
