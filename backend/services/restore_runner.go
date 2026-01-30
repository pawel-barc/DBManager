package services

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"safebase/db"
	"safebase/utils"
	"strings"
	"time"
)

// RestoreBackup exécute de maniere synchrone une restauration à partir d'un fichier de backup
func RestoreBackup(backupID, userID int) (restoreID int, dbName string, err error) {
	// Insérer une entrée dans la table restores avec le statut initial
	if err := db.DB.QueryRow(db.QueryInsertRestore, backupID, userID).Scan(&restoreID); err != nil {
		utils.LogError("Impossible d'enregistrer la restauration:", err)
		return 0, "", fmt.Errorf("impossible d'enregistrer la restauration: %v", err)
	}

	// Récupérer les informations du backup
	var databaseID int
	var dbType, host, port, username, password, backupPath string
	if err := db.DB.QueryRow(db.QueryGetBackupByID, backupID, userID).Scan(
		&backupID, &databaseID, &backupPath, &dbType, &host, &port, &username, &password, &dbName); err != nil {
		msg := fmt.Sprintf("Backup introuvable ou accès refusé: %v", err)
		db.DB.Exec(db.QueryUpdateRestoreError, msg, restoreID)
		return restoreID, "", err
	}

	// Restaurer selon le type de base de données
	switch dbType {
	case "postgres":
		return restorePostgres(host, port, dbName, username, password, backupPath, restoreID, backupID)
	case "mysql":
		return restoreMySQL(host, port, dbName, username, password, backupPath, restoreID)
	default:
		errMsg := "Type de base non supporté"
		db.DB.Exec(db.QueryUpdateRestoreError, errMsg, restoreID)
		utils.LogError(errMsg, nil)
		return restoreID, "", fmt.Errorf(errMsg)
	}
}

// restorePostgres restaure une base de données PostgreSQL
func restorePostgres(host, port, dbName, username, password, backupPath string, restoreID, backupID int) (int, string, error) {
	// Get admin credentials
	adminUser := os.Getenv("PG_ADMIN_USER")
	adminPassword := os.Getenv("PG_ADMIN_PASSWORD")
	if adminUser == "" || adminPassword == "" {
		return restoreID, "", fmt.Errorf("admin database credentials missing")
	}

	// 1. Drop and recreate database
	if err := recreatePostgresDB(host, port, dbName, adminUser, adminPassword, username); err != nil {
		db.DB.Exec(db.QueryUpdateRestoreError, err.Error(), restoreID)
		utils.LogError("Database recreation failed", err)
		return restoreID, "", err
	}

	// 2. Clean and restore backup
	if err := restorePostgresBackup(host, port, dbName, adminUser, adminPassword, backupPath, restoreID); err != nil {
		return restoreID, "", err
	}

	// 3. Grant privileges to app user
	if err := grantPostgresPrivileges(host, port, dbName, adminUser, adminPassword, username); err != nil {
		utils.LogInfo(fmt.Sprintf("Warning: Failed to grant user privileges: %v", err))
	}

	// Success
	db.DB.Exec(db.QueryUpdateRestoreSuccess, "PostgreSQL restore completed", restoreID)
	utils.LogInfo(fmt.Sprintf("Restore backup #%d terminé avec succès", backupID))
	return restoreID, dbName, nil
}

// restorePostgresBackup exécute la restauration du fichier de backup
func restorePostgresBackup(host, port, dbName, adminUser, adminPassword, backupPath string, restoreID int) error {
	// Create cleaned backup file
	cleanedBackupPath, cleanup, err := cleanBackupFile(backupPath)
	if err != nil {
		db.DB.Exec(db.QueryUpdateRestoreError, err.Error(), restoreID)
		utils.LogError("Failed to clean backup file", err)
		return err
	}
	defer cleanup()

	// Run psql as admin
	os.Setenv("PGPASSWORD", adminPassword)
	cmd := exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", adminUser,
		"-d", dbName,
		"-v", "ON_ERROR_STOP=0",
		"--single-transaction",
		"-f", cleanedBackupPath,
	)
	cmd.Env = append(os.Environ(), "PGCLIENTENCODING=UTF-8", "PGPASSWORD="+adminPassword)

	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	utils.LogInfo(fmt.Sprintf("Restore output: %s", outputStr))

	if err != nil {
		// Check if restore was partially successful
		if isPartialSuccess(outputStr) {
			utils.LogInfo("Restore completed with warnings")
			return nil
		}
		
		logMsg := fmt.Sprintf("Restore error: %v, output: %s", err, outputStr)
		db.DB.Exec(db.QueryUpdateRestoreError, logMsg, restoreID)
		utils.LogError(logMsg, err)
		return fmt.Errorf(logMsg)
	}

	return nil
}

