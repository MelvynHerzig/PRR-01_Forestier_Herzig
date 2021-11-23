# PRR-02_Forestier_Herzig
Respository du laboratoire 02 pour le cours PRR

# Étudiants
- Forestier Quentin
- Herzig Melvyn

# Installation

* Cloner le répertoire.
> `$ git clone https://github.com/MelvynHerzig/PRR-01_Forestier_Herzig.git`

* Remplir le fichier de configuration _config.json_ à la racine du projet.
  * debug ( booléen, true/false ): Pour lancer les serveurs en mod debug
  * nbRooms ( nombre, entre 1 et N ): Pour définir le nombre de chambres supportées par l'hôtel
  * nbNights ( nombre, entre 1 et N): Pour définir le nombre de nuits supportées par l'hôtel
  * servers ( ip et port, minimum 1 serveur): Pour définir les adresses et ports du cluster de serveurs de gestion de l'hôtel
```
{
  "debug": true,
  "nbRooms":10,
  "nbNights":10,
  "servers": [
    {
      "ip": "127.0.0.1",
      "port": 3000
    },
    {
      "ip": "127.0.0.1",
      "port": 3001
    },
    {
      "ip": "127.0.0.1",
      "port": 3002
    }
  ]
}
```
> La configuration précédente est un exemple.

* Démarrer le(s) serveur(s). Un argument est nécessaire.
  * Entre 0 et N-1 avec N = nombres de serveurs configurés dans _config.json_

> Depuis le dossier <i>server</i>.
>
> En admettant le fichier de configuration précédent:\
> `$ go run . 1`\
> `$ go run . 2`\
> `$ go run . 0`
>
> L'ordre de démarrage importe peu. Durant cette étape, les serveurs s'inter-connectent. En conséquence, tant que tous ne sont pas allumés et connectés, ils n'acceptent que des connexions ayant une adresse IP source appartenant au fichier de configuration.

* Démarrer le(s) client(s). Un argument est facultatif.
  * Numéro du serveur distant auquel se connecter.

> Depuis le dossier <i>client</i>.
>
> Pour lancer un client qui se connecte à serveur 0.\
> `$ go run . 0`
>
> Pour lancer un client qui se connecte à un serveur aléatoire trouvable dans le fichier _config.json_\
> `$ go run .`

# Utilisation
__Fonctionne__\
Toutes les fonctionnalités de la donnée ont été implémentées avec succès.

__Ne fonctionne pas__\
Rien à notre connaissance.

## Serveur
Une fois le serveur lancé, aucune action supplémentaire est nécessaire.

## Client
Au démarrage, les clients reçoivent la bienvenue du serveur sous cette forme:

` Welcome in the FH Hostel ! Nb rooms: 10, nb nights: 10 ` <br>
`Available commands:` <br>
`-  LOGIN userName` <br>
`-  LOGOUT` <br>
`-  BOOK roomNumber arrivalNight nbNights` <br>
`-  ROOMLIST night` <br>
`-  FREEROOM arrivalNight nbNights` <br>

> Remarquez:
>* le nombre de chambres supportées: de 1 à 10.
>* le nombre de nuits supportées: de 1 à 10.
>* la liste des commandes: LOGIN, LOGOUT, BOOK, ROOMLIST et FREEROOM.

### LOGIN
` LOGIN <userName>`</br>
Première commande à effectuer. Les autres commandes ne fonctionnement pas tant que l'authentification avec un nom d'utilisateur n'a pas été exécutée. Les noms d'utilisateur sont supposés unique. En conséquence, deux utilisateurs avec le même nom ne peuvent pas s'authentifier simultanément.

En cas de succès, l'utilisateur reçoit:
> `Login success` </br>

Sinon il reçoit un message d'erreur avec une explication.

