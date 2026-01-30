# ------ DOCKER POUR LE DÉVELOPPEMENT BACKEND GO ---------#

# Utilise une version stable et supportée (trouvée dans le doc GO)
FROM golang:1.25.5-alpine AS base

# Définir le répertoire de travail principal du conteneur
WORKDIR /app

# Copier les fichiers go.mod et go.sum(dépendences) dans le conteneur
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier tous les fichiers et dossier du backend dans le conteneur
COPY . .

# Exposer le port utilisé par l'API GO
EXPOSE 8080

# Démarrer le serveur de développement du backend
CMD ["go", "run", "main.go"]
