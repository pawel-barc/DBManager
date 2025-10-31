package db

// Requêtes SQL brutes pour les opérations sur les bases et backups

var (
	// ---- Ajout d'une base de données ----
	QueryInsertDatabase = `
		INSERT INTO databases (user_id, name, type, host, port, db_username, db_password, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id
	`

	// ---- Vérification d'appartenance (sécurité) ----
	QueryCheckDatabaseOwnership = `
		SELECT COUNT(*) FROM databases
		WHERE id = $1 AND user_id = $2
	`

	// ---- Insertion d'un nouveau backup ----
	QueryInsertBackup = `
		INSERT INTO backups (database_id, name, status, version)
		VALUES ($1, $2, 'queued', $3)
		RETURNING id
	`

	// ---- Liste des backups pour une base ----
	QueryListBackupsByDatabase = `
		SELECT id, database_id, name, file_path, file_size, backup_date, status, version, log
		FROM backups
		WHERE database_id = $1
		ORDER BY id DESC
	`

	// ---- Récupérer le statut d'un backup (avec ownership) ----
	QueryGetBackupStatusWithOwnership = `
		SELECT b.id, b.status
		FROM backups b
		JOIN databases d ON d.id = b.database_id
		WHERE b.id = $1 AND d.user_id = $2
	`

	// ---- Récupérer le log d'un backup ----
	QueryGetBackupLogWithOwnership = `
		SELECT b.id, b.log
		FROM backups b
		JOIN databases d ON d.id = b.database_id
		WHERE b.id = $1 AND d.user_id = $2
	`

	// ---- Récupérer le chemin du fichier pour téléchargement ----
	QueryGetBackupPathWithOwnership = `
		SELECT b.file_path, b.name
		FROM backups b
		JOIN databases d ON d.id = b.database_id
		WHERE b.id = $1 AND d.user_id = $2 AND b.status = 'done'
	`

	// ---- Mise à jour pendant exécution ----
	QueryMarkBackupRunning = `
		UPDATE backups
		SET status = 'running'
		WHERE id = $1
	`

	QueryMarkBackupDone = `
		UPDATE backups
		SET status = 'done',
		    file_path = $2,
		    file_size = $3,
		    backup_date = NOW(),
		    log = COALESCE(log, '')
		WHERE id = $1
	`

	QueryMarkBackupError = `
		UPDATE backups
		SET status = 'error',
		    log = COALESCE(log, '') || E'\n' || $2
		WHERE id = $1
	`
)
