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
	*sqlx.DB
}

type CheckUser struct {
	Uuid  string
	IsNew bool
}

func New(cfg *config.Config) *Repository {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	return &Repository{conn}
}

func (r *Repository) Close() {
	_ = r.DB.Close()
}

func (r *Repository) GetPeerFollow(ctx context.Context, userUUID string) ([]string, error) {
	sqlString, args, err := sq.Select("invited").
		From("friends").
		Where(sq.Eq{"initiator": userUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query string: %w", err)
	}

	var follows []string
	if err := r.Chk(ctx).SelectContext(ctx, &follows, sqlString, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch follows: %w", err)
	}
	return follows, nil
}

func (r *Repository) GetWhoFollowPeer(ctx context.Context, userUUID string) ([]string, error) {
	sqlString, args, err := sq.Select("initiator").
		From("friends").
		Where(sq.Eq{"invited": userUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query string: %w", err)
	}

	var follows []string
	if err := r.Chk(ctx).SelectContext(ctx, &follows, sqlString, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch follows: %w", err)
	}
	return follows, nil
}

func (r *Repository) GetSubscribersCount(ctx context.Context, userUUID string) (int64, error) {
	sqlString, args, err := sq.Select("COUNT(initiator)").
		From("friends").
		Where(sq.Eq{"invited": userUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build query string: %v", err)
	}

	var count int64
	if err := r.Chk(ctx).GetContext(ctx, &count, sqlString, args...); err != nil {
		return 0, fmt.Errorf("failed to get subscribers count: %w", err)
	}

	return count, nil
}

func (r *Repository) GetSubscriptionCount(ctx context.Context, userUUID string) (int64, error) {
	sqlString, args, err := sq.Select("COUNT(initiator)").
		From("friends").
		Where(sq.Eq{"initiator": userUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build query string: %v", err)
	}

	var count int64
	if err := r.Chk(ctx).GetContext(ctx, &count, sqlString, args...); err != nil {
		return 0, fmt.Errorf("failed to get subscription count: %w", err)
	}

	return count, nil
}

func (r *Repository) SetFriends(ctx context.Context, peer1, peer2 string) error {
	sqlString, args, err := sq.Insert("friends").
		Columns("initiator", "invited").
		Values(peer1, peer2).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query string: %v", err)
	}

	_, err = r.Chk(ctx).ExecContext(ctx, sqlString, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	return nil
}

func (r *Repository) RemoveFriends(ctx context.Context, peer1, peer2 string) error {
	sqlString, args, err := sq.Delete("friends").
		Where(sq.Eq{"initiator": peer1, "invited": peer2}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query string: %v", err)
	}
	_, err = r.Chk(ctx).ExecContext(ctx, sqlString, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	return nil
}

func (r *Repository) CheckFriendship(ctx context.Context, peer1, peer2 string) (bool, error) {
	var exists bool

	sqlString, args, err := sq.Select("COUNT(1) > 0").
		From("friends").
		Where(sq.Eq{"initiator": peer1, "invited": peer2}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = r.Chk(ctx).GetContext(ctx, &exists, sqlString, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check friends existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) UpdateUserAvatar(uuid, link string) error {
	query := `UPDATE users SET last_avatar_link = $1 WHERE uuid = $2`
	_, err := r.Exec(query, link, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) IsUserExistByUUID(uuid string) (bool, error) {
	var exists bool
	row := r.QueryRow("SELECT 1 FROM users WHERE uuid=$1", uuid)
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
	err = r.Get(&userUuid, "SELECT uuid FROM users WHERE email=$1", email)
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
	tx, err := r.Beginx()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	_, err = tx.Exec("INSERT INTO users (login, uuid, email, last_avatar_link) VALUES ($1, $2, $3, $4)", nickname, uuid_.String(), email, defaultAvatar)
	if err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("failed to get id of inserted row: %v", err)
	}
	_, err = tx.Exec("INSERT INTO data (user_uuid) VALUES ($1)", uuid_)
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
		JOIN data d ON d.user_uuid = u.uuid
		where u.uuid = $1
	`
	var result []model.UserInfo
	err := r.Chk(ctx).SelectContext(ctx, &result, query, uuid)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("failed to get user info: %v", err)
	}
	if len(result) == 0 {
		return model.UserInfo{}, errors.New("user not found")
	}
	return result[0], nil
}

func (r *Repository) GetUsersByUUID(uuid string) (model.UserInfoMin, error) { // TODO: wrong naming, it gets 1 user info, not list
	query := "select uuid, login, last_avatar_link, " +
		"COALESCE(data.name, '') name, COALESCE(data.surname, '') as surname from users " +
		"join data ON users.uuid = data.user_uuid " +
		"where uuid = $1"
	var result model.UserInfoMin
	err := r.Get(&result, query, uuid)
	if err != nil {
		return model.UserInfoMin{}, fmt.Errorf("failed to get user info: %v", err)
	}
	return result, nil
}

func (r *Repository) GetUserWithLimit(uuid, nickname string, limit int64, offset int64) ([]model.UserWithLimit, int64, error) { // TODO: wrong naming, it gets list of user, not 1 user
	var userWithLimit []model.UserWithLimit
	likeNick := "%" + nickname + "%"
	err := r.Select(&userWithLimit, "SELECT users.login, users.uuid, users.last_avatar_link, COALESCE(data.name, '') as name, COALESCE(data.surname, '') as surname "+
		"FROM users "+
		"JOIN data on users.uuid = data.user_uuid "+
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
	err = r.Get(&total, query, uuid, likeNick)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user total: %v", err)
	}
	return userWithLimit, total, nil
}

func (r *Repository) GetLoginByUuid(ctx context.Context, uuid string) (string, error) {
	query := `
		SELECT login 
		FROM users 
		WHERE uuid=$1
	`
	var result string
	err := r.Get(&result, query, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get login by uuid: %v", err)
	}
	return result, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, data model.ProfileData, userUuid string) error {
	tx, err := r.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	log.Println(data.Birthdate)
	query, args, err := sq.Update("data").
		Set("name", data.Name).
		Set("birthdate", data.Birthdate).
		Set("git", data.Git).
		Set("telegram", data.Telegram).
		Set("os_id", data.OsId).
		Where(sq.Eq{"user_uuid": userUuid}).
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

func (r *Repository) GetUserForCreation(ctx context.Context, email string) (*model.UserAuthInfo, error) {
	query, args, err := sq.
		Select(`uuid`, `login`).
		From(`users`).
		Where(sq.Eq{`email`: email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var user model.UserAuthInfo
	err = r.Chk(ctx).GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user by email: %v", err)
	}

	return &user, nil
}

func (r *Repository) CheckNicknameAvailability(ctx context.Context, nickname string) (bool, error) {
	query, args, err := sq.Select(`login`).
		From("users").
		Where(sq.Eq{"login": nickname}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %v", err)
	}

	var login string
	err = r.Chk(ctx).GetContext(ctx, &login, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("failed to execute query: %v", err)
	}
	return false, nil
}

func (r *Repository) CreateUser(ctx context.Context, userUUID string, email string, nickname string) error {
	query, args, err := sq.
		Insert(`users`).
		Columns(`login`, `email`, `uuid`, `last_avatar_link`).
		Values(nickname, email, userUUID, defaultAvatar).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	_, err = r.Chk(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	return nil
}

func (r *Repository) CreateUserMeta(ctx context.Context, userUUID string) error {
	query, args, err := sq.
		Insert(`data`).
		Columns(`user_uuid`).
		Values(userUUID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	_, err = r.Chk(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	return nil
}

func (r *Repository) CreatePost(ctx context.Context, ownerUUID uuid.UUID, content string) (string, error) {
	query, args, err := sq.Insert("posts").
		Columns("user_uuid", "content").
		Values(ownerUUID, content).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("failed to build insert query: %v", err)
	}

	var newPostUUID string
	err = r.Chk(ctx).GetContext(ctx, &newPostUUID, query, args...)

	if err != nil {
		return "", fmt.Errorf("failed to create post: %v", err)
	}

	return newPostUUID, nil
}

func (r *Repository) GetPostsByIds(ctx context.Context, uuids []string) (*model.PostInfoList, error) {
	var posts model.PostInfoList

	query, args, err := sq.
		Select(
			"cast(posts.id as varchar) as id",
			"users.login as login",
			"coalesce(data.name, '') as name",
			"coalesce(data.surname, '') as surname",
			"users.last_avatar_link as last_avatar_link",
			"posts.content as content",
			"posts.created_at as created_at",
			"posts.updated_at as updated_at").
		From("posts").
		Join("users ON users.uuid = posts.user_uuid").
		Join("data ON data.user_uuid = users.uuid").
		Where(sq.And{
			sq.Eq{"posts.id": uuids},
			sq.Eq{"posts.deleted_at": nil}}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	err = r.Chk(ctx).SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return &posts, nil
}
