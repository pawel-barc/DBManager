🧩 Description générale

⚠️ IMPORTANT
La branche feature/configuration contient l’architecture prsque complète du projet(manque de styles et tests dans le front).
Après avoir cloné le projet, basculez sur cette branche — elle contient également le script de migration pour la création de la base de données.

1. Le projet est développé en Go, React (Vite.js) et PostgreSQL.

2. Le serveur backend utilise Chi, un framework HTTP léger (équivalent à Express.js en Node).

3. Aucune ORM n’a été utilisée — la communication avec la base de données se fait en SQL brut.

🚀 Lancement du projet

1.  Installer Go et configurer les variables d’environnement (PATH).

2.  Préparation de la base de données

Ouvrez pgAdmin4 et connectez-vous à votre serveur PostgreSQL.

Créez la base de données : CREATE DATABASE safebase;

2.  Ouvrez le fichier main.go (dans le dossier backend) et décommentez temporairement la section suivante pour exécuter
    la migration initiale :
    // Migration: création de la base des données, Commentez après avoir utiliser
    // sqlBytes, err := os.ReadFile("db/migrations/001_init.sql")
    // if err != nil {
    // log.Fatal("Impossible de lire le fichier de migration:", err)
    // }

        // _, err = db.DB.Exec(string(sqlBytes))
        // if err != nil {
        // 	log.Fatal("Erreur lors de l'exécution de la migration", err)
        // }

        // log.Println("Migration exécutée avec succès")

3.  Démarrer les serveurs :

Backend → go run main.go Si dans le terminal: ( Connexion à la base des donnès réussie, Migration exécutée avec succès et Demarrage du serveur sur :8080...) c'est ok.
Si vous rencontrez une erreur de connexion à la base, vérifiez les identifiants dans le fichier db/db.go
connStr := "user=postgres password=root dbname=safebase sslmode=disable" db/db.go

ensuite:
Frontend → npm i package et
Frontend → npm run dev

4. Toutes les dépendances Go se trouvent dans le fichier go.mod.
   executer s'il faut :

go mod tidy

1. La gestion des réponses d’erreur et de succès se trouve dans utils/response.go.
   → Utilisez ces fonctions dans tout le projet afin d’unifier le format des réponses et d’éviter les répétitions de code.

2. Importations et portée :
   Les fonctions dont le nom commence par une majuscule sont exportées et peuvent être utilisées depuis d’autres packages.
   Exemple :

import "safebase/controllers"

controllers.Register()

3. Base de données :
   Les identifiants actuels sont configurés avec mes propres paramètres.
   On peut les mettre aussi dans .env et ca pourras marcher

4. Le projet utilise le fichier .env pour stocker les variables d’environnement.
   Un exemple est fourni dans .env.example.

5. Les routes sont divisées en publiques et privées.
   → Le middleware AuthMiddleware décide si un utilisateur a accès à une route protégée.

6. Fonctionnement du middleware d’authentification :
   Il lit le token dans le cookie access_token, vérifie sa validité et en extrait l’ID de l’utilisateur, qui est ensuite stocké dans le contexte de la requête.

serID := int(userIDFloat)

// Stockage de l'ID utilisateur dans le contexte de la requête
ctx := context.WithValue(r.Context(), UserIDKey, userID)

7. Récupération de l’ID utilisateur depuis le contexte :
   Exemple :
   import (
   "net/http"
   "safebase/db"
   "safebase/middleware"
   "safebase/utils"
   )

userID := r.Context().Value(middleware.UserIDKey).(int)

8. Tokens :

Access token → valide pendant 15 minutes

Refresh token → stocké dans la base de données, valide pendant 7 jours

Le frontend est copier/coller et un peut modifier, et est a corriger s'il faut
Chez moi tout marche comme il faut jusquau la

===================IMPORTANT===========================
pgAdmin 4 Querrytool rights for user!!!!

1.  psql -U postgres
2.  CREATE DATABASE exampledb;
3.  \c exampledb
4.  -- users
    CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- payments
CREATE TABLE payments (
id SERIAL PRIMARY KEY,
user_id INT REFERENCES users(id) ON DELETE CASCADE,
amount NUMERIC(10,2) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

5.  CREATE USER example_user WITH PASSWORD 'securepass';
6.  ALTER TABLE users OWNER TO example_user;
    ALTER TABLE payments OWNER TO example_user;
7.  GRANT ALL PRIVILEGES ON DATABASE exampledb TO example_user;

GRANT USAGE ON SCHEMA public TO example_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO example_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO example_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO example_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO example_user;

=============IMPORTANT CHANGEMENT DANS LA BASE DE DONNES============
AJOUTER dans le querytool pgadmin:
DATABASE CAHANGEMENT :

ALTER TABLE scheduled_tasks
ADD COLUMN user_id INT NOT NULL REFERENCES users(id);

ALTER TABLE databases
    DROP CONSTRAINT IF EXISTS databases_user_id_fkey,
    ADD CONSTRAINT fk_databases_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE;

ALTER TABLE backups
    DROP CONSTRAINT IF EXISTS backups_database_id_fkey,
    ADD CONSTRAINT fk_backups_database
        FOREIGN KEY (database_id)
        REFERENCES databases(id)
        ON DELETE CASCADE;

ALTER TABLE restores
    DROP CONSTRAINT IF EXISTS restores_backup_id_fkey,
    ADD CONSTRAINT fk_restores_backup
        FOREIGN KEY (backup_id)
        REFERENCES backups(id)
        ON DELETE CASCADE;


ALTER TABLE restores
    DROP CONSTRAINT IF EXISTS restores_user_id_fkey,
    ADD CONSTRAINT fk_restores_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE;

ALTER TABLE scheduled_tasks
    DROP CONSTRAINT IF EXISTS scheduled_tasks_database_id_fkey,
    ADD CONSTRAINT fk_scheduled_tasks_database
        FOREIGN KEY (database_id)
        REFERENCES databases(id)
        ON DELETE CASCADE;

ALTER TABLE alerts
    DROP CONSTRAINT IF EXISTS alerts_user_id_fkey,
    ADD CONSTRAINT fk_alerts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE;
