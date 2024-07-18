-- name: FindPersonByDocument :one
SELECT uuid, name, document, created_at, updated_at FROM person.person 
WHERE document = $1 AND active = true
LIMIT 1;


-- name: InsertPerson :exec
INSERT INTO person.person (uuid, name, document, active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);
