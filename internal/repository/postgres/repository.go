package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/s21platform/user-service/internal/config"
	"log"
)

type Repository struct {
	conn *sql.DB
}

type CheckUser struct {
	Uuid  string
	IsNew bool
}

func (r *Repository) IsUserExistByUUID(uuid string) (bool, error) {
	var exists bool
	row := r.conn.QueryRow("SELECT 1 FROM users WHERE uuid=$1", uuid)
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
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
	row := r.conn.QueryRow("SELECT uuid FROM users WHERE email=$1", email)
	if err := row.Scan(&userUuid); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("For user: %s - not found row in DB. Creating...\n", email)
			uuid_, err := r.createUser(email)
			if err != nil {
				log.Printf("For user: %s - unknown error while create new uuid\n", email)
				return nil, err
			}
			log.Printf("For user: %s - user created\n", email)
			return &CheckUser{Uuid: uuid_, IsNew: true}, nil
		}
		log.Printf("For user: %s - unknown error\n", email)
		return nil, err
	}
	log.Printf("For user: %s - exist. ok: %s!\n", email, userUuid)
	return &CheckUser{Uuid: userUuid, IsNew: false}, nil
}

func (r *Repository) createUser(email string) (string, error) {
	// TODO необходимо сюда прокидывать логин
	uuid_, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	_, err = r.conn.Exec("INSERT INTO users (login, uuid, email) VALUES ('Test', $1, $2)", uuid_.String(), email)
	if err != nil {
		return "", err
	}
	return uuid_.String(), nil
}

func (r *Repository) Close() {
	_ = r.conn.Close()
}

func New(cfg *config.Config) *Repository {
	//Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal("error connect: ", err)
	}

	if err := conn.Ping(); err != nil {
		log.Println("error ping: ", err)
	}
	return &Repository{conn}
}
