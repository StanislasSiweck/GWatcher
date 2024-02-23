# Utiliser une image de base qui contient Go
FROM golang:1.20-alpine

# Définir le répertoire de travail à l'intérieur du conteneur
WORKDIR /app

# Copier le contenu du répertoire actuel (où se trouve votre code Go) dans le conteneur
COPY . .

# Compiler le code Go
RUN go build -o main .

# Exposer le port sur lequel votre application Go écoute
EXPOSE 8080

# Commande pour exécuter votre application Go une fois que le conteneur est démarré
CMD ["./main"]