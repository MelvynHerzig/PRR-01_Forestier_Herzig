# PRR-03_Forestier_Herzig
Respository du laboratoire 03 pour le cours PRR

# Étudiants
- Forestier Quentin
- Herzig Melvyn

# Installation

* Cloner le répertoire.
> `$ git clone https://github.com/MelvynHerzig/PRR-01_Forestier_Herzig.git`

* Remplir le fichier de configuration _config.json_ à la racine du projet.
  * debug ( booléen, true/false ): Pour lancer les serveurs en mod debug
  * nbRooms ( nombre, entre 1 et N ): Pour définir le nombre de chambres supportées par l'hôtel
  * nbNights ( nombre, entre 1 et M ): Pour définir le nombre de nuits supportées par l'hôtel
  * servers ( ip, port et numéro du parent [0, Nb serveurs], minimum 1 serveur): Pour définir les adresses, les ports et l'arborescence du cluster de serveurs de gestion de l'hôtel
```
{
  "debug": true,
  "nbRooms": 10,
  "nbNights": 10,
  "servers": [
    {
      "ip": "127.0.0.1",
      "port": 3000,
      "parent": 2
    },
    {
      "ip": "127.0.0.1",
      "port": 3001,
      "parent": 1
    },
    {
      "ip": "127.0.0.1",
      "port": 3002,
      "parent": 1
    },
    {
      "ip": "127.0.0.1",
      "port": 3003,
      "parent": 2
    },
    {
      "ip": "127.0.0.1",
      "port": 3004,
      "parent": 1
    }
  ]
}
```
> La configuration précédente est un exemple avec cinq serveurs.\
> Sachant que le premier serveur dans la liste est le serveur 0 et le dernier, le serveur 4.\
> L'arborescence serait la suivante:\
>![Sans titre](https://user-images.githubusercontent.com/34660483/146408671-b4042b0a-5ddf-4226-a552-fef85f3ba00c.png)

* Démarrer le(s) serveur(s). Un argument est nécessaire.
  * Entre 0 et N-1 avec N = nombres de serveurs configurés dans _config.json_

> Depuis le dossier <i>server</i>.
>
> En admettant le fichier de configuration précédent:\
> `$ go run . 0`\
> `$ go run . 2`\
> `$ go run . 4`\ 
> `$ go run . 1`\
> `$ go run . 3`
>
> L'ordre de démarrage 'est pas important. Durant cette étape, les serveurs s'inter-connectent. En conséquence, tant que tous ne sont pas allumés et connectés, ils n'acceptent que des connexions ayant une adresse IP source appartenant au fichier de configuration.

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

## Qui parle en premier (attente de la mise en ligne entre serveurs) ? 
Au démarrage, un serveur commence par attendre que ses enfants (si il en a) se connectent. Les enfants parlent en premier en envoyant leur numéro de serveur entre 0 et N-1.
Quand tous les enfants d'un serveur sont connectés et qu'ils ont envoyé leur numéro, le serveur se connectent à son parent.
Lorsque tous les enfants direct du noeud racine sont connectés et ont envoyés leur numéro, la racine envoie un ordre de démarrage ("GO") à ses enfants qui est propagés aux enfants des enfants et ainsi de suite.

## Qui ferme la connexion et quand?
Les serveurs ferment la connexion que si tout le cluster doit s'arrêter. Théoriquement, jamais.

## Qu'est ce qui se passe quand un message est reçu (après mise en ligne de tous les serveur)?
### Serveur
Le serveur récupère le premier mot de la requête et vérifie si il correspond à:
* Premier mot d'une requête client-serveur: il applique immédiatement la requête sur ses données.
* Confirmation de réplication: il signal au travers d'un canal spécialisé qu'une confirmation de réplication a été reçue.
* Premier mot d'un message de l'algorithme de Raymond: il forme le message et le transmet au mutex.

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

## Syntaxe des messages de Raymond
| Utilité | Syntaxe |
|---|----|
| Request | req {sender server id} CRLF |
| Token | token {sender server id} CRLF |