### BOOK
` BOOK roomNumber arrivalNight nbNights` <br>
Cette commande sert à réserver une chambre de numéro <i>roomNumber</i> à partir de la nuit <i>arrivalNight</i> durant un nombre de nuits <i>nbNights</i>. Cette commande est disponible seulement après un <i>LOGIN</i> avec succès.

Si la commande `BOOK 1 2 3` est effectuée avec succès, l'utilisateur reçoit:
>`You successfully booked room  1  for  3  night(s), starting night 2`

Sinon il reçoit un message d'erreur avec une explication.

### ROOMLIST 
` ROOMLIST night` <br>
Cette commande permet de voir l'état des chambres dans l'hôtel pour une nuit donnée <i>night</i>. Cette commande est disponible seulement après un <i>LOGIN</i> avec succès.

Si la commande `ROOMLIST 2` est effectuée avec succès, l'utilisateur reçoit:
>`ROOMLIST 2` <br>
`Room no : state` <br>
`1  :  Self reserved` <br>
`2  :  Free` <br>
`3  :  Free` <br>
`4  :  Occupied` <br>
`5  :  Free` <br>
`6  :  Free` <br>
`7  :  Free` <br>
`8  :  Free` <br>
`9  :  Free` <br>
`10  :  Free` <br>

> Pour la nuit 2, nous appercevons que toutes les chambres sont libres sauf la chambre 1, réservée par l'utilisateur lui même, et la chambre 4, réservée par un autre utilisateur.

Sinon il reçoit un message d'erreur avec une explication.

### FREEROOM
`FREEROOM arrivalNight nbNights` <br>
Cette commande permet de chercher la première chambre libre à partir d'une nuit <i>arrivalNight</i> pendant un nombre de nuits <i>nbNights</i>. Cette commande est disponible seulement après un <i>LOGIN</i> avec succès.

Si la commande FREEROOM 2 1 est effectuée avec succès, l'utilisateur reçoit:
>`Room  2  is free from night  2  during  1  night(s).`

> Effectuée au moment du résultat de l'exécution de <i>ROOMLIST</i> précédent.

Si aucune chambre est disponible, l'utilisateur reçoit:
> `No rooms free from night  2  for  1  night(s).`

Sinon il reçoit un message d'erreur avec une explication.

### LOGOUT
` LOGOUT`</br>
Cette commande permet à un utilisateur de se déconnecter. Cette commande est disponible seulement après un <i>LOGIN</i> réussi.

En cas de succès, l'utilisateur reçoit:
> `Logout success` </br>

Sinon il reçoit un message d'erreur avec une explication.

# Compatibilité
L'application a été testée et validée sous les environnements Windows et Linux.

La compatibilité MacOS n'a pas pu être contrôlée mais devrait être possible. Seules les fonctionnalités de base de Golang ont été utilisées.

# Protocole de communication TCP client-serveur
## Comment le client trouve le serveur (adresses et ports)?
Le client sélectionne un serveur dans le fichier _config.json_ aléatoirement ou manuellement avec son numéro.

## Qui parle en premier ? 
Le serveur parle en premier lorsque le client parvient à se connecter.</br>
Il envoie un message de bienvenue ainsi qu'une liste de commandes.

## Qui ferme la connexion et quand?
Le client ferme la connexion lorsqu'il termine son exécution.

## Qu'est ce qui se passe quand un message est reçu?
### Serveur
Le serveur récupère le premier mot de la requête. Si le premier mot est syntaxiquement:
* Inconnu: Envoie une erreur. (context non concurrent)
* Connu: Tente de former la suite de la requête. (context non concurrent)

Si la requête peut être formée:
* Execute et retourne le résultat. (context concurrent)
* Sinon envoie une erreur. (context non concurrent)


### Client
Le client récupère le premier mot de la réponse. Il détermine si c'est:
* Un résultat et affiche les détails.
* Une erreure et affiche la raison.

