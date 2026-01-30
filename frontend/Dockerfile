# ------ DOCKER POUR LE DÉVELOPPEMENT EN REACT ---------#

# Utilise la dernière version LTS de node.js (LONG TERM SUPPORT)
FROM node:24-alpine AS build

# Définir le répertoire de travail principal du conteneur
WORKDIR /app

# Copier les fichiers package.json et package-lock.json dans le conteneur
COPY package*.json ./ 

# Installer les dépendances
RUN npm install

# Copier tous les fichiers et dossier du frontend dans le conteneur
COPY . .

# Exposer le port utilisé par le frontend React 
EXPOSE 5173

# Démarrer le serveur de développement du frontend
CMD ["npm", "run", "dev"]