> La version optimisée a été implémentée: un ack est envoyé par un serveur i seulement si son dernier message envoyé n'est pas un req.

## Exemple d'une conversation entre 2 serveurs Server 2 et Server 1 (en fonction de la configuration initiale)

_Synchronisation (une fois au démarrage)_ 

// Server 1 est déjà connecté à serveur 4\
Server 2 : <br>
`2`
Server 1 :\
`GO`


_Demande de mutex_

Server 2 : <br>
`req 2`\
Server 1 : <br> 
`token 1`

_Réplication_

Server 2 : <br>
`BOOK 1 1 2 Pierre`\
> Serveur 1 à reçu la demande de réplication, l'a effectué, l'a transmise à 4 qui a confirmé la réplication\

Server 1 : <br> 
`OK`

# Tests automatique
Un program de tests a été mis en place. Le programme de tests se comporte comme un seul client. Toutes les requêtes, qu'un client peut émettre, sont testées avec leurs retours
positifs (pas d’erreurs) et négatifs (paramètres incorrects, indisponibilités, …). 

Pour que le programme de tests fonctionne, il faut lancer les 5 serveurs conformément à la section **Installation** tout en utilisant le fichier de configuration _config.json_ qui y est présenté. 

Pour lancer les tests, depuis le dossier <i>server</i>, exécuter la commande: `$go test -v`

_Attention_, certains critères sont nécessaires pour le bon fonctionnement des tests :
- Les serveurs doivent être fraichement lancé, avec un hotel propre.
- Le nombre de chambres doit être d'exactement 2
- Le nombre de nuits ne doit pas excéder 10
- Il faut avoir un minimum de 2 serveurs
- Ne pas activer le mode debug

Le lancement des tests vérifie ces critères et indique les problèmes potentiels.

# Debug de la concurrence
Comme présenté dans la rubrique "Installation", le serveur peut être lancé en mode debug grâce au paramètre debug dans le fichier de configuration. Ce mode de fonctionnement affiche les événements dans la console. 
De plus, lorsqu'un processus entre dans le mutex, il se met en pause afin d'avoir le temps de lancer une autre commande depuis un autre client.
De ce fait, il est possible de vérifier que la gestion du mutex est conforme.

