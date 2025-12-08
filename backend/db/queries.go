package db

// Requêtes SQL brutes pour les operations comme backup etc


var (
				//  =============     DATABASES     =============  //

// ---- L'ajout d'une base des données -----
	QueryInsertDatabase =
		`INSERT INTO databases (user_id, name, type, host, port, db_username, db_password, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id`
		
// ---- Récupérer les bases d'un utilisateur -----
	QuerySelectDatabases = 
		`SELECT id, name, type, host, port, db_username FROM databases WHERE user_id = $1 ORDER BY id DESC`	
// ---- Suppression d'une database		
	QueryDeleteDatabase = 
	`DELETE FROM databases WHERE id=$1 AND user_id = $2`

// ---- Récupération informations database
	QuerySelectDatabaseInfo =
	`SELECT type, host, port, db_username, db_password, name FROM databases WHERE id = $1 AND user_id = $2`

// ---- Récupération des bases sans user_id pour CRON
	QuerySelectDatabaseInfoNoUser =
	`SELECT type, host, port, db_username, db_password, name FROM databases WHERE id =$1`	

	
					//   =============    BACKUPS    =============    //
// ---- Nouveau sauvegarde
	QueryInsertBackup =
	`INSERT INTO backups (database_id, name, status, version, file_path, backup_date)
	VALUES ($1, $2, 'pending', $3, $4, NOW())
	RETURNING id `
	
// ---- Mis à jour comme completé
	QueryUpdateBackupSuccess = 
	`UPDATE backups SET status = 'success', file_path = $1, file_size = $2, log = $3, backup_date = NOW() WHERE id = $4`

// ---- Mis à jour en cas d'erreur 
	QueryUpdateBackupFail =
	`UPDATE backups SET status = 'error', log = $1, backup_date = NOW() WHERE id = $2`

// ---- Récupération du backup et l'atacher au propriétaire 
	QuerySelectBackupWithOwnership =
	`SELECT b.id, b.status FROM backups b JOIN databases d ON b.database_id = d.id WHERE b.id = $1 AND d.user_id = $2 LIMIT 1`

// ---- Liste des backups d'un utilisateur
	QuerySelectAllBackups =
	`SELECT b.id, b.database_id, b.name, b.file_path, b.file_size, b.backup_date, b.status, b.version, b.log FROM backups b
	JOIN databases d ON d.id = database_id WHERE user_id = $1 ORDER BY b.backup_date DESC`

// ---- Liste des backups d'une seule base de données 
	QuerySelectBackupsByDatabase =
	`SELECT id, database_id, name, file_path, file_size, backup_date, status, version, log FROM backups WHERE database_id = $1 ORDER BY backup_date DESC`

// ---- Vérification de la propriété d'une base
	QueryCheckDatabaseOwnership =
	`SELECT COUNT(1) FROM databases WHERE id = $1 AND user_id = $2`

// ---- Récupérer le chemin du backup avec vérification propriétaire
	QueryGetBackupPathWithOwnership = 
	`SELECT b.file_path, b.name FROM backups b JOIN databases d ON b.database_id = d.id WHERE b.id = $1 AND d.user_id = $2 LIMIT 1 `

// ---- Récupération du chemin d'accès du fichier backup
	QueryFindBackupPathFile =
	`SELECT b.file_path FROM backups b JOIN databases d ON database_id = d.id WHERE b.id = $1 AND d.user_id = $2`


					//     =============    CRON    =============    //

// ---- Ajouter une nouvelle tâche planifiée
	QueryCreateScheduledTask = 
	`INSERT INTO scheduled_tasks (user_id, database_id, cron_expression, is_active) VALUES ($1, $2, $3, true) RETURNING id;`		

// ---- Récupérer toutes les tâches d'un utilisateur	donné
	QueryGetUserScheduledTasks =
	`SELECT id, user_id, database_id, cron_expression, is_active, last_run_at FROM scheduled_tasks WHERE user_id = $1;`

// ---- Mettre à jour l'expression CRON d'une tâche
	QueryUpdateCronExpression =
	`UPDATE scheduled_tasks SET cron_expression = $1 WHERE id = $2;`

// ---- Activer ou désactiver une tâche planifiée
	QueryToggleTaskActive =
	`UPDATE scheduled_tasks SET is_active = $1 WHERE id = $2;`
	
// ---- Mettre à jour la date d'exécution d'une tâche
	QueryUpdateLastTaskRun =
	`UPDATE scheduled_tasks SET last_run_at = NOW() WHERE id = $1;`
	
// ---- Récupèrer toutes les tâches actives
	QueryGetAllActiveTasks =
	`SELECT id, user_id, database_id, cron_expression FROM scheduled_tasks WHERE is_active = true;`	

// ---- Suppression d'une tâche planifiée
	QueryDeleteScheduledTask = 
	`DELETE from scheduled_tasks WHERE id = $1`	

					//     =============    RESTORES    =============    //

// ---- Eféctuer une restauration
	QueryInsertRestore =
	`INSERT INTO restores (backup_id, user_id, status) VALUES ($1, $2, 'restoring') RETURNING id;`
	
// ---- Mettre à jour le status de la restauration => succès 
	QueryUpdateRestoreSuccess =
	`UPDATE restores SET status = 'success', restored_at = NOW(), log = $1 WHERE id = $2;`
	
// ----	Mettre à jour le status de la restauration => échec
	QueryUpdateRestoreError = 
	`UPDATE restores SET status = 'error', restored_at = NOW(), log = $1 WHERE id = $2;`

// ---- Récupération d'informations du backup utilisé pour la restauration
	QueryGetBackupByID =
	`SELECT b.id, b.database_id, b.file_path, d.type, d.host, d.port, d.db_username, d.db_password, d.name AS database_name 
	FROM backups b JOIN databases d ON b.database_id = d.id WHERE b.id = $1 AND user_id=$2;`
	
					//     =============    ALERTS    =============    //

// ---- Récupération des notifications
	QueryGetAllAlerts =
	`SELECT id, alert_type, message, created_at, is_read FROM alerts WHERE user_id = $1 ORDER BY created_at DESC;`
	
// ---- Récupération de la date du dernier backup pour chaque base, par utilisateur, retourne: user_id, db_name, last_backup
	QueryGetDatabasesLastBackup =
	`SELECT u.id AS user_id, d.name AS db_name, MAX(b.backup_date) AS last_backup FROM databases d JOIN users u ON d.user_id = u.id
	LEFT JOIN backups b ON b.database_id = d.id GROUP BY u.id, d.name;`
)
