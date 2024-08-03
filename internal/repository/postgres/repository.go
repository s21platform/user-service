package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/s21platform/user-service/internal/config"
	"log"
)

type Repository struct {
	conn *sql.DB
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