## Distinctions
Il existe deux types de log:
* RISK: log une requête effectuée dans la zone partagée localement entre les clients. (goroutine gérant l'hôtel, hostelManager).
* SAFE: log une réception ou un envoi depuis/vers le client (goroutine gérant la communication spécifique à chaque client, clientHandler). Ce type de log peut apparaître au milieu d'un passage en zone concurrente sans problème.
* MUTEX: log les interactions avec le mutex. Le log affiche la demande, l'attente, l'entrée et la sortie du mutex. 

Théoriquement, pour une bonne gestion de la concurrence, les logs qui indiquent un passage (entrée puis sortie) en zone partagé ne doivent pas se chevaucher.

__Correct__ \
`1  DEBUG >> Nov 25 15:25:37 SAFE) From 127.0.0.1:1280: LOGIN Alec`\
`2  DEBUG >> Nov 25 15:25:37 RISK) --------- Enter shared zone ---------`\
`3  DEBUG >> Nov 25 15:25:35 SAFE) From 127.0.0.1:1074: LOGIN Melvyn`\
`4  DEBUG >> Nov 25 15:25:37 RISK) LOGIN with username: Alec HANDLING`\
`5  DEBUG >> Nov 25 15:25:37 MUTEX) --------- Asking ---------`\
`6  DEBUG >> Nov 25 15:25:37 MUTEX) --------- Waiting ---------`\
`7  DEBUG >> Nov 25 15:25:49 MUTEX) --------- Entering ---------`\
`8  DEBUG >> Nov 25 15:26:04 RISK) RESULT_LOGIN SUCCESS`\
`9  DEBUG >> Nov 25 15:26:04 MUTEX) --------- Leaving ---------`\
`10  DEBUG >> Nov 25 15:26:04 RISK) --------- Leave shared zone ---------`\
`11 DEBUG >> Nov 25 15:26:04 SAFE) To 127.0.0.1:1280: RESULT_LOGIN`



> Ligne 2, nous voyons que la goroutine gérant les accès concurrents est entrée en zone partagée, prête à traiter la prochaine demande.\
> Ensuite en ligne 3, par le préfix SAFE nous voyons que le client 127.0.0.1:1074 a envoyé la requête LOGIN Melvyn et que sa goroutine dédiée a reçu sa requête.\
> Ligne 4, nous voyons que la requête a été transmise et que la goroutine qui gère l'hôtel traite la demande.\
Enfin, en ligne 5-6, on voit que la requête demande l'accès au mutex, et doit attendre.\
> Ligne 7, la demande peut enfin être traitée après 12 secondes d'attentes.
Finalement lignes 8 à 11, le traitement est terminé et le mutex ainsi que la zone concurrente sont quittés.\
Cet exemple montre un cas d'exécution correct. Il n'y a qu'une exécution au sein du même passage en zone critique, ainsi qu'un seul serveur dans le mutex. De plus, aucune nouvelle entrée en section critique, ou dans le mutex, n'est effectuée tant que la première n'est pas sortie.

## Vérification manuelle

Pour vérifier manuellement la concurrence suivez les étapes suivantes:
* Editer le fichier de configuration (config.json) afin d'activer le mode debug

* Démarrer tous les serveurs présents dans le fichier de configuration debug\
(au minimum 2 serveurs)\
`$go run . 0`\
`$go run . 1`\
`$go run . 2`\
 ...

* Démarrer trois clients (A, B et C) \
`$go run . 0` (A) \
`$go run . 0` (B) \
`$go run . 1` (C) 

* Identifier les clients
  * A \
  `LOGIN A`
  * B \
  `LOGIN B`
  * C \
  `LOGIN C`
> Lors du traitement de la première requête, celle-ci va s'arrêter pendant 15 secondes 
pour laisser le temps aux 2 autres requêtes d'arriver.

* Envoyer la même requête depuis les trois clients
  * A \
  `BOOK 1 1 1`
  * B \
  `BOOK 1 1 1`
  * C \
  `BOOK 1 1 1`

* Attendre la fin des 20 secondes.
> Normalement, un client reçoit une validation et les autres une erreur.

### Analyse d'une sortie des logs.

__Server 0__\
`1 DEBUG >> Nov 25 18:54:07 SAFE) To 127.0.0.1:2109: WELCOME { … }`\
`2 DEBUG >> Nov 25 18:54:09 SAFE) To 127.0.0.1:2116: WELCOME { … }`\
`3 DEBUG >> Nov 25 18:54:14 SAFE) From 127.0.0.1:2109: LOGIN A`\
`4 DEBUG >> Nov 25 18:54:14 RISK) --------- Enter shared zone ---------`\
`5 DEBUG >> Nov 25 18:54:14 RISK) LOGIN with username: A HANDLING`\
`6 DEBUG >> Nov 25 18:54:14 MUTEX) --------- Asking ---------`\
`7 DEBUG >> Nov 25 18:54:14 MUTEX) --------- Waiting ---------`\
`8 DEBUG >> Nov 25 18:54:14 MUTEX) --------- Entering ---------`\
`9 DEBUG >> Nov 25 18:54:17 SAFE) From 127.0.0.1:2116: LOGIN B`\
`10 DEBUG >> Nov 25 18:54:29 RISK) RESULT_LOGIN SUCCESS`\
`11 DEBUG >> Nov 25 18:54:29 SAFE) To 127.0.0.1:2109: RESULT_LOGIN`\
`12 DEBUG >> Nov 25 18:54:29 MUTEX) --------- Leaving ---------`\
`13 DEBUG >> Nov 25 18:54:29 RISK) --------- Leave shared zone ---------`\
`14 DEBUG >> Nov 25 18:54:29 RISK) --------- Enter shared zone ---------`\
`15 DEBUG >> Nov 25 18:54:29 RISK) LOGIN with username: B HANDLING`\
`16 DEBUG >> Nov 25 18:54:29 MUTEX) --------- Asking ---------`\
`17 DEBUG >> Nov 25 18:54:29 MUTEX) --------- Waiting ---------`\
`18 DEBUG >> Nov 25 18:54:33 SAFE) From 127.0.0.1:2109: BOOK 1 1 1`\
`19 DEBUG >> Nov 25 18:54:44 RISK) LOGIN with username: C REPLICATING`\
`20 DEBUG >> Nov 25 18:54:44 MUTEX) --------- Entering ---------`\
`21 DEBUG >> Nov 25 18:54:59 RISK) RESULT_LOGIN SUCCESS`\
`22 DEBUG >> Nov 25 18:54:59 MUTEX) --------- Leaving ---------`\
`23 DEBUG >> Nov 25 18:54:59 RISK) --------- Leave shared zone ---------`\
`24 DEBUG >> Nov 25 18:54:59 RISK) --------- Enter shared zone ---------`\
`25 DEBUG >> Nov 25 18:54:59 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`26 DEBUG >> Nov 25 18:54:59 SAFE) To 127.0.0.1:2116: RESULT_LOGIN`\
`27 DEBUG >> Nov 25 18:54:59 MUTEX) --------- Asking ---------`\
`28 DEBUG >> Nov 25 18:54:59 MUTEX) --------- Waiting ---------`\
`29 DEBUG >> Nov 25 18:55:02 SAFE) From 127.0.0.1:2116: BOOK 1 1 1`\
`30 DEBUG >> Nov 25 18:55:14 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING`\
`31 DEBUG >> Nov 25 18:55:14 MUTEX) --------- Entering ---------`\
`32 DEBUG >> Nov 25 18:55:29 RISK) Room already booked ERROR`\
`33 DEBUG >> Nov 25 18:55:29 MUTEX) --------- Leaving ---------`\
`34 DEBUG >> Nov 25 18:55:29 SAFE) To 127.0.0.1:2109: ERROR Room already booked`\
`35 DEBUG >> Nov 25 18:55:29 RISK) --------- Leave shared zone ---------`\
`36 DEBUG >> Nov 25 18:55:29 RISK) --------- Enter shared zone ---------`\
`37 DEBUG >> Nov 25 18:55:29 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`38 DEBUG >> Nov 25 18:55:29 MUTEX) --------- Asking ---------`\
`39 DEBUG >> Nov 25 18:55:29 MUTEX) --------- Waiting ---------`\
`40 DEBUG >> Nov 25 18:55:29 MUTEX) --------- Entering ---------`\
`41 DEBUG >> Nov 25 18:55:44 RISK) Room already booked ERROR`\
`42 DEBUG >> Nov 25 18:55:44 MUTEX) --------- Leaving ---------`\
`43 DEBUG >> Nov 25 18:55:44 SAFE) To 127.0.0.1:2116: ERROR Room already booked`\
`44 DEBUG >> Nov 25 18:55:44 RISK) --------- Leave shared zone ---------`\



