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

func (mysql *MysqlRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := mysql.db.ExecContext(ctx, "INSERT INTO posts(id, post_content, user_id) VALUES (?,?,?)", post.Id, post.PostContent, post.UserId)
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

func (mysql *MysqlRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := mysql.db.QueryContext(ctx, "SELECT * FROM posts WHERE posts.id = ?", id)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var post = models.Post{}
	for rows.Next() {
		if err = rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId); err == nil {
			return &post, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &post, nil
}

func (repo *MysqlRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = ? WHERE id = ? and user_id = ?", post.PostContent, post.Id, post.UserId)
	return err
}

func (repo *MysqlRepository) DeletePostById(ctx context.Context, id string, userId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE posts.id = ? AND posts.user_id = ?", id, userId)
	return err
}

func (repo *MysqlRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	//can you configurate how much numbers of element select by pagination
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts LIMIT ? OFFSET ?", 5, page*5)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var posts []*models.Post
	for rows.Next() {
		var post = models.Post{}
		if err = rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt); err == nil {
			posts = append(posts, &post)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
