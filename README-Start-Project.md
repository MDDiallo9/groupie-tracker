# 🎸 Groupie Tracker

Un projet web interactif qui permet de visualiser les groupes de musique, leurs membres, leurs dates de concerts et leurs localisations, en s'appuyant exclusivement sur une API publique.

---

## 📂 GitHub de l'exercice

👉 [Sujet sur GitHub](https://github.com/01-edu/public/tree/master/subjects/groupie-tracker)

---

## 🧠 Premières réflexions

- Les données seront exclusivement récupérées via l'API, sans sauvegarde locale sur le serveur ou l'ordinateur.
- L'embellissement du site sera traité en dernier, mais la structure de présentation et l'expérience utilisateur doivent être pensées dès le début pour construire un backend solide.

---

## 🗺️ Carte des API

📍 **URL racine de l'API** :  
`https://groupietrackers.herokuapp.com/api`

### 🎤 Liste des artistes  
`https://groupietrackers.herokuapp.com/api/artists`

- ID du groupe : `1` à `52`
- Image du groupe : `https://groupietrackers.herokuapp.com/api/images/{nom_du_groupe attaché en minuscule}.jpeg`
- Nom du groupe
- Membres du groupe (min. 1, max inconnu)
- Date de création : `AAAA`
- Date du premier album : `JJ-MM-AAAA`
- API des localisations : `https://groupietrackers.herokuapp.com/api/locations/{ID}`
- API des dates : `https://groupietrackers.herokuapp.com/api/dates/{ID}`
- API des concerts (dates + localisations) : `https://groupietrackers.herokuapp.com/api/relation/{ID}`

---

### 📅 Liste des concerts avec localisations et dates  
`https://groupietrackers.herokuapp.com/api/relation`

- ID du groupe : `1` à `52`
- Localisations et dates :
  - Format : `"NomVille ou NomVille_Composé" - "NomPays ou NomPays_Composé" "JJ-MM-AAAA"`
  - ⚠️ Plusieurs dates possibles pour une même localisation

---

### 📆 Liste des dates de concerts  
`https://groupietrackers.herokuapp.com/api/dates`

- ID du groupe : `1` à `52`
- Dates :
  - Format : `"JJ-MM-AAAA"` ou `"*JJ-MM-AAAA"`
  - `*` indique plusieurs dates dans une même ville

---

### 🌍 Liste des localisations des concerts  
`https://groupietrackers.herokuapp.com/api/locations`

- ID du groupe : `1` à `52`
- Localisations :
  - Format : `"NomVille ou NomVille_Composé" - "NomPays ou NomPays_Composé"`
  - ⚠️ Une ville n’est listée qu’une seule fois, peu importe le nombre de dates
- Lien vers les dates : `https://groupietrackers.herokuapp.com/api/dates/{ID}`

---

## 🛠️ Fonctionnalités à intégrer

### 🔍 Barre de recherche  
[Documentation](https://github.com/01-edu/public/tree/master/subjects/groupie-tracker/search-bar)

### 📌 Géolocalisation  
[Documentation](https://github.com/01-edu/public/tree/master/subjects/groupie-tracker/geolocalization)

### 🧰 Filtres  
[Documentation](https://github.com/01-edu/public/tree/master/subjects/groupie-tracker/filters)

- Date de création du groupe
- Date du premier album
- Nombre de membres
- Localisation des concerts

#### 🧩 Options UI à intégrer

- Slider : [Exemples sur Dribbble](https://dribbble.com/search/filter-slider)
- Checkbox : [Documentation MDN](https://developer.mozilla.org/fr/docs/Web/HTML/Reference/Elements/input/checkbox)
