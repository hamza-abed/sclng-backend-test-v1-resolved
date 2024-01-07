# Canvas for Backend Technical Test at Scalingo

## Instructions

* From this canvas, respond to the project which has been communicated to you by our team
* Feel free to change everything

## Execution

```
docker compose up
```

Application will be then running on port `5000`

## Test

```
$ curl localhost:5000/ping
{ "status": "pong" }
```

## 
@todo: create indexes

## Archi
The API need to have its own database to not depend on github api and its limitations "ratelimit for example"

We need a Worker to store results in a local database 

waitgroup to treat one page of repos at once 

Mutex to protect transactions in db

Save Repo method execute one transaction to respect ACID and rollback the entire inserts if error

L'API est publique donc pas besoin de middleware pour vérifier si l'utilisateur est connecté

Je n'ai pas utilisé d'ORM pour ne pas utiliser des couches qui peuvent ralentir l'exécution de l'API

Dans mes tests je n'ai pas mocké les données pour d'abord pour aller vite et puis cela me permet de tester la connexion à ma DB et l'exécution des queries


## Simpler to go faster

On peut utiliser des contextes pour ajouter du tracing, mais ici je ne l'ai pas fait pour aller vite
On peut rajouter un swagger/openAPI pour faciliter les tests de l'api