## Syntaxe des messages envoyé par le client au serveur
| Utilité | Syntaxe |
|---|----|
| S'identifier à l'hôtel | LOGIN {nom de l'utilisateur} CRLF |
| Réserver une chambre | BOOK {numéro de chambre} {nuit d'arrivée} {nombre de nuits} CRLF  |
| Récupérer la liste des disponnibilités pour une nuit. | ROOMLIST {numéro de nuit} CRLF  |
| Recevoir un numéro de chambre libre pour un nombre de nuits à partir d'une nuit d'arrivée. | FREEROOM {nuit d'arrivée} {nombre de nuit} CRLF  |
| Se déconnecter. | LOGOUT CRLF |

## Syntaxe des messages envoyé par le serveur au client
| Utilité | Syntaxe |
|---|----|
| Réponse positive à LOGIN | RESULT_LOGIN CRLF |
| Réponse positive à BOOK | RESULT_BOOK {numéro de chambre} {nuit d'arrivée} {nombre de nuits} CRLF  |
| Réponse positive à ROOMLIST | RESULT_ROOMLIST {état chambre1}, {état chambre2} ...  CRLF  |
| Réponse positive à FREEROOM | RESULT_FREEROOM {no chambre libre ou 0} CRLF |
| Réponse positive à LOGOUT | RESULT_LOGOUT {numéro de chambre ou 0} {nuit d'arrivée} {nombre de nuits} CRLF |
| Erreur | ERROR {message} CRLF |

## Exemple d'une conversation entre client et serveur tcp
Server : </br>
` Welcome in the FH Hotel ! Nb rooms: 10, nb nights: 10 ` <br>
`Available commands:` <br>
`-  LOGIN userName` <br>
`-  LOGOUT` <br>
`-  BOOK roomNumber arrivalNight nbNights` <br>
`-  ROOMLIST night` <br>
`-  FREEROOM arrivalNight nbNights CRLF` <br>
Client : <br> 
`LOGIN John CRLF`\
Server : <br>
`RESULT_LOGIN CRLF`\
Client : <br> 
`ROOMLIST 1 CRLF`\
Server :\
`RESULT_ROOMLIST Free, Occupied, Free, Self reserved CRLF`\
Client :\
`FREEROOM 1 2 CRLF`\
Server :\
`RESULT_FREEROOM 1 1 2 CRLF`\
Client :\
`BOOK 1 1 2 CRLF`\
Server :\
`RESULT_BOOK OK CRLF`\
Client :\
`BOOK 1 1 2 CRLF`\
Server :\
`ERROR room already booked CRLF`\
Client : <br> 
`LOGOUT CRLF`\
Server : <br>
`RESULT_LOGOUT CRLF`

# Protocole de communication TCP serveur-serveur
Au démarrage, les serveurs doivent s'attendre pour ouvrir des connexions avant de commencer à traiter les clients.

## Comment serveur trouve un autre serveur (adresses et ports)?
Le serveur interroge le fichier _config.json_.

## Qui parle en premier ? 
Les serveurs de numéros M se connectent aux serveur de numéro 0 - N avec N < M. Les serveurs de numéro M transmettent leur numéro M de serveur.

## Qui ferme la connexion et quand?
Les serveurs ne ferment la connexion que si tout le cluster doit s'arrêter. Théoriquement, jamais.

## Qu'est ce qui se passe quand un message est reçu?
### Serveur
Le serveur récupère le premier mot de la requête et vérifie si il correspond à:
* Premier mot d'une requête client-serveur: il applique immédiatement la requête sur ses données.
* Confirmation de réplication: il signal au travers d'un canal spécialisé qu'une confirmation de réplication a été reçue.
* Premier mot d'un message de l'algorithme de lamport: il forme le message et le transmet au mutex.

## Syntaxe des messages de réplication
### Requête
| Utilité | Syntaxe |
|---|----|
| S'identifier à l'hôtel | LOGIN {nom de l'utilisateur} CRLF |
| Réserver une chambre | BOOK {numéro de chambre} {nuit d'arrivée} {nombre de nuits} {nom de l'utilisateur mendant} CRLF  |
| Récupérer la liste des disponnibilités pour une nuit. | ROOMLIST {numéro de nuit} {nom de l'utilisateur mendant} CRLF  |
| Recevoir un numéro de chambre libre pour un nombre de nuits à partir d'une nuit d'arrivée. | FREEROOM {nuit d'arrivée} {nombre de nuit} {nom de l'utilisateur mendant} CRLF  |
| Se déconnecter. | LOGOUT {nom de l'utilisateur mendant} CRLF |

> À noter: Seules les requêtes de LOGIN, BOOK et LOGOUT sont répliquées. De plus, les requêtes qui ne peuvent pas être traitées localement ne seront pas répliquées.

### Réponse
| Utilité | Syntaxe |
|---|----|
| Réponse positive à une réplication | OK CRLF |

## Syntaxe des messages de Lamport
| Utilité | Syntaxe |
|---|----|
| Acknowledgement | ACK {server local timestamp} {server local number} CRLF |
| Mutex demand | REQ {server local timestamp} {server local number} CRLF |
| Mutex release | REL {server local timestamp} {server local number} CRLF |

> La version optimisée a été implémentée: un ack est envoyé par un serveur i seulement si son dernier message envoyé n'est pas un req.

## Exemple d'une conversation entre 2 serveurs Server0 et Server1

_Synchronisation (une fois au démarrage)_ 

Server1 : <br>
`1`

_Demande de mutex_

Server0 : <br>
`REQ 1 0`\
Server1 : <br> 
`ACK 2 1`

_Réplication_

Server0 : <br>
`BOOK 1 1 2 Pierre`\
Server1 : <br> 
`OK`

_Relâchement de mutex_

Server0 : <br>
`REL 5 0`

# Tests automatique
Un program de tests a été mis en place. Le programme de tests se comporte comme un seul client. Toutes les requêtes, qu'un client peut émettre, sont testées avec leurs retours
positifs (pas d’erreurs) et négatifs (paramètres incorrects, indisponibilités, …). 

Pour que le programme de tests fonctionne, il faut lancer trois serveurs conformément à la section **Installation** tout en utilisant le fichier de configuration _config.json_ qui y est présenté.

Pour lancer les tests, depuis le dossier <i>server</i>, exécuter la commande: `$go test -v`

# Debug de la concurrence
Comme présenté dans la rubrique "Installation", le serveur peut être lancé en mode debug grâce au paramètre "-debug". Ce mode de fonctionnement affiche les événements dans la console. De plus, à chaque fois que deux utilisateurs sont connectés (login) avec succès, la goroutine qui gère les ressources partagées se met en pause pendant 20 secondes dans le but de laisser suffisament de temps pour créer une situation de concurrence. 

## Distinctions
Il existe deux types de log:
* RISK: log une requête effectuée dans la zone partagée (goroutine gérant l'hôtel, hostelManager).
* SAFE: log une réception ou un envoi depuis/vers le client (goroutine gérant la communication spécifique à chaque client, clientHandler). Ce type de log peut apparaître au milieu d'un passage en zone concurrente sans problème.

Théoriquement, pour une bonne gestion de la concurrence, les logs qui indiquent un passage (entrée puis sortie) en zone partagé ne doivent pas se chevaucher.

__Correct__ \
`1 DEBUG >>  RISK)  --------- Enter shared zone ---------`\
`2 DEBUG >>  SAFE) From 127.0.0.1:5155 : BOOK 1 1 1`\
`3 DEBUG >>  RISK) From 127.0.0.1:5155 BOOK room 1 from night1 for 1 night(s) HANDLING`\
`4 DEBUG >>  RISK) From 127.0.0.1:5155 BOOK room 1 from night1 for 1 night(s) SUCCESS`\
`5 DEBUG >>  RISK) --------- Leave shared zone ---------`


> Ligne 1, nous voyons que la goroutine gérant les accès concurrents est entrée en zone partagée, prête à traiter la prochaine demande. Ensuite en ligne 2, par le préfix SAFE nous voyons que le client 127.0.0.1:5155 a envoyé la requête BOOK 1 1 1 et que sa goroutine dédiée a reçu sa requête. Ligne 3, nous voyons que la requête a été transmise et que la goroutine qui gère l'hôtel traite la demande.
Finalement lignes 4 et 5, le traitement est terminé et la zone concurrente est quittée.
Cet exemple montre un cas d'exécution correct. Il n'y a qu'une exécution au sein du même passage en zone critique. De plus, aucune nouvelle entrée en section critique est effectuée tant que la première n'est pas sortie.

__Faux__ \
`1  DEBUG >>  RISK)  --------- Enter shared zone --------- `\
`2  DEBUG >>  RISK)  --------- Enter shared zone --------- `\
`3  DEBUG >>  SAFE) From 127.0.0.1:5155 : BOOK 1 1 1 `\
`4  DEBUG >>  SAFE) From 127.0.0.1:5156 : BOOK 1 1 1 `\
`5  DEBUG >>  RISK) From 127.0.0.1:5155 BOOK room 1 from night1 for 1 night(s) HANDLING`\
`6  DEBUG >>  RISK) From 127.0.0.1:5155 BOOK room 1 from night1 for 1 night(s) SUCCESS`\
`7  DEBUG >>  RISK) From 127.0.0.1:5156 BOOK room 1 from night1 for 1 night(s) HANDLING`\
`8  DEBUG >>  RISK) From 127.0.0.1:5156 BOOK room 1 from night1 for 1 night(s) SUCCESS`\
`9  DEBUG >>  RISK) --------- Leave shared zone ---------` \
`10 DEBUG >>  RISK) --------- Leave shared zone ---------`

> Contrairement à l'exemple précédent, celui-ci montre un cas qui ne doit pas arriver. Nous pouvons voir que deux traitements critiques ont été executés en même temps. En effet, une entrée en zone partagée a été effectuée alors que la précédente entrée n'a pas terminé son traitement. Les lignes 2,7,8 et 10 devraient être exécutées après les lignes 1,5,6 et 9.

## Vérification manuelle

Pour vérifier manuellement la concurrence suivez les étapes suivantes:

* Démarrer le serveur en mode debug\
`$go run . 10 10 -debug`

* Démarrer deux clients (A et B) \
`$go run . localhost` (A) \
`$go run . localhost` (B)

* Identifier les clients
  * A \
  `LOGIN melvyn`
  * B \
  `LOGIN quentin`

> A ce moment la goroutine du serveur qui gère l'exécution des requêtes s'endore pendant 20 secondes. Envoyer deux requêtes pendant ces 20 secondes.

* Envoyer la même requête depuis les deux clients\
  * A \
  `BOOK 1 1 1`
  * B \
  `BOOK 1 1 1`

* Attendre la fin des 20 secondes.
> Normalement, un client reçoit une validation et l'autre une erreur.

### Analyse d'une sortie des logs.

`1  DEBUG >>  RISK) --------- Enter shared zone --------- ` \
`2  DEBUG >>  SAFE) To 127.0.0.1:13120 : WELCOME {...} ` \
`3  DEBUG >>  SAFE) To 127.0.0.1:13124 : WELCOME {...} ` \
`4  DEBUG >>  SAFE) From 127.0.0.1:13120 : LOGIN melvyn ` \
`5  DEBUG >>  RISK) From 127.0.0.1:13120 login as melvyn HANDLING ` \
`6  DEBUG >>  RISK) From 127.0.0.1:13120 login as melvyn SUCCESS ` \
`7  DEBUG >>  RISK) --------- Leave shared zone --------- ` \
`8  DEBUG >>  RISK) --------- Enter shared zone --------- ` \
`9  DEBUG >>  SAFE) To 127.0.0.1:13120 : RESULT_LOGIN ` \
`10 DEBUG >>  SAFE) From 127.0.0.1:13124 : LOGIN quentin ` \
`11 DEBUG >>  RISK) From 127.0.0.1:13124 login as quentin HANDLING ` \
`12 DEBUG >>  RISK) Server request handler suspended. Resume in 20s.` \
`13 DEBUG >>  SAFE) To 127.0.0.1:13124 : RESULT_LOGIN ` \
`14 DEBUG >>  SAFE) From 127.0.0.1:13120 : BOOK 1 1 1 ` \
`15 DEBUG >>  SAFE) From 127.0.0.1:13124 : BOOK 1 1 1 ` \
`16 DEBUG >>  RISK) Server request handler resumed. ` \
`17 DEBUG >>  RISK) From 127.0.0.1:13124 login as quentin SUCCESS ` \
`18 DEBUG >>  RISK) --------- Leave shared zone --------- ` \
`19 DEBUG >>  RISK) --------- Enter shared zone --------- ` \
`20 DEBUG >>  RISK) From 127.0.0.1:13120 BOOK room 1 from night1 for 1 night(s) HANDLING ` \
`21 DEBUG >>  RISK) From 127.0.0.1:13120 BOOK room 1 from night1 for 1 night(s) SUCCESS ` \
`22 DEBUG >>  RISK) --------- Leave shared zone --------- ` \
`23 DEBUG >>  RISK) --------- Enter shared zone --------- ` \
`24 DEBUG >>  RISK) From 127.0.0.1:13124 BOOK room 1 from night1 for 1 night(s) HANDLING ` \
`25 DEBUG >>  RISK) From 127.0.0.1:13124 BOOK room 1 from night1 for 1 night(s) ERROR ` \
`26 DEBUG >>  RISK) --------- Leave shared zone --------- ` \
`27 DEBUG >>  RISK) --------- Enter shared zone --------- ` \
`28 DEBUG >>  SAFE) To 127.0.0.1:13120 : RESULT_BOOK 1 1 1 ` \
`29 DEBUG >>  SAFE) To 127.0.0.1:13124 : ERROR room already booked ` \


> __Lignes 1-18)__ Ces lignes montrent l'enchaînement des appels lorsque les deux clients se connectent au serveur et que les deux utilisateurs s'authentifient. Comme annoncé, lorsque deux utilisateurs sont authentifiés, la go routine de traîtement des requêtes se met en pause, c'est ce que nous voyons en ligne 12. Ensuite, les deux clients envoient la même requête. Les lignes 13 et 14 indiquent que les goroutines qui communiquent avec les clients ont bien reçu leur requête. La ligne 16 indique la reprise de la goroutine de traitement. Finalement, les lignes 17 et 18 terminent le traitement du login du client 127.0.0.1:13124. 

>__Lignes 19-22)__ La requête BOOK du client 127.0.0.1:13120 est traitée.

>__Lignes 23-26)__ La requête BOOK du client From 127.0.0.1:13124 est traitée. Ligne 25, confirmation de l'échec. Cet échec est dû au fait que la chambre a déjà été réservée par le client précédent.

>__Ligne 27)__ Attente d'un nouveau traitement concurrent.

>__Ligne 28-29)__ Envoi des réponses aux clients.

Comme nous pouvons le voir, aucune entrée en zone partagée n'est effectuée tant que l'entrée précédente n'est pas terminée. Chaque entrée partagée contient le traitement d'une seule requête. En conclusion, la gestion des accès concurrents est correcte.

## go race
L'application a passé le test du go race.

La procédure précédente (Debug de la concurrence -> vérification manuelle) a été exécutée en démarrant le serveur avec l'argument <i>-race</i>. 
> `$ go run -race . 10 10 -debug`

Aucune concurrence n'a été détectée.