> __Ligne 3)__ Le serveur une requête de login du client ayant l'adresse 127.0.0.1:27358

> __Lignes 4-5)__ Le serveur entre dans la section critique pour la requête

> __Lignes 6-8)__ Le serveur demande l'accès au mutex, attend les réponses et entre dans le mutex
 
> __Ligne 9)__ Le serveur reçoit une requête de login du client ayant l'adresse 127.0.0.1:27359, mais le met en attente

> __Lignes 10-12)__ Le serveur fini le traitement de la requête, quitte le mutex et la section critique

> __Lignes 13-14)__ La 2ème requête reçue passe en section critique

> __Lignes 15-20)__ Le serveur demande l'accès au mutex, et attend la réponse. Pendant cette attente, il renvoie le résultat au client A et reçoit une nouvelle requête. On voit que le serveur a dû attendre 15 secondes sur le mutex. Entre temps, le serveur a répliqué une requête fait à un autre serveur (Ligne 19)

> __Ligne 25)__ La requête BOOK du client A est traitée et demande l'accès au mutex.

> __Ligne 30)__ Le serveur réplique la requête actuellement dans le mutex

> __Lignes 32-33)__ Le traitement de la requête BOOK du client est terminé, et une erreur est survenue car la chambre a déjà été réservée sur un autre serveur
 
