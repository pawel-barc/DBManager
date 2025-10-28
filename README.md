🧩 Description générale

1. Le projet est développé en Go, React (Vite.js) et PostgreSQL.

2. Le serveur backend utilise Chi, un framework HTTP léger (équivalent à Express.js en Node).

3. Aucune ORM n’a été utilisée — la communication avec la base de données se fait en SQL brut.

🚀 Lancement du projet

1. Installer Go et configurer les variables d’environnement (PATH).

2. Démarrer les serveurs :

Backend → go run main.go

Frontend → npm run dev

3. Toutes les dépendances Go se trouvent dans le fichier go.mod.
   Après un git pull, exécutez simplement :

go mod tidy

pour installer automatiquement les dépendances nécessaires.

--------Informations importantes-------

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
