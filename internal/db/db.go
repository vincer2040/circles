package db

import (
	"context"
	"database/sql"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/vincer2040/circles/internal/post"
	"github.com/vincer2040/circles/internal/user"
	_ "modernc.org/sqlite"
)

type CirclesDB struct {
	db  *sql.DB
	ctx context.Context
}

func New(url string) (*CirclesDB, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return &CirclesDB{db, ctx}, nil
}

func (cdb *CirclesDB) CreateUserTable() error {
	_, err := cdb.exec(
		`CREATE TABLE IF NOT EXISTS
        users(
            first TEXT NOT NULL,
            last TEXT NOT NULL,
            email TEXT NOT NULL,
            password TEXT NOT NULL,
            PRIMARY KEY(email)
        )`,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) InsertUser(user *user.User) error {
	_, err := cdb.exec(
		`INSERT INTO
        users(first, last, email, password)
        VALUES (?, ?, ?, ?)`,
		user.First, user.Last, user.Email, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) GetUser(email string) (*user.User, error) {
	var user user.User
	query := "SELECT first, last, email, password FROM users WHERE email = ?"
	row := cdb.db.QueryRow(query, email)
	err := row.Scan(&user.First, &user.Last, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (cdb *CirclesDB) DeleteUser(email string) error {
	_, err := cdb.exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) DropUserTable() error {
	_, err := cdb.exec("DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) CreateCirclesTable() error {
	_, err := cdb.exec(
		`CREATE TABLE IF NOT EXISTS
        circles(
            name TEXT NOT NULL,
            creator TEXT NOT NULL,
            PRIMARY KEY(name),
            FOREIGN KEY(creator) REFERENCES users(email) ON DELETE CASCADE
        )
        `,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) InsertCircle(name, creator string) error {
	_, err := cdb.exec(
		`INSERT INTO
        circles(name, creator)
        VALUES(?, ?);
        INSERT INTO
        circlesUsers(name, email)
        VALUES(?, ?)
        `,
		name,
		creator,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) GetCreatorCircles(email string) ([]string, error) {
	query :=
		`SELECT name
        FROM circles
        WHERE creator = ?
        `
	rows, err := cdb.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var circles []string
	for rows.Next() {
		var circle string
		if err := rows.Scan(&circle); err != nil {
			return nil, err
		}
		circles = append(circles, circle)
	}
	return circles, nil
}

func (cdb *CirclesDB) DropCirclesTable() error {
	_, err := cdb.exec("DROP TABLE IF EXISTS circles")
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) CreateCircleUsersTable() error {
	_, err := cdb.exec(
		`CREATE TABLE IF NOT EXISTS
        circlesUsers(
            name TEXT NOT NULL,
            email TEXT NOT NULL,
            FOREIGN KEY(name) REFERENCES circles(name) ON DELETE CASCADE,
            FOREIGN KEY(email) REFERENCES circles(email) ON DELETE CASCADE
        )
        `,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) InsertUserToCircle(circle, email string) error {
	_, err := cdb.exec(
		`INSERT INTO
        circlesUsers(name, email)
        VALUES(?, ?)
        `,
		circle,
		email,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) GetUsersCircles(email string) ([]string, error) {
	query :=
		`SELECT name
    FROM circlesUsers
    WHERE email = ?
    `
	rows, err := cdb.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var circles []string
	for rows.Next() {
		var circle string
		if err := rows.Scan(&circle); err != nil {
			return nil, err
		}

		circles = append(circles, circle)
	}
	return circles, nil
}

func (cdb *CirclesDB) UserIsInCircle(circle, email string) (bool, error) {
	query :=
		`SELECT email
    FROM circlesUsers
    WHERE name = ?
    AND email = ?
    `

	rows, err := cdb.db.Query(query, circle, email)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	var emailFromDB string

	for rows.Next() {
		if err := rows.Scan(&emailFromDB); err != nil {
			return false, err
		}
	}
	if emailFromDB == email {
		return true, nil
	}
	return false, nil
}

func (cdb *CirclesDB) DropCirclesUsersTable() error {
	_, err := cdb.exec("DROP TABLE IF EXISTS circlesUsers")
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) CreatePostsTable() error {
	_, err := cdb.exec(
		`CREATE TABLE IF NOT EXISTS
        posts(
            circle TEXT NOT NULL,
            author TEXT NOT NULL,
            description TEXT NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(circle) REFERENCES circles(name) ON DELETE CASCADE,
            FOREIGN KEY(author) REFERENCES users(email) ON DELETE CASCADE
        )
        `,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) InsertPost(circle, author, description, timestamp string) error {
	_, err := cdb.exec(
		`INSERT INTO
        posts(circle, author, description, timestamp)
        VALUES (?, ?, ?, ?)
        `,
		circle,
		author,
		description,
		timestamp,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) DeletePost(id int64, author string) error  {
	_, err := cdb.exec(
        `DELETE FROM posts
        WHERE rowid = ?
        `,
		id,
	)

	if err != nil {
		return err
	}
    return nil
}

func (cdb *CirclesDB) GetPostsForCircle(circle string) ([]post.PostFromDB, error) {
	query :=
		`SELECT posts.rowid, first, description, timestamp
        FROM posts
        INNER JOIN users on posts.author=users.email
        WHERE circle = ?
        ORDER BY
        timestamp DESC
        `
	rows, err := cdb.db.Query(query, circle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []post.PostFromDB
	for rows.Next() {
		var post post.PostFromDB
		if err := rows.Scan(&post.ID, &post.Author, &post.Description, &post.TimeStamp); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (cdb *CirclesDB) GetPostsForUser(email string) ([]post.UserPost, error) {
	query :=
		`SELECT posts.rowid, circle, description, timestamp
        FROM posts
        WHERE author = ?
        ORDER BY
        timestamp DESC
        `
	rows, err := cdb.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	var posts []post.UserPost
	defer rows.Close()
	for rows.Next() {
		var post post.UserPost
		if err := rows.Scan(&post.ID, &post.Circle, &post.Description, &post.TimeStamp); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (cdb *CirclesDB) DropPostsTable() error {
	_, err := cdb.exec("DROP TABLE IF EXISTS posts")
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) Close() {
	cdb.db.Close()
}

func (cdb *CirclesDB) exec(stmt string, args ...any) (*sql.Result, error) {
	res, err := cdb.db.ExecContext(cdb.ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (cdb *CirclesDB) query(stmt string, args ...any) (*sql.Rows, error) {
	res, err := cdb.db.QueryContext(cdb.ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