// isPartialSuccess vérifie si la restauration a partiellement réussi
func isPartialSuccess(output string) bool {
	return strings.Contains(output, "CREATE TABLE") || 
	       strings.Contains(output, "ALTER TABLE") ||
	       strings.Contains(output, "CREATE SEQUENCE") ||
	       strings.Contains(output, "COPY")
}

// recreatePostgresDB recrée la base de données PostgreSQL
func recreatePostgresDB(host, port, dbName, adminUser, adminPassword, appUser string) error {
	os.Setenv("PGPASSWORD", adminPassword)
	defer os.Unsetenv("PGPASSWORD")

	utils.LogInfo(fmt.Sprintf("Recreating PostgreSQL database: %s", dbName))

	// Terminate connections
	if err := terminateDBConnections(host, port, dbName, adminUser); err != nil {
		utils.LogInfo(fmt.Sprintf("Warning terminating connections: %v", err))
	}

	// Drop database
	if err := dropDatabase(host, port, dbName, adminUser); err != nil {
		utils.LogInfo(fmt.Sprintf("Warning dropping database: %v", err))
	}

	time.Sleep(2 * time.Second)

	// Create database
	if err := createDatabase(host, port, dbName, adminUser, "postgres"); err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// Grant initial privileges
	if err := grantDatabaseConnect(host, port, dbName, adminUser, appUser); err != nil {
		utils.LogInfo(fmt.Sprintf("Warning granting connect: %v", err))
	}

	utils.LogInfo(fmt.Sprintf("Successfully recreated database: %s", dbName))
	return nil
}

// terminateDBConnections termine toutes les connexions à la base de données
func terminateDBConnections(host, port, dbName, adminUser string) error {
	cmd := exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", adminUser,
		"-d", "postgres",
		"-v", "ON_ERROR_STOP=0",
		"-c", fmt.Sprintf(`
			REVOKE CONNECT ON DATABASE "%s" FROM PUBLIC;
			SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid();
			SELECT pg_sleep(2);
		`, dbName, dbName),
	)
	
	output, err := cmd.CombinedOutput()
	utils.LogInfo(fmt.Sprintf("Terminate connections: %s", string(output)))
	return err
}

// dropDatabase supprime la base de données
func dropDatabase(host, port, dbName, adminUser string) error {
	// First attempt
	cmd := exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", adminUser,
		"-d", "postgres",
		"-v", "ON_ERROR_STOP=0",
		"-c", fmt.Sprintf(`DROP DATABASE IF EXISTS "%s";`, dbName),
	)
	
	output, err := cmd.CombinedOutput()
	utils.LogInfo(fmt.Sprintf("Drop database: %s", string(output)))
	
	if err == nil {
		return nil
	}

	// Second attempt - more aggressive
	utils.LogInfo("First drop failed, trying aggressive approach...")
	cmd = exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", adminUser,
		"-d", "postgres",
		"-v", "ON_ERROR_STOP=0",
		"-c", fmt.Sprintf(`
			UPDATE pg_database SET datallowconn = 'false' WHERE datname = '%s';
			SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s';
			DROP DATABASE IF EXISTS "%s";
		`, dbName, dbName, dbName),
	)
	
	output, err = cmd.CombinedOutput()
	utils.LogInfo(fmt.Sprintf("Aggressive drop: %s", string(output)))
	return err
}

// createDatabase crée une nouvelle base de données
func createDatabase(host, port, dbName, adminUser, owner string) error {
	cmd := exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", adminUser,
		"-d", "postgres",
		"-v", "ON_ERROR_STOP=1",
		"-c", fmt.Sprintf(`CREATE DATABASE "%s" WITH OWNER = %s;`, dbName, owner),
	)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create database failed: %v\n%s", err, string(output))
	}
	return nil
}

