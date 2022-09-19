-- name: CreateBankAccount :one
INSERT INTO bank_accounts (owner, balance, currency) VALUES ($1, $2, $3) RETURNING *;

-- name: GetBankAccount :one
SELECT * FROM bank_accounts WHERE id = $1 LIMIT 1;

-- name: GetBankAccountForUpdate :one
SELECT * FROM bank_accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListBankAccounts :many
SELECT * FROM bank_accounts WHERE owner = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: UpdateBankAccount :one
UPDATE bank_accounts SET balance = $2 WHERE id = $1 RETURNING *;

-- name: AddBankAccountBalance :one
UPDATE bank_accounts SET balance = balance + sqlc.arg(amount) WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteBankAccount :exec
DELETE FROM bank_accounts WHERE id = $1;
