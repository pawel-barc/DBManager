package db

// Requêtes SQL brutes pour les operations comme backup etc


var (
// ---- L'ajout d'une base des données -----
	QueryInsertDatabase =
		`INSERT INTO databases (user_id, name, type, host, port, db_username, db_password, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id`
		
// ---- Récupérer les bases d'un utilisateur -----
	QuerySelectDatabases = 
		`SELECT id, name, type, host, port, db_username FROM databases WHERE user_id = $1 ORDER BY id DESC`	
		QueryDeleteDatabase = 
		`DELETE FROM databases WHERE id=$1 AND user_id = $2`
)