// grantDatabaseConnect accorde le droit de connexion
func grantDatabaseConnect(host, port, dbName, adminUser, appUser string) error {
	cmd := exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", adminUser,
		"-d", "postgres",
		"-v", "ON_ERROR_STOP=0",
		"-c", fmt.Sprintf(`GRANT CONNECT ON DATABASE "%s" TO %s;`, dbName, appUser),
	)
	
	output, err := cmd.CombinedOutput()
	utils.LogInfo(fmt.Sprintf("Grant connect: %s", string(output)))
	return err
}

// grantPostgresPrivileges accorde les privilèges complets à l'utilisateur
func grantPostgresPrivileges(host, port, dbName, adminUser, adminPassword, appUser string) error {
	os.Setenv("PGPASSWORD", adminPassword)
	defer os.Unsetenv("PGPASSWORD")
	
	cmd := exec.Command(
		"psql",
		"-h", host,
		"-p", port,
		"-U", adminUser,
		"-d", dbName,
		"-v", "ON_ERROR_STOP=1",
		"-c", fmt.Sprintf(`
			GRANT CONNECT ON DATABASE "%s" TO %s;
			GRANT USAGE, CREATE ON SCHEMA public TO %s;
			GRANT ALL ON ALL TABLES IN SCHEMA public TO %s;
			GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO %s;
			GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO %s;
			ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO %s;
			ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO %s;
			ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO %s;
		`, dbName, appUser, appUser, appUser, appUser, appUser, appUser, appUser, appUser),
	)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("grant privileges failed: %v\n%s", err, string(output))
	}
	return nil
}

// restoreMySQL restaure une base de données MySQL
func restoreMySQL(host, port, dbName, username, password, backupPath string, restoreID int) (int, string, error) {
	file, err := os.Open(backupPath)
	if err != nil {
		db.DB.Exec(db.QueryUpdateRestoreError, err.Error(), restoreID)
		utils.LogError("Impossible d'ouvrir le fichier de backup mysql", err)
		return restoreID, "", err
	}
	defer file.Close()

	cmd := exec.Command(
		"mysql",
		"-h", host,
		"-P", port,
		"-u", username,
		fmt.Sprintf("-p%s", password),
		dbName,
	)
	cmd.Stdin = file

	output, err := cmd.CombinedOutput()
	if err != nil {
		logMsg := fmt.Sprintf("MySQL restore error: %v, output: %s", err, string(output))
		db.DB.Exec(db.QueryUpdateRestoreError, logMsg, restoreID)
		utils.LogError(logMsg, err)
		return restoreID, "", err
	}

	db.DB.Exec(db.QueryUpdateRestoreSuccess, string(output), restoreID)
	return restoreID, dbName, nil
}

// cleanBackupFile nettoie le fichier de backup PostgreSQL
func cleanBackupFile(backupPath string) (string, func(), error) {
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read backup file: %v", err)
	}

	contentStr := string(content)
	
	// Remove empty COPY statements
	pattern := `COPY [^;]+ FROM stdin;\s*\\\.`
	re := regexp.MustCompile(pattern)
	cleanedContent := re.ReplaceAllString(contentStr, "")
	
	// Remove problematic lines
	lines := strings.Split(cleanedContent, "\n")
	var finalLines []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip problematic lines
		if strings.HasPrefix(trimmed, "\\") || strings.Contains(line, "more") {
			continue
		}
		
		// Remove owner statements
		if shouldSkipLine(trimmed) {
			continue
		}
		
		finalLines = append(finalLines, line)
	}
	
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "cleaned_backup_*.sql")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	
	finalContent := strings.Join(finalLines, "\n")
	if _, err := tmpFile.WriteString(finalContent); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("failed to write temp file: %v", err)
	}
	tmpFile.Close()
	
	return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }, nil
}

// shouldSkipLine détermine si une ligne doit être ignorée
func shouldSkipLine(line string) bool {
	return (strings.HasPrefix(line, "ALTER TABLE") && strings.Contains(line, "OWNER TO postgres")) ||
	       (strings.HasPrefix(line, "ALTER SEQUENCE") && strings.Contains(line, "OWNER TO postgres"))
}