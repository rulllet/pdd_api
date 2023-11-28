// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: query.sql

package data

import (
	"context"
)

const getAnswers = `-- name: GetAnswers :many
SELECT id, title, correct_answer, question_id
FROM answer
WHERE question_id = ?
`

func (q *Queries) GetAnswers(ctx context.Context, questionID string) ([]Answer, error) {
	rows, err := q.db.QueryContext(ctx, getAnswers, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Answer
	for rows.Next() {
		var i Answer
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.CorrectAnswer,
			&i.QuestionID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuestion = `-- name: GetQuestion :one
SELECT id, ticket, number, title, help, image, category_id
FROM question
WHERE category_id = ? AND ticket = ? AND number = ?
`

type GetQuestionParams struct {
	CategoryID string `json:"category_id"`
	Ticket     int64  `json:"ticket"`
	Number     int64  `json:"number"`
}

func (q *Queries) GetQuestion(ctx context.Context, arg GetQuestionParams) (Question, error) {
	row := q.db.QueryRowContext(ctx, getQuestion, arg.CategoryID, arg.Ticket, arg.Number)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.Ticket,
		&i.Number,
		&i.Title,
		&i.Help,
		&i.Image,
		&i.CategoryID,
	)
	return i, err
}