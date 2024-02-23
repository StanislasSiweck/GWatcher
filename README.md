# Bot Watchman

Ce bot discord est pour surveiller des serveur de jeux (principallement garry's mod) sans besoin de lancer le jeux en question pour voir si le serveur est en ligne, le nombre de joueur, le nom du serveur, ect...
Envoyer les info sur un channel discord et les mettre a jour toute les 1 minutes, ce projet étais pour m'entrainer a faire des bots discord en golang.
Pour ce faire j'ai utiliser la librairie [discordgo](github.com/bwmarrin/discordgo) pour l'api discord, [Source A2S Queries](github.com/rumblefrog/go-a2s) pour les requêtes A2S pour les [Serveur Source](https://developer.valvesoftware.com/wiki/Server_queries).
Les commentaires ont été générés de manière automatique.

## Installation

```bash
git clone https://github.com/StanislasSiweck/BotWatchman.git
```

## Usage

Ajouter les variables d'environnement suivante (la base de donnée est optionnel):
```bash
export DISCORD_TOKEN=token
export DISCORD_CHANEL_ID=channel_id

export DB_HOST=host
export DB_PORT=port
export DB_USERNAME=username
export DB_PASSWORD=password
export DB_DATABASE=name
```

Lancer le bot:
```bash
go run main.go
```

## Docker

Récupérer l'image:
```bash
docker pull docker.pkg.github.com/stanislassiweck/botwatchman/watchman:latest
```

Lancer le container:
```bash
docker run -d --name botwatchman -e DISCORD_TOKEN=token -e DISCORD_CHANEL_ID=channel_id
```