package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/s21platform/user-service/internal/config"
	"log"
	"strings"
)

const defaultAvatar = "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg"

type Repository struct {
	conn *sqlx.DB
}

type CheckUser struct {
	Uuid  string
	IsNew bool
}

func (r *Repository) IsUserExistByUUID(uuid string) (bool, error) {
	var exists bool
	row := r.conn.QueryRow("SELECT 1 FROM users WHERE uuid=$1", uuid)
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("For user: %s - not found row in DB\n", uuid)
			return false, nil
		}
		log.Printf("For user: %s - unknown error\n", uuid)
		return false, err
	}
	log.Printf("For user: %s - exist. ok!\n", uuid)
	return exists, nil
}

func (r *Repository) GetOrSetUserByLogin(email string) (*CheckUser, error) {
	var userUuid string
	nickname, err := checkEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %v", err)
	}
	err = r.conn.Get(&userUuid, "SELECT uuid FROM users WHERE email=$1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uuid_, err := r.createUser(nickname, email)
			if err != nil {
				return nil, fmt.Errorf("failed create user: %v", err)
			}
			return &CheckUser{Uuid: uuid_, IsNew: true}, nil
		}
		return nil, fmt.Errorf("error checking user: %v", err)
	}
	log.Printf("For user: %s - exist. ok: %s!\n", email, userUuid)
	return &CheckUser{Uuid: userUuid, IsNew: false}, nil
}

func checkEmail(email string) (string, error) {
	res := strings.Split(email, "@")
	if len(res) != 2 {
		return "", fmt.Errorf("checkEmail, %s is not email", email)
	}
	// TODO Тут пропускаются только со школьной почтой
	if res[1] != "student.21-school.ru" {
		return "", fmt.Errorf("checkEmail, %s is not 21-school email", email)
	}
	return res[0], nil
}

func (r *Repository) createUser(nickname, email string) (string, error) {
	uuid_, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	tx, err := r.conn.Beginx()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	var lastId int
	err = tx.QueryRowx("INSERT INTO users (login, uuid, email, last_avatar_link) VALUES ($1, $2, $3, $4) RETURNING id", nickname, uuid_.String(), email, defaultAvatar).Scan(&lastId)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to get id of inserted row: %v", err)
	}
	_, err = r.conn.Exec("INSERT INTO data (user_id) VALUES ($1)", lastId)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to insert data: %v", err)
	}
	tx.Commit()
	return uuid_.String(), nil
}

func (r *Repository) Close() {
	_ = r.conn.Close()
}

func New(cfg *config.Config) *Repository {
	//Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatal("error connect: ", err)
	}

	if err := conn.Ping(); err != nil {
		log.Println("error ping: ", err)
	}
	return &Repository{conn}
}
