# Projet-Groupie-Tracker

Une application web dédiée à l'exploration des données spatiales de la NASA. Découvrez et explorez des images, des catégories et des ressources liées à l'espace et aux missions spatiales.



## 🚀 Fonctionnalités principales

- Navigation par catégories de contenu spatial
- Recherche avancée de ressources et d'images
- Visualisation détaillée des collections
- Accès aux informations sur les missions et projets de la NASA
- Sauvegarde de favoris pour un accès rapide

## 📋 Prérequis

- Go (version 1.16 ou supérieure)
- Git

## 🛠️ Installation et démarrage

1. Clonez le dépôt GitHub :
   ```bash
   git clone https://github.com/votre-compte/projet-nasa.git
   ```

2. Accédez au répertoire du projet :
   ```bash
   cd projet-nasa
   ```

3. Installez les dépendances :
   ```bash
   go mod download
   ```

4. Démarrez le serveur :
   ```bash
   go run main.go
   ```

5. Accédez à l'application dans votre navigateur :
   ```
   http://localhost:8080
   ```

## 🌐 Routes implémentées

| Route | Méthode | Description |
|-------|---------|-------------|
| `/` | GET | Page d'accueil |
| `/categories` | GET | Liste de toutes les catégories |
| `/category/:id` | GET | Détails d'une catégorie spécifique |
| `/collection/:id` | GET | Affichage d'une collection spécifique |
| `/resource/:id` | GET | Détails d'une ressource spécifique |
| `/search` | GET | Recherche de ressources |
| `/about` | GET | Page à propos du projet |
| `/favoris` | GET | Accès aux éléments favoris |

## 🛰️ API NASA utilisée

Ce projet utilise l'API publique de la NASA pour récupérer des données authentiques sur l'espace et les missions spatiales.

- **Base URL**: `https://api.nasa.gov/` `https://images-api.nasa.gov/` `https://gibs.earthdata.nasa.gov/`
- **Endpoints utilisés**:
  - `/planetary/apod` - Image astronomique du jour
  - `/mars-photos/api/v1/rovers` - Photos des rovers martiens
  -`/wmts/epsg4326/best/MODIS_Terra_CorrectedReflectance_TrueColor/default/{date}/250m/{zoom}/{y}/{x}.jpg` - Images satellite de la Terre
  -`https://api.rss2json.com/v1/api.json` - Conversion du flux RSS des actualités NASA
  -`https://api.mapbox.com/styles/v1/mapbox/satellite-v9/static/{lon},{lat},{zoom},0/600x400` - Images satellite Mapbox
  -`https://maps.googleapis.com/maps/api/staticmap` - Solution de secours pour images satellite
## 👥 Structure du projet

L'application suit une architecture MVC (Model-View-Controller) avec les composants suivants :

- **Models**: Structures de données et interaction avec l'API NASA
- **Views**: Templates HTML et styles CSS
- **Controllers**: Logique de traitement des requêtes et préparation des données

## 📅 Méthodologie de développement

Le projet a été développé en utilisant une approche agile avec des sprints hebdomadaires et une organisation des tâches basée sur la matrice MoSCoW (Must have, Should have, Could have, Won't have).

## 🧪 Tests

Pour exécuter les tests unitaires :

```bash
go test ./...
```

## 📚 Documentation

Pour plus d'informations sur l'utilisation de l'API NASA, consultez la [documentation officielle](https://api.nasa.gov/).

## 🤝 Contribuer

1. Forkez le projet
2. Créez votre branche de fonctionnalité (`git checkout -b feature/nouvelle-fonctionnalite`)
3. Committez vos changements (`git commit -am 'Ajout d'une nouvelle fonctionnalité'`)
4. Poussez vers la branche (`git push origin feature/nouvelle-fonctionnalite`)
5. Créez une nouvelle Pull Request

## 📝 Licence

Ce projet est sous licence MIT. Voir le fichier LICENSE pour plus de détails.

## 🙏 Remerciements

- La NASA pour l'accès à ses API publiques et ses ressources
- L'équipe de développement pour son travail acharné et sa passion pour l'exploration spatiale
