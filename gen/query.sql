-- name: GetQuestion :one
SELECT *
FROM question
WHERE category_id = ? AND ticket = ? AND number = ?;

-- name: GetAnswers :many
SELECT *
FROM answer
WHERE question_id = ?;
