package database

import (
	"context"
	"database/sql"
	"log"

	"curso-rest.com/go/rest/models"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type MysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(url string) (*MysqlRepository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	return &MysqlRepository{db}, nil
}

func (mysql *MysqlRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := mysql.db.ExecContext(ctx, "INSERT INTO users(id,email, password) VALUES (?,?,?)", user.Id, user.Email, user.Password)
	return err
}

func (mysql *MysqlRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := mysql.db.QueryContext(ctx, "SELECT id, email FROM users WHERE users.id = ?", id)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (mysql *MysqlRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := mysql.db.QueryContext(ctx, "SELECT id,email, password FROM users WHERE users.email = ?", email)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}
func (mysql *MysqlRepository) Close() error {
	return mysql.db.Close()
}
