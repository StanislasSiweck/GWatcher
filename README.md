# Bot Watchman

Ce bot discord est pour surveiller des serveur de jeux (principallement garry's mod) sans besoin de lancer le jeux en question pour voir si le serveur est en ligne, le nombre de joueur, le nom du serveur, ect...
Envoyer les info sur un channel discord et les mettre a jour toute les 1 minutes, ce projet étais pour m'entrainer a faire des bots discord en golang.
Pour ce faire j'ai utiliser la librairie [discordgo](github.com/bwmarrin/discordgo) pour l'api discord, [Source A2S Queries](github.com/rumblefrog/go-a2s) pour les requêtes A2S pour les [Serveur Source](https://developer.valvesoftware.com/wiki/Server_queries).

## Fonctionnement

Pour ajouter un serveur il suffit de faire la commande `/server add ip:IP port:PORT` et le bot va ajouter le serveur a la liste des serveur a surveiller.
Pour retirer un serveur il suffit de faire la commande `/server remove ip:IP port:PORT` et le bot va retirer le serveur de la liste des serveur a surveiller.
Pour changer le channel du message il suffit de faire la commande `/server message ` et le bot va crée un nouveau message dans le channel actuel.
Sous le message il y a des bouton pour changer de page et refresh les info des serveur.


Le bot peut fonctionner en 2 mode:

### DB mode

En mode DB le bot va chercher les serveur a surveiller dans une base de donnée et les mettra a jour toute les 1 minutes pour chaque serveur.

### Local mode

En mode local le bot va chercher les serveur a surveiller en cache (si le bot restart vous perdez la liste des serveur) et les mettra a jour toute les 1 minutes pour chaque serveur.


## Fonctionnalité

- [x] Ajouter un serveur
- [x] Retirer un serveur
- [x] Changer le channel du message
- [x] Mode DB
- [x] Mode Local
- [x] Message avec bouton pour changer de page
- [x] Message avec bouton pour refresh les info des serveur

## Installation

```bash
git clone https://github.com/StanislasSiweck/BotWatchman.git
```

## Usage

Ajouter les variables d'environnement suivante:
```bash
export DISCORD_TOKEN=token

-- Optionnel Local Server
export DISCORD_CHANEL_ID=channel_id
export DISCORD_GUILD_ID=guild_id
-- Optionnel MESSAGE
export DISCORD_MESSAGE_ID=message_id

-- DB
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
docker run -d --name botwatchman -e DISCORD_TOKEN=token -e DISCORD_CHANEL_ID=channel_id -e DISCORD_GUILD_ID=guild_id
```