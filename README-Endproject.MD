# 🎸 Groupie Tracker

## 📌 Présentation

L'objectif était de concevoir un site web interactif en exploitant une API fournie, sans sauvegarde locale, afin de visualiser les données relatives à des groupes de musique : artistes, dates de concerts, lieux et relations entre ces éléments.

Nous avons **implémentés toutes les fonctionnalités demandées** dans le sujet, en respectant les contraintes techniques et les bonnes pratiques de développement.

**Lien GITHUB de l'exercice** : https://github.com/01-edu/public/tree/787074ad7908b3ad598cdd07720fb49863245a41/subjects/groupie-tracker

## 🚀 Fonctionnalités

- **Affichage complet des données de l'API** :
  - Informations sur les artistes : nom, image, année de début, premier album, membres.
  - Dates et lieux des concerts.
  - Relations entre artistes, dates et lieux.

- **Visualisation dynamique des données** :
  - Utilisation de cartes de géolocalisations, structures, listes et pages dédiées pour une navigation fluide.
  - Interface utilisateur claire et responsive.

- **Système d’événements côté client** :
  - Mise en place d’interactions qui déclenchent des appels serveur (client-server).
  - Exemple : suggestions dynamiques, filtres par date ou localisation.

- **Backend en Go** :
  - Serveur robuste, sans crash.
  - Gestion des erreurs et des cas limites.
  - Respect des bonnes pratiques Go, commentaire et structuration du code.

## 🛠️ Technologies utilisées

- **Langage backend** : Go (Golang)
- **Frontend** : HTML, CSS, JavaScript
- **API utilisée** : [Groupie Tracker API](https://groupietrackers.herokuapp.com/api)

## 📦 Packages autorisés

Nous avons uniquement utilisé les **packages standards de Go**, conformément aux consignes du projet.

## ✅ Statut du projet

✔️ **Terminé** — Toutes les options demandées ont été implémentées.