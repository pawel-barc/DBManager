package db

// Requêtes SQL brutes pour les operations comme backup etc

// ---- L'ajout d'une base des données -----
var (

	QueryInsertDatabase =
		`INSERT INTO databases (user_id, name, type, host, port, db_username, db_password, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id`
)

