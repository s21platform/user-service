package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/model"
)

const defaultAvatar = "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg"

type Repository struct {
	conn *sqlx.DB
}

type CheckUser struct {
	Uuid  string
	IsNew bool
}

func (r *Repository) SetFriends(ctx context.Context, peer1, peer2 string) (bool, error) {
	res, err := r.isRowFriendsExist(ctx, peer1, peer2)
	if err != nil {
		return false, fmt.Errorf("failed to check an RowExist: %v", err)
	}
	if res {
		return false, nil
	}

	query := sq.Insert("friends").
		Columns("initiator", "user_id").
		Values(peer1, peer2).
		PlaceholderFormat(sq.Dollar)

	sqlString, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query string: %v", err)
	}

	_, err = r.conn.ExecContext(ctx, sqlString, args...)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %v", err)
	}
	return true, nil
}

func (r *Repository) RemoveFriends(ctx context.Context, peer1, peer2 string) (bool, error) {
	res, err := r.isRowFriendsExist(ctx, peer1, peer2)
	if err != nil {
		return false, fmt.Errorf("failed to check an RowExist: %v", err)
	}
	if !res {
		return false, nil
	}

	query := sq.Delete("friends").
		Where(sq.Eq{"initiator": peer1, "user_id": peer2}).
		PlaceholderFormat(sq.Dollar)

	sqlString, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query string: %v", err)
	}
	_, err = r.conn.ExecContext(ctx, sqlString, args...)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %v", err)
	}
	return true, nil
}

func (r *Repository) isRowFriendsExist(ctx context.Context, peer1, peer2 string) (bool, error) {
	var exists bool

	query := sq.Select("COUNT(1) > 0").
		From("friends").
		Where(sq.Eq{"initiator": peer1, "user_id": peer2}).
		PlaceholderFormat(sq.Dollar)

	sqlString, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = r.conn.GetContext(ctx, &exists, sqlString, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check friends existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) UpdateUserAvatar(uuid, link string) error {
	query := `UPDATE users SET last_avatar_link = $1 WHERE uuid = $2`
	_, err := r.conn.Exec(query, link, uuid)
	if err != nil {
		return err
	}

	return nil
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
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to get id of inserted row: %v", err)
	}
	_, err = tx.Exec("INSERT INTO data (user_id) VALUES ($1)", lastId)
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to insert data: %v", err)
	}
	_ = tx.Commit()
	return uuid_.String(), nil
}

func (r *Repository) GetUserInfoByUUID(ctx context.Context, uuid string) (model.UserInfo, error) {
	query := `
		SELECT 
		    u.login,
		    u.last_avatar_link,
		    d.name,
		    d.surname,
		    d.birthdate,
		    d.phone,
		    d.telegram,
		    d.git,
		    d.city_id,
		    d.os_id,
		    d.work_id,
		    d.university_id
		FROM users u 
		JOIN data d ON d.user_id = u.id
		where u.uuid = $1
	`
	var result []model.UserInfo
	err := r.conn.Select(&result, query, uuid)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("failed to get user info: %v", err)
	}
	if len(result) == 0 {
		return model.UserInfo{}, errors.New("user not found")
	}
	return result[0], nil
}

func (r *Repository) GetUsersByUUID(uuid string) (model.UserInfoMin, error) {
	query := "select uuid, login, last_avatar_link, " +
		"COALESCE(data.name, '') name, COALESCE(data.surname, '') as surname from users " +
		"join data ON users.id = data.user_id " +
		"where uuid = $1"
	var result model.UserInfoMin
	err := r.conn.Get(&result, query, uuid)
	if err != nil {
		return model.UserInfoMin{}, fmt.Errorf("failed to get user info: %v", err)
	}
	return result, nil
}

func (r *Repository) GetUserWithLimit(uuid, nickname string, limit int64, offset int64) ([]model.UserWithLimit, int64, error) {
	var userWithLimit []model.UserWithLimit
	likeNick := "%" + nickname + "%"
	err := r.conn.Select(&userWithLimit, "SELECT users.login, users.uuid, users.last_avatar_link, COALESCE(data.name, '') as name, COALESCE(data.surname, '') as surname "+
		"FROM users "+
		"JOIN data on users.id = data.user_id "+
		"WHERE users.uuid != $1 AND users.login LIKE $2 LIMIT $3 OFFSET $4;", uuid, likeNick, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user info: %v", err)
	}
	var total int64
	query := `
		SELECT
			count(uuid)
		FROM users 
		WHERE uuid != $1 AND login LIKE $2
	`
	err = r.conn.Get(&total, query, uuid, likeNick)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user total: %v", err)
	}
	return userWithLimit, total, nil
}

func (r *Repository) GetLoginByUuid(ctx context.Context, uuid string) (string, error) {
	query := `
		SELECT
			login 
		FROM users 
		WHERE uuid=$1
	`
	var result string
	err := r.conn.Get(&result, query, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get login by uuid: %v", err)
	}
	return result, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, data model.ProfileData, userUuid string) error {
	tx, err := r.conn.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	query, args, err := sq.Select("id").
		From("users").Where(sq.Eq{"uuid": userUuid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to build query: %v", err)
	}
	var userId int64
	err = tx.GetContext(ctx, &userId, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get user id: %v", err)
	}
	log.Println(data.Birthdate)
	query, args, err = sq.Update("data").
		Set("name", data.Name).
		Set("birthdate", data.Birthdate).
		Set("git", data.Git).
		Set("telegram", data.Telegram).
		Set("os_id", data.OsId).
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update user: %v", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update user: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
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
		log.Fatal("error ping: ", err)
	}
	return &Repository{conn}
}
