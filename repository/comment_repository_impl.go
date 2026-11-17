package repository

import (
	"belajar-golang-database/entity"
	"context"
	"database/sql"
	"errors"
	"strconv"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

// FindAll implements [CommentRepository].
func (c *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments"
	rows, err := c.DB.QueryContext(ctx, script)

	comments := []entity.Comment{}
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		comment := entity.Comment{}
		if err := rows.Scan(&comment.Id, &comment.Email, &comment.Comment); err != nil {
			return comments, err
		}
		comments = append(comments, comment)
		if err := rows.Err(); err != nil {
			return comments, err
		}
	}
	return comments, nil
}

// FindById implements [CommentRepository].
func (c *commentRepositoryImpl) FindById(ctx context.Context, id int32) (entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments WHERE id = ? LIMIT 1"
	rows, err := c.DB.QueryContext(ctx, script, id)

	comment := entity.Comment{}
	if err != nil {
		return comment, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		return comment, nil
	} else {
		return comment, errors.New("id " + strconv.Itoa(int(id)) + " Not found")
	}
}

func (c *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	script := "INSERT INTO comments(email, comment) VALUES (?,?)"
	result, err := c.DB.ExecContext(ctx, script, comment.Email, comment.Comment)

	if err != nil {
		return comment, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}

	comment.Id = int32(insertId)
	return comment, nil
}