__Server 1__\
`1 DEBUG >> Nov 25 18:54:11 SAFE) To 127.0.0.1:2117: WELCOME { … }`\
`2 DEBUG >> Nov 25 18:54:21 SAFE) From 127.0.0.1:2117: LOGIN C`\
`3 DEBUG >> Nov 25 18:54:21 RISK) --------- Enter shared zone ---------`\
`4 DEBUG >> Nov 25 18:54:21 RISK) LOGIN with username: C HANDLING`\
`5 DEBUG >> Nov 25 18:54:21 MUTEX) --------- Asking ---------`\
`6 DEBUG >> Nov 25 18:54:21 MUTEX) --------- Waiting ---------`\
`7 DEBUG >> Nov 25 18:54:29 RISK) LOGIN with username: A REPLICATING`\
`8 DEBUG >> Nov 25 18:54:29 MUTEX) --------- Entering ---------`\
`9 DEBUG >> Nov 25 18:54:44 RISK) RESULT_LOGIN SUCCESS`\
`10 DEBUG >> Nov 25 18:54:44 MUTEX) --------- Leaving ---------`\
`11 DEBUG >> Nov 25 18:54:44 RISK) --------- Leave shared zone ---------`\
`12 DEBUG >> Nov 25 18:54:44 SAFE) To 127.0.0.1:2117: RESULT_LOGIN`\
`13 DEBUG >> Nov 25 18:54:47 SAFE) From 127.0.0.1:2117: BOOK 1 1 1`\
`14 DEBUG >> Nov 25 18:54:47 RISK) --------- Enter shared zone ---------`\
`15 DEBUG >> Nov 25 18:54:47 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`16 DEBUG >> Nov 25 18:54:47 MUTEX) --------- Asking ---------`\
`17 DEBUG >> Nov 25 18:54:47 MUTEX) --------- Waiting ---------`\
`18 DEBUG >> Nov 25 18:54:59 RISK) LOGIN with username: B REPLICATING`\
`19 DEBUG >> Nov 25 18:54:59 MUTEX) --------- Entering ---------`\
`20 DEBUG >> Nov 25 18:55:14 RISK) RESULT_BOOK 1 1 1 SUCCESS`\
`21 DEBUG >> Nov 25 18:55:14 MUTEX) --------- Leaving ---------`\
`22 DEBUG >> Nov 25 18:55:14 RISK) --------- Leave shared zone ---------`\
`23 DEBUG >> Nov 25 18:55:14 SAFE) To 127.0.0.1:2117: RESULT_BOOK 1 1 1`\

> __Lignes 13-23)__ Le serveur reçoit une requête BOOK du client C. Il entre en section critique et attend sur le mutex. 
Après avoir obtenu l'accès au mutex, une réponse positive est renvoyé au client. La chambre est réservée. 

__Server 2__\
`1 DEBUG >> Nov 25 18:54:29 RISK) LOGIN with username: A REPLICATING`\
`2 DEBUG >> Nov 25 18:54:44 RISK) LOGIN with username: C REPLICATING`\
`3 DEBUG >> Nov 25 18:54:59 RISK) LOGIN with username: B REPLICATING`\
`4 DEBUG >> Nov 25 18:55:14 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING`\

> __Lignes 1-4)__ Le serveur n'a reçu de requête d'aucun client. Cependant, il a répliqué toutes les requêtes des autres serveurs ayant obtenu une réponse positive.

Comme nous pouvons le voir, aucune entrée en zone partagée n'est effectuée tant que l'entrée précédente n'est pas terminée. Chaque entrée partagée contient le traitement d'une seule requête. En conclusion, la gestion des accès concurrents est correcte.

## go race
L'application a passé le test du go race.

La procédure précédente (Debug de la concurrence -> vérification manuelle) a été exécutée en démarrant les serveurs avec l'argument <i>-race</i>. 
> `$ go run -race . {no du serveur}`

Aucune concurrence n'a été détectée.
