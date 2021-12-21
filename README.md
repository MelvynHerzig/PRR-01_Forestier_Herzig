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


## Exemple d'une conversation entre 2 serveurs Server 2 et Server 1 (en fonction de la configuration initiale)

_Synchronisation (une fois au démarrage)_ 

(Server 1 est déjà connecté à serveur 4)\
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
`BOOK 1 1 2 Pierre`
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
Il existe quatres types de log:
* RISK: log une requête effectuée dans la zone partagée localement entre les clients. (goroutine gérant l'hôtel, hostelManager).
* SAFE: log une réception ou un envoi depuis/vers le client (goroutine gérant la communication spécifique à chaque client, clientHandler). Ce type de log peut apparaître au milieu d'un passage en zone concurrente sans problème.
* MUTEX: log les interactions avec le mutex. Le log affiche la demande, l'attente, l'entrée et la sortie du mutex.
* SERVER: log les interactions entre les serveurs. Notamment les réplications et les transmisions du token.

Théoriquement, pour une bonne gestion de la concurrence, les logs qui indiquent un passage (entrée puis sortie) en zone partagé ne doivent pas se chevaucher.

Principes (Log simplifiés) \

`1  DEBUG >> Nov 25 15:25:37 SAFE) From 127.0.0.1:1280: LOGIN Quentin`\
`2  DEBUG >> Nov 25 15:25:37 RISK) --------- Enter shared zone ---------`\
`3  DEBUG >> Nov 25 15:25:35 SAFE) From 127.0.0.1:1074: LOGIN Melvyn`\
`4  DEBUG >> Nov 25 15:25:37 RISK) LOGIN with username: Quentin HANDLING`\
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

### __Scénario d'example__

Pour vérifier manuellement la concurrence suivez les étapes suivantes:
* Editer le fichier de configuration (config.json) afin d'activer le mode debug

* Démarrer tous les serveurs présents dans le fichier de configuration debug\
(au minimum 5 serveurs)\
`$go run . 0`\
`$go run . 1`\
`$go run . 2`\
`$go run . 3`\
`$go run . 4`\
 ...

* Démarrer quatre clients (A, B, C et D) \
`$go run . 0` (A) \
`$go run . 0` (B) \
`$go run . 1` (C) \
`$go run . 3` (D) 

* Identifier les clients
  * A \
  `LOGIN A`
  * B \
  `LOGIN B`
  * C \
  `LOGIN C`
  * D \
  `LOGIN D`

> Lors du traitement de la première requête, celle-ci va s'arrêter pendant 15 secondes 
pour laisser le temps aux 3 autres requêtes d'arriver.

* Envoyer la même requête depuis les quatres clients dès que possible
  * A \
  `BOOK 1 1 1`
  * B \
  `BOOK 1 1 1`
  * C \
  `BOOK 1 1 1`
  * D \
  `BOOK 1 1 1`


> Normalement, un client reçoit une validation et les autres une erreur.

### Analyse d'une sortie des logs.

__Server 0__\
`1  Parent 2 connected `\
`2  DEBUG >> Dec 19 17:29:24 SERVER) 0 to server 2`\
`3  DEBUG >> Dec 19 17:29:28 SERVER) GO to servers [ ]`\
`4  Server 0 ready to handle clients`\
`5  DEBUG >> Dec 19 17:31:10 SAFE) To 127.0.0.1:5152: WELCOME { ... }`\
`6  DEBUG >> Dec 19 17:31:13 SAFE) To 127.0.0.1:5154: WELCOME { ... }`\
`7  DEBUG >> Dec 19 17:31:46 SAFE) From 127.0.0.1:5152: LOGIN A`\
`8  DEBUG >> Dec 19 17:31:46 RISK) --------- Enter shared zone ---------`\
`9  DEBUG >> Dec 19 17:31:46 RISK) LOGIN with username: A HANDLING`\
`10  DEBUG >> Dec 19 17:31:46 MUTEX) --------- Asking ---------`\
`11  DEBUG >> Dec 19 17:31:46 MUTEX) --------- Waiting ---------`\
`12  DEBUG >> Dec 19 17:31:46 SERVER) req 0 to server 2`\
`13  DEBUG >> Dec 19 17:31:46 SERVER) token 2 received`\
`14  DEBUG >> Dec 19 17:31:46 MUTEX) --------- Entering ---------`\
`15  DEBUG >> Dec 19 17:31:49 SAFE) From 127.0.0.1:5154: LOGIN B`\
`16  DEBUG >> Dec 19 17:31:52 SERVER) req 2 received`\
`17  DEBUG >> Dec 19 17:32:01 RISK) RESULT_LOGIN SUCCESS `\
`18  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A to servers [ 2 ]`\
`19  DEBUG >> Dec 19 17:32:01 SERVER) OK received`\
`20  DEBUG >> Dec 19 17:32:01 MUTEX) --------- Leaving ---------`\
`21  DEBUG >> Dec 19 17:32:01 RISK) --------- Leave shared zone ---------`\
`22  DEBUG >> Dec 19 17:32:01 RISK) --------- Enter shared zone ---------`\
`23  DEBUG >> Dec 19 17:32:01 RISK) LOGIN with username: B HANDLING`\
`24  DEBUG >> Dec 19 17:32:01 MUTEX) --------- Asking ---------`\
`25  DEBUG >> Dec 19 17:32:01 SAFE) To 127.0.0.1:5152: RESULT_LOGIN`\
`26  DEBUG >> Dec 19 17:32:01 SERVER) token 0 to server 2`\
`27  DEBUG >> Dec 19 17:32:01 SERVER) req 0 to server 2`\
`28  DEBUG >> Dec 19 17:32:01 MUTEX) --------- Waiting ---------`\
`29  DEBUG >> Dec 19 17:32:09 SAFE) From 127.0.0.1:5152: BOOK 1 1 1`\
`30  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C received`\
`31  DEBUG >> Dec 19 17:32:16 RISK) LOGIN with username: C REPLICATING `\
`32  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C to servers [ ]`\
`33  DEBUG >> Dec 19 17:32:16 SERVER) OK to server 2`\
`34  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D received`\
`35  DEBUG >> Dec 19 17:32:31 RISK) LOGIN with username: D REPLICATING `\
`36  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D to servers [ ]`\
`37  DEBUG >> Dec 19 17:32:31 SERVER) OK to server 2`\
`38  DEBUG >> Dec 19 17:32:31 SERVER) token 2 received`\
`39  DEBUG >> Dec 19 17:32:31 MUTEX) --------- Entering ---------`\
`40  DEBUG >> Dec 19 17:32:31 SERVER) req 2 received`\
`41  DEBUG >> Dec 19 17:32:46 RISK) RESULT_LOGIN SUCCESS `\
`42  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B to servers [ 2 ]`\
`43  DEBUG >> Dec 19 17:32:46 SERVER) OK received`\
`44  DEBUG >> Dec 19 17:32:46 MUTEX) --------- Leaving ---------`\
`45  DEBUG >> Dec 19 17:32:46 RISK) --------- Leave shared zone ---------`\
`46  DEBUG >> Dec 19 17:32:46 RISK) --------- Enter shared zone ---------`\
`47  DEBUG >> Dec 19 17:32:46 SAFE) To 127.0.0.1:5154: RESULT_LOGIN`\
`48  DEBUG >> Dec 19 17:32:46 SERVER) token 0 to server 2`\
`49  DEBUG >> Dec 19 17:32:46 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`50  DEBUG >> Dec 19 17:32:46 MUTEX) --------- Asking ---------`\
`51  DEBUG >> Dec 19 17:32:46 MUTEX) --------- Waiting ---------`\
`52  DEBUG >> Dec 19 17:32:46 SERVER) req 0 to server 2`\
`53  DEBUG >> Dec 19 17:32:51 SAFE) From 127.0.0.1:5154: BOOK 1 1 1`\
`54  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C received`\
`55  DEBUG >> Dec 19 17:33:01 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`56  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C to servers [ ]`\
`57  DEBUG >> Dec 19 17:33:01 SERVER) OK to server 2`\
`58  DEBUG >> Dec 19 17:33:16 SERVER) token 2 received`\
`59  DEBUG >> Dec 19 17:33:16 MUTEX) --------- Entering ---------`\
`60  DEBUG >> Dec 19 17:33:31 RISK) Room already booked ERROR `\
`61  DEBUG >> Dec 19 17:33:31 MUTEX) --------- Leaving ---------`\
`62  DEBUG >> Dec 19 17:33:31 RISK) --------- Leave shared zone ---------`\
`63  DEBUG >> Dec 19 17:33:31 RISK) --------- Enter shared zone ---------`\
`64  DEBUG >> Dec 19 17:33:31 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`65  DEBUG >> Dec 19 17:33:31 MUTEX) --------- Asking ---------`\
`66  DEBUG >> Dec 19 17:33:31 MUTEX) --------- Waiting ---------`\
`67  DEBUG >> Dec 19 17:33:31 MUTEX) --------- Entering ---------`\
`68  DEBUG >> Dec 19 17:33:31 SAFE) To 127.0.0.1:5152: ERROR Room already booked`\
`69  DEBUG >> Dec 19 17:33:46 RISK) Room already booked ERROR `\
`70  DEBUG >> Dec 19 17:33:46 MUTEX) --------- Leaving ---------`\
`71  DEBUG >> Dec 19 17:33:46 RISK) --------- Leave shared zone ---------`\
`72  DEBUG >> Dec 19 17:33:46 SAFE) To 127.0.0.1:5154: ERROR Room already booked`\

> __Ligne 1-4)__ Mise en place de la connexion avec les parents/enfants. On peut voir que ce serveur n'a pas d'enfants au démarrage.

> __Lignes 7-14)__ Le serveur reçoit une demande, entre dans la zone à risque locale, demande le mutex, reçoit le token et entre en SC

> __Ligne 16)__ Reçoit une demande du serveur 2

> __Lignes 18-19)__ Le serveur demande à ses enfants de répliquer la requête. L'enfant répond OK car la réplication a bien été effectuée.

> __Lignes 20-28)__ Le serveur quitte la SC, donne le token car ce n'est pas lui le prochain dans la queue, puis commence à traiter une nouvelle requête. Il redemande l'accès au mutex.

> __Lignes 30-33)__ Le serveur reçoit une requête a répliqué, la réplique, demande à ses enfants de répliquer la requête, reçoit les OK des enfants, puis envoi OK au parent.

> __Lignes 38-39)__ Le serveur reçoit le token et peut donc entrer en SC.

> __Lignes 56-60)__ Le serveur réplique la réservation `BOOK 1 1 1` pour l'utilisateur C. Il obtient ensuite le mutex pour la requête `BOOK 1 1 1` de l'utilisateur A. Cependant, cela est maintenant impossible.
 
__Server 1__\
`1  server 2 connected`\
`2  server 4 connected`\
`3  DEBUG >> Dec 19 17:29:28 SERVER) GO to servers [ 2 4 ]`\
`4  Server 1 ready to handle clients`\
`5  DEBUG >> Dec 19 17:31:24 SAFE) To 127.0.0.1:5174: WELCOME { ... }`\
`6  DEBUG >> Dec 19 17:31:46 SERVER) req 2 received`\
`7  DEBUG >> Dec 19 17:31:46 SERVER) token 1 to server 2`\
`8  DEBUG >> Dec 19 17:31:52 SAFE) From 127.0.0.1:5174: LOGIN C`\
`9  DEBUG >> Dec 19 17:31:52 RISK) --------- Enter shared zone ---------`\
`10  DEBUG >> Dec 19 17:31:52 RISK) LOGIN with username: C HANDLING`\
`11  DEBUG >> Dec 19 17:31:52 MUTEX) --------- Asking ---------`\
`12  DEBUG >> Dec 19 17:31:52 MUTEX) --------- Waiting ---------`\
`13  DEBUG >> Dec 19 17:31:52 SERVER) req 1 to server 2`\
`14  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A received`\
`15  DEBUG >> Dec 19 17:32:01 RISK) LOGIN with username: A REPLICATING `\
`16  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A to servers [ 4 ]`\
`17  DEBUG >> Dec 19 17:32:01 SERVER) OK received`\
`18  DEBUG >> Dec 19 17:32:01 SERVER) OK to server 2`\
`19  DEBUG >> Dec 19 17:32:01 SERVER) token 2 received`\
`20  DEBUG >> Dec 19 17:32:01 MUTEX) --------- Entering ---------`\
`21  DEBUG >> Dec 19 17:32:01 SERVER) req 2 received`\
`22  DEBUG >> Dec 19 17:32:16 RISK) RESULT_LOGIN SUCCESS `\
`23  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C to servers [ 2 4 ]`\
`24  DEBUG >> Dec 19 17:32:16 SERVER) OK received`\
`25  DEBUG >> Dec 19 17:32:16 SERVER) OK received`\
`26  DEBUG >> Dec 19 17:32:16 MUTEX) --------- Leaving ---------`\
`27  DEBUG >> Dec 19 17:32:16 RISK) --------- Leave shared zone ---------`\
`28  DEBUG >> Dec 19 17:32:16 SERVER) token 1 to server 2`\
`29  DEBUG >> Dec 19 17:32:16 SAFE) To 127.0.0.1:5174: RESULT_LOGIN`\
`30  DEBUG >> Dec 19 17:32:20 SAFE) From 127.0.0.1:5174: BOOK 1 1 1`\
`31  DEBUG >> Dec 19 17:32:20 RISK) --------- Enter shared zone ---------`\
`32  DEBUG >> Dec 19 17:32:20 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`33  DEBUG >> Dec 19 17:32:20 MUTEX) --------- Asking ---------`\
`34  DEBUG >> Dec 19 17:32:20 MUTEX) --------- Waiting ---------`\
`35  DEBUG >> Dec 19 17:32:20 SERVER) req 1 to server 2`\
`36  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D received`\
`37  DEBUG >> Dec 19 17:32:31 RISK) LOGIN with username: D REPLICATING `\
`38  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D to servers [ 4 ]`\
`39  DEBUG >> Dec 19 17:32:31 SERVER) OK received`\
`40  DEBUG >> Dec 19 17:32:31 SERVER) OK to server 2`\
`41  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B received`\
`42  DEBUG >> Dec 19 17:32:46 RISK) LOGIN with username: B REPLICATING `\
`43  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B to servers [ 4 ]`\
`44  DEBUG >> Dec 19 17:32:46 SERVER) OK received`\
`45  DEBUG >> Dec 19 17:32:46 SERVER) OK to server 2`\
`46  DEBUG >> Dec 19 17:32:46 SERVER) token 2 received`\
`47  DEBUG >> Dec 19 17:32:46 MUTEX) --------- Entering ---------`\
`48  DEBUG >> Dec 19 17:32:46 SERVER) req 2 received`\
`49  DEBUG >> Dec 19 17:33:01 RISK) RESULT_BOOK 1 1 1 SUCCESS `\
`50  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C to servers [ 2 4 ]`\
`51  DEBUG >> Dec 19 17:33:01 SERVER) OK received`\
`52  DEBUG >> Dec 19 17:33:01 SERVER) OK received`\
`53  DEBUG >> Dec 19 17:33:01 MUTEX) --------- Leaving ---------`\
`54  DEBUG >> Dec 19 17:33:01 RISK) --------- Leave shared zone ---------`\
`55  DEBUG >> Dec 19 17:33:01 SERVER) token 1 to server 2`\
`56  DEBUG >> Dec 19 17:33:01 SAFE) To 127.0.0.1:5174: RESULT_BOOK 1 1 1`\

> __Lignes 1-4)__ Ce serveur étant la racine, il n'essaye pas de se connecter à son parent. Cependant, on peut voir qu'il a des enfants.

__Server 2__\
`1  server 0 connected`\
`2  server 3 connected`\
`3  Parent 1 connected `\
`4  DEBUG >> Dec 19 17:29:26 SERVER) 2 to server 1`\
`5  DEBUG >> Dec 19 17:29:28 SERVER) GO to servers [ 0 3 ]`\
`6  Server 2 ready to handle clients`\
`7  DEBUG >> Dec 19 17:31:46 SERVER) req 0 received`\
`8  DEBUG >> Dec 19 17:31:46 SERVER) req 2 to server 1`\
`9  DEBUG >> Dec 19 17:31:46 SERVER) token 1 received`\
`10  DEBUG >> Dec 19 17:31:46 SERVER) token 2 to server 0`\
`11  DEBUG >> Dec 19 17:31:52 SERVER) req 1 received`\
`12  DEBUG >> Dec 19 17:31:52 SERVER) req 2 to server 0`\
`13  DEBUG >> Dec 19 17:31:55 SERVER) req 3 received`\
`14  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A received`\
`15  DEBUG >> Dec 19 17:32:01 RISK) LOGIN with username: A REPLICATING `\
`16  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A to servers [ 1 3 ]`\
`17  DEBUG >> Dec 19 17:32:01 SERVER) OK received`\
`18  DEBUG >> Dec 19 17:32:01 SERVER) OK received`\
`19  DEBUG >> Dec 19 17:32:01 SERVER) OK to server 0`\
`20  DEBUG >> Dec 19 17:32:01 SERVER) token 0 received`\
`21  DEBUG >> Dec 19 17:32:01 SERVER) token 2 to server 1`\
`22  DEBUG >> Dec 19 17:32:01 SERVER) req 2 to server 1`\
`23  DEBUG >> Dec 19 17:32:01 SERVER) req 0 received`\
`24  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C received`\
`25  DEBUG >> Dec 19 17:32:16 RISK) LOGIN with username: C REPLICATING `\
`26  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C to servers [ 0 3 ]`\
`27  DEBUG >> Dec 19 17:32:16 SERVER) OK received`\
`28  DEBUG >> Dec 19 17:32:16 SERVER) OK received`\
`29  DEBUG >> Dec 19 17:32:16 SERVER) OK to server 1`\
`30  DEBUG >> Dec 19 17:32:16 SERVER) token 1 received`\
`31  DEBUG >> Dec 19 17:32:16 SERVER) token 2 to server 3`\
`32  DEBUG >> Dec 19 17:32:16 SERVER) req 2 to server 3`\
`33  DEBUG >> Dec 19 17:32:20 SERVER) req 1 received`\
`34  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D received`\
`35  DEBUG >> Dec 19 17:32:31 RISK) LOGIN with username: D REPLICATING `\
`36  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D to servers [ 0 1 ]`\
`37  DEBUG >> Dec 19 17:32:31 SERVER) OK received`\
`38  DEBUG >> Dec 19 17:32:31 SERVER) OK received`\
`39  DEBUG >> Dec 19 17:32:31 SERVER) OK to server 3`\
`40  DEBUG >> Dec 19 17:32:31 SERVER) token 3 received`\
`41  DEBUG >> Dec 19 17:32:31 SERVER) token 2 to server 0`\
`42  DEBUG >> Dec 19 17:32:31 SERVER) req 2 to server 0`\
`43  DEBUG >> Dec 19 17:32:40 SERVER) req 3 received`\
`44  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B received`\
`45  DEBUG >> Dec 19 17:32:46 RISK) LOGIN with username: B REPLICATING `\
`46  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B to servers [ 3 1 ]`\
`47  DEBUG >> Dec 19 17:32:46 SERVER) OK received`\
`48  DEBUG >> Dec 19 17:32:46 SERVER) OK received`\
`49  DEBUG >> Dec 19 17:32:46 SERVER) OK to server 0`\
`50  DEBUG >> Dec 19 17:32:46 SERVER) token 0 received`\
`51  DEBUG >> Dec 19 17:32:46 SERVER) token 2 to server 1`\
`52  DEBUG >> Dec 19 17:32:46 SERVER) req 2 to server 1`\
`53  DEBUG >> Dec 19 17:32:46 SERVER) req 0 received`\
`54  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C received`\
`55  DEBUG >> Dec 19 17:33:01 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`56  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C to servers [ 3 0 ]`\
`57  DEBUG >> Dec 19 17:33:01 SERVER) OK received`\
`58  DEBUG >> Dec 19 17:33:01 SERVER) OK received`\
`59  DEBUG >> Dec 19 17:33:01 SERVER) OK to server 1`\
`60  DEBUG >> Dec 19 17:33:01 SERVER) token 1 received`\
`61  DEBUG >> Dec 19 17:33:01 SERVER) token 2 to server 3`\
`62  DEBUG >> Dec 19 17:33:01 SERVER) req 2 to server 3`\
`63  DEBUG >> Dec 19 17:33:16 SERVER) token 3 received`\
`64  DEBUG >> Dec 19 17:33:16 SERVER) token 2 to server 0`\

> Aucune nouvelle sitatution particulière ne s'est produite sur ce serveur.


__Server 3__\
`1  Parent 2 connected `\
`2  DEBUG >> Dec 19 17:29:26 SERVER) 3 to server 2`\
`3  DEBUG >> Dec 19 17:29:28 SERVER) GO to servers [ ]`\
`4  Server 3 ready to handle clients`\
`5  DEBUG >> Dec 19 17:31:32 SAFE) To 127.0.0.1:5185: WELCOME { ... }`\
`6  DEBUG >> Dec 19 17:31:55 SAFE) From 127.0.0.1:5185: LOGIN D`\
`7  DEBUG >> Dec 19 17:31:55 RISK) --------- Enter shared zone ---------`\
`8  DEBUG >> Dec 19 17:31:55 RISK) LOGIN with username: D HANDLING`\
`9  DEBUG >> Dec 19 17:31:55 MUTEX) --------- Asking ---------`\
`10  DEBUG >> Dec 19 17:31:55 MUTEX) --------- Waiting ---------`\
`11  DEBUG >> Dec 19 17:31:55 SERVER) req 3 to server 2`\
`12  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A received`\
`13  DEBUG >> Dec 19 17:32:01 RISK) LOGIN with username: A REPLICATING `\
`14  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A to servers [ ]`\
`15  DEBUG >> Dec 19 17:32:01 SERVER) OK to server 2`\
`16  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C received`\
`17  DEBUG >> Dec 19 17:32:16 RISK) LOGIN with username: C REPLICATING `\
`18  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C to servers [ ]`\
`19  DEBUG >> Dec 19 17:32:16 SERVER) OK to server 2`\
`20  DEBUG >> Dec 19 17:32:16 SERVER) token 2 received`\
`21  DEBUG >> Dec 19 17:32:16 SERVER) req 2 received`\
`22  DEBUG >> Dec 19 17:32:16 MUTEX) --------- Entering ---------`\
`23  DEBUG >> Dec 19 17:32:31 RISK) RESULT_LOGIN SUCCESS `\
`24  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D to servers [ 2 ]`\
`25  DEBUG >> Dec 19 17:32:31 SERVER) OK received`\
`26  DEBUG >> Dec 19 17:32:31 MUTEX) --------- Leaving ---------`\
`27  DEBUG >> Dec 19 17:32:31 RISK) --------- Leave shared zone ---------`\
`28  DEBUG >> Dec 19 17:32:31 SERVER) token 3 to server 2`\
`29  DEBUG >> Dec 19 17:32:31 SAFE) To 127.0.0.1:5185: RESULT_LOGIN`\
`30  DEBUG >> Dec 19 17:32:40 SAFE) From 127.0.0.1:5185: BOOK 1 1 1`\
`31  DEBUG >> Dec 19 17:32:40 RISK) --------- Enter shared zone ---------`\
`32  DEBUG >> Dec 19 17:32:40 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`33  DEBUG >> Dec 19 17:32:40 MUTEX) --------- Asking ---------`\
`34  DEBUG >> Dec 19 17:32:40 MUTEX) --------- Waiting ---------`\
`35  DEBUG >> Dec 19 17:32:40 SERVER) req 3 to server 2`\
`36  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B received`\
`37  DEBUG >> Dec 19 17:32:46 RISK) LOGIN with username: B REPLICATING `\
`38  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B to servers [ ]`\
`39  DEBUG >> Dec 19 17:32:46 SERVER) OK to server 2`\
`40  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C received`\
`41  DEBUG >> Dec 19 17:33:01 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`42  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C to servers [ ]`\
`43  DEBUG >> Dec 19 17:33:01 SERVER) OK to server 2`\
`44  DEBUG >> Dec 19 17:33:01 SERVER) token 2 received`\
`45  DEBUG >> Dec 19 17:33:01 SERVER) req 2 received`\
`46  DEBUG >> Dec 19 17:33:01 MUTEX) --------- Entering ---------`\
`47  DEBUG >> Dec 19 17:33:16 RISK) Room already booked ERROR `\
`48  DEBUG >> Dec 19 17:33:16 MUTEX) --------- Leaving ---------`\
`49  DEBUG >> Dec 19 17:33:16 RISK) --------- Leave shared zone ---------`\
`50  DEBUG >> Dec 19 17:33:16 SERVER) token 3 to server 2`\
`51  DEBUG >> Dec 19 17:33:16 SAFE) To 127.0.0.1:5185: ERROR Room already booked`\

> Aucune nouvelle sitatution particulière ne s'est produite sur ce serveur.

__Server 4__\
`1  Parent 1 connected `\
`2  DEBUG >> Dec 19 17:29:28 SERVER) 4 to server 1`\
`3  DEBUG >> Dec 19 17:29:28 SERVER) GO to servers [ ]`\
`4  Server 4 ready to handle clients`\
`5  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A received`\
`6  DEBUG >> Dec 19 17:32:01 RISK) LOGIN with username: A REPLICATING `\
`7  DEBUG >> Dec 19 17:32:01 SERVER) LOGIN A to servers [ ]`\
`8  DEBUG >> Dec 19 17:32:01 SERVER) OK to server 1`\
`9  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C received`\
`10  DEBUG >> Dec 19 17:32:16 RISK) LOGIN with username: C REPLICATING `\
`11  DEBUG >> Dec 19 17:32:16 SERVER) LOGIN C to servers [ ]`\
`12  DEBUG >> Dec 19 17:32:16 SERVER) OK to server 1`\
`13  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D received`\
`14  DEBUG >> Dec 19 17:32:31 RISK) LOGIN with username: D REPLICATING `\
`15  DEBUG >> Dec 19 17:32:31 SERVER) LOGIN D to servers [ ]`\
`16  DEBUG >> Dec 19 17:32:31 SERVER) OK to server 1`\
`17  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B received`\
`18  DEBUG >> Dec 19 17:32:46 RISK) LOGIN with username: B REPLICATING `\
`19  DEBUG >> Dec 19 17:32:46 SERVER) LOGIN B to servers [ ]`\
`20  DEBUG >> Dec 19 17:32:46 SERVER) OK to server 1`\
`21  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C received`\
`22  DEBUG >> Dec 19 17:33:01 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`23  DEBUG >> Dec 19 17:33:01 SERVER) BOOK 1 1 1 C to servers [ ]`\
`24  DEBUG >> Dec 19 17:33:01 SERVER) OK to server 1`\

> Aucune nouvelle sitatution particulière ne s'est produite sur ce serveur.

__Globalement__
- On peut voir que tous les serveurs répliquent les requêtes. 
- On vérifie que tous se passe bien pour chaque enfant avant d'envoyer le OK au parent.
- Une entrée en SC arrive toujours après une sortie (sauf la première fois).

### __Branches d'exécution possibles__
Voici la liste des situations possible :

Lors d'une demande :
- A) En tant que noeud non racine, ayant déjà fait une demande ( state = demanding) -> mettre dans la queue  la demande.
- B) En tant que noeud non racine, n'ayant pas déjà fait une demande (state = noDemand) -> mettre dans la queue la demande et passage de state à inDemand
- C) En tant que racine -> accède au mutex

Lors de la sortie de section critique :
- D) La queue est vide -> repasse à noDemand et la queue est vide
- E) La queue contient 1 demande -> repasse à noDemand et dépile la prochaine requête depuis la queue, envoie le token, la queue est vide
- F) La queue contient 2 ou plus demandes -> repasse à noDemand et dépile la prochaine requête depuis la queue, envoie le token, la queue n'est pas vide, envoie donc une REQ pour récupérer le token

Traitement d'une requête en provenance d'un autre serveur :
- G) le serveur est racine avec un état noDemand -> l'expéditeur devient sont parent et il lui envoie le token
- H) Le serveur est en SC OU le serveur n'est pas racine mais a déjà envoyé une demande -> stocke la requête dans sa queue.
- I) le serveur n'est pas racine et n'a pas déjà envoyé une demande -> stocke la requête dans sa queue et envoie une demande à son parent.

Traîtement du token en provenance de son parent :
- J) dépile la queue, la requête dépilée vient du serveur courrant -> le serveur devient la racine et entre en SC
- k) dépile la queue, la requête provient d'un enfant, la queue est vide -> l'enfant devient le parent, le token lui est passé, passage à l'état noDemand
- L) dépile la queue, la requête provient d'un enfant, la queue n'est pas vide -> l'enfant devient le parent, le token lui es passé, passage à l'état demanding et envoie une requête à son parent.


Avec l'arborescence des serveurs ainsi que les demandes suivantes :
![Arboresence et requêtes](https://user-images.githubusercontent.com/61196479/146681150-c8d098f1-c451-4cd9-8616-29ad467cf31e.png)

> À noter : Nous ne pouvons pas garantir à 100% l'ordre d'exécution. Dans les logs qui suivent, la requête `LOGIN 0` a été exécuté APRÈS la requête `LOGIN 4` malgré l'ordre.

Nous obtenons les logs suivants :

__Serveur 0__\
`1  Parent 2 connected `\
`2  DEBUG >> Dec 16 17:55:45 SERVER) 0 to server 2`\
`3  DEBUG >> Dec 16 17:55:49 SERVER) GO to servers [ ]`\
`4  Server 0 ready to handle clients`\
`5  DEBUG >> Dec 16 17:55:58 SAFE) To 127.0.0.1:28422: WELCOME { ... }`\
`6  DEBUG >> Dec 16 17:56:34 SAFE) From 127.0.0.1:28422: LOGIN 0`\
`7  DEBUG >> Dec 16 17:56:34 RISK) --------- Enter shared zone ---------`\
`8  DEBUG >> Dec 16 17:56:34 RISK) LOGIN with username: 0 HANDLING `\
`9  DEBUG >> Dec 16 17:56:34 MUTEX) --------- Asking ---------`\
`10  DEBUG >> Dec 16 17:56:34 MUTEX) --------- Waiting ---------`\
`12  DEBUG >> Dec 16 17:56:34 SERVER) Scénario B`\
`13  DEBUG >> Dec 16 17:56:34 SERVER) req 0 to server 2`\
`14  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 received`\
`15  DEBUG >> Dec 16 17:56:43 RISK) LOGIN with username: 3 REPLICATING `\
`16  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 to servers [ ]`\
`17  DEBUG >> Dec 16 17:56:43 SERVER) OK to server 2`\
`18  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 received`\
`19  DEBUG >> Dec 16 17:56:58 RISK) LOGIN with username: 1 REPLICATING `\
`20  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 to servers [ ]`\
`21  DEBUG >> Dec 16 17:56:58 SERVER) OK to server 2`\
`22  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 received`\
`23  DEBUG >> Dec 16 17:57:13 RISK) LOGIN with username: 4 REPLICATING `\
`24  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 to servers [ ]`\
`25  DEBUG >> Dec 16 17:57:13 SERVER) OK to server 2`\
`26  DEBUG >> Dec 16 17:57:13 SERVER) token 2 received`\
`27  DEBUG >> Dec 16 17:57:13 SERVER) Scénario J`\
`28  DEBUG >> Dec 16 17:57:13 MUTEX) --------- Entering ---------`\
`29  DEBUG >> Dec 16 17:57:28 RISK) RESULT_LOGIN SUCCESS `\
`30  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 to servers [ 2 ]`\
`31  DEBUG >> Dec 16 17:57:28 SERVER) OK received`\
`32  DEBUG >> Dec 16 17:57:28 MUTEX) --------- Leaving ---------`\
`33  DEBUG >> Dec 16 17:57:28 RISK) --------- Leave shared zone ---------`\
`34  DEBUG >> Dec 16 17:57:28 SERVER) Scénario D`\
`35  DEBUG >> Dec 16 17:57:28 SAFE) To 127.0.0.1:28422: RESULT_LOGIN`\
`36  DEBUG >> Dec 16 18:00:12 SAFE) From 127.0.0.1:28422: BOOK 1 1 1`\
`37  DEBUG >> Dec 16 18:00:12 RISK) --------- Enter shared zone ---------`\
`38  DEBUG >> Dec 16 18:00:12 RISK) BOOK room 1 from night 1 for 1 night(s) HANDLING`\
`39  DEBUG >> Dec 16 18:00:12 MUTEX) --------- Asking ---------`\
`40  DEBUG >> Dec 16 18:00:12 MUTEX) --------- Waiting ---------`\
`41  DEBUG >> Dec 16 18:00:12 SERVER) Scénario C`\
`42  DEBUG >> Dec 16 18:00:12 MUTEX) --------- Entering ---------`\
`43  DEBUG >> Dec 16 18:00:27 RISK) RESULT_BOOK 1 1 1 SUCCESS `\
`44  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 to servers [ 2 ]`\
`45  DEBUG >> Dec 16 18:00:27 SERVER) OK received`\
`46  DEBUG >> Dec 16 18:00:27 MUTEX) --------- Leaving ---------`\
`47  DEBUG >> Dec 16 18:00:27 RISK) --------- Leave shared zone ---------`\
`48  DEBUG >> Dec 16 18:00:27 SAFE) To 127.0.0.1:28422: RESULT_BOOK 1 1 1`\
`49  DEBUG >> Dec 16 18:00:27 SERVER) Scénario D`

> Scénario B : Le serveur n'est pas racine, n'a pas encore de demande. Il fait donc la demande (ligne  13).\
  Scénario J : La prochaine requête dans la queue est la sienne. Il passe donc en SC (ligne 28) après avoir reçu le token (ligne 26) \
  Scénario D : Le serveur n'a plus rien dans sa queue et ne fait qu'attendre à la sortie de SC (ligne 33). \
  Scénario C : Le token a été reçu (ligne 26) et n'a pas été renvoyé. Il peut donc passer directement en SC.

__Server 1__\
`1   server 2 connected`\
`2   server 4 connected`\
`3   DEBUG >> Dec 16 17:55:49 SERVER) GO to servers [ 2 4 ]`\
`4   Server 1 ready to handle clients`\
`5   DEBUG >> Dec 16 17:55:57 SAFE) To 127.0.0.1:28418: WELCOME { ... }`\
`6   DEBUG >> Dec 16 17:56:28 SERVER) req 2 received`\
`7   DEBUG >> Dec 16 17:56:28 SERVER) Scénario G`\
`8   DEBUG >> Dec 16 17:56:28 SERVER) token 1 to server 2`\
`9   DEBUG >> Dec 16 17:56:31 SAFE) From 127.0.0.1:28418: LOGIN 1`\
`10  DEBUG >> Dec 16 17:56:31 RISK) --------- Enter shared zone ---------`\
`11  DEBUG >> Dec 16 17:56:31 RISK) LOGIN with username: 1 HANDLING `\
`12  DEBUG >> Dec 16 17:56:31 MUTEX) --------- Asking ---------`\
`13  DEBUG >> Dec 16 17:56:31 MUTEX) --------- Waiting ---------`\
`14  DEBUG >> Dec 16 17:56:31 SERVER) Scénario B`\
`15  DEBUG >> Dec 16 17:56:31 SERVER) req 1 to server 2`\
`16  DEBUG >> Dec 16 17:56:36 SERVER) req 4 received`\
`17  DEBUG >> Dec 16 17:56:36 SERVER) Scénario H`\
`18  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 received`\
`19  DEBUG >> Dec 16 17:56:43 RISK) LOGIN with username: 3 REPLICATING `\
`20  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 to servers [ 4 ]`\
`21  DEBUG >> Dec 16 17:56:43 SERVER) OK received`\
`22  DEBUG >> Dec 16 17:56:43 SERVER) OK to server 2`\
`23  DEBUG >> Dec 16 17:56:43 SERVER) token 2 received`\
`24  DEBUG >> Dec 16 17:56:43 SERVER) Scénario J`\
`25  DEBUG >> Dec 16 17:56:43 MUTEX) --------- Entering ---------`\
`26  DEBUG >> Dec 16 17:56:43 SERVER) req 2 received`\
`27  DEBUG >> Dec 16 17:56:43 SERVER) Scénario H`\
`28  DEBUG >> Dec 16 17:56:58 RISK) RESULT_LOGIN SUCCESS `\
`29  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 to servers [ 2 4 ]`\
`30  DEBUG >> Dec 16 17:56:58 SERVER) OK received`\
`31  DEBUG >> Dec 16 17:56:58 SERVER) OK received`\
`32  DEBUG >> Dec 16 17:56:58 MUTEX) --------- Leaving ---------`\
`33  DEBUG >> Dec 16 17:56:58 RISK) --------- Leave shared zone ---------`\
`34  DEBUG >> Dec 16 17:56:58 SERVER) token 1 to server 4`\
`35  DEBUG >> Dec 16 17:56:58 SAFE) To 127.0.0.1:28418: RESULT_LOGIN`\
`36  DEBUG >> Dec 16 17:56:58 SERVER) req 1 to server 4`\
`37  DEBUG >> Dec 16 17:56:58 SERVER) Scénario F`\
`38  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 received`\
`39  DEBUG >> Dec 16 17:57:13 RISK) LOGIN with username: 4 REPLICATING `\
`40  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 to servers [ 2 ]`\
`41  DEBUG >> Dec 16 17:57:13 SERVER) OK received`\
`42  DEBUG >> Dec 16 17:57:13 SERVER) OK to server 4`\
`43  DEBUG >> Dec 16 17:57:13 SERVER) token 4 received`\
`44  DEBUG >> Dec 16 17:57:13 SERVER) token 1 to server 2`\
`45  DEBUG >> Dec 16 17:57:13 SERVER) Scénario K`\
`46  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 received`\
`47  DEBUG >> Dec 16 17:57:28 RISK) LOGIN with username: 0 REPLICATING `\
`48  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 to servers [ 4 ]`\
`49  DEBUG >> Dec 16 17:57:28 SERVER) OK received`\
`50  DEBUG >> Dec 16 17:57:28 SERVER) OK to server 2`\
`51  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 received`\
`52  DEBUG >> Dec 16 18:00:27 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`53  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 to servers [ 4 ]`\
`54  DEBUG >> Dec 16 18:00:27 SERVER) OK received`\
`55  DEBUG >> Dec 16 18:00:27 SERVER) OK to server 2`

> Scénario G : Le serveur est racine, mais est en attente. Il reçoit une demande et donc transmet son token (ligne 6).\
  Scénario H : Le serveur est en SC (ligne 25) et reçoit une demande (ligne 26). \
  Scénario F : Le serveur donne le token (ligne 34), puis effectue une demande (ligne 36). 

__Server 2__\
`1   server 0 connected`\
`2   server 3 connected`\
`3   Parent 1 connected `\
`4   DEBUG >> Dec 16 17:55:47 SERVER) 2 to server 1`\
`5   DEBUG >> Dec 16 17:55:49 SERVER) GO to servers [ 0 3 ]`\
`6   Server 2 ready to handle clients`\
`7   DEBUG >> Dec 16 17:56:28 SERVER) req 3 received`\
`8   DEBUG >> Dec 16 17:56:28 SERVER) Scénario I`\
`9   DEBUG >> Dec 16 17:56:28 SERVER) req 2 to server 1`\
`10  DEBUG >> Dec 16 17:56:28 SERVER) token 1 received`\
`11  DEBUG >> Dec 16 17:56:28 SERVER) token 2 to server 3`\
`12  DEBUG >> Dec 16 17:56:28 SERVER) Scénario K`\
`13  DEBUG >> Dec 16 17:56:31 SERVER) req 1 received`\
`14  DEBUG >> Dec 16 17:56:31 SERVER) Scénario I`\
`15  DEBUG >> Dec 16 17:56:31 SERVER) req 2 to server 3`\
`16  DEBUG >> Dec 16 17:56:34 SERVER) req 0 received`\
`17  DEBUG >> Dec 16 17:56:34 SERVER) Scénario H`\
`18  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 received`\
`19  DEBUG >> Dec 16 17:56:43 RISK) LOGIN with username: 3 REPLICATING `\
`20  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 to servers [ 1 0 ]`\
`21  DEBUG >> Dec 16 17:56:43 SERVER) OK received`\
`22  DEBUG >> Dec 16 17:56:43 SERVER) OK received`\
`23  DEBUG >> Dec 16 17:56:43 SERVER) OK to server 3`\
`24  DEBUG >> Dec 16 17:56:43 SERVER) token 3 received`\
`25  DEBUG >> Dec 16 17:56:43 SERVER) token 2 to server 1`\
`26  DEBUG >> Dec 16 17:56:43 SERVER) Scénario L`\
`27  DEBUG >> Dec 16 17:56:43 SERVER) req 2 to server 1`\
`28  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 received`\
`29  DEBUG >> Dec 16 17:56:58 RISK) LOGIN with username: 1 REPLICATING `\
`30  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 to servers [ 0 3 ]`\
`31  DEBUG >> Dec 16 17:56:58 SERVER) OK received`\
`32  DEBUG >> Dec 16 17:56:58 SERVER) OK received`\
`33  DEBUG >> Dec 16 17:56:58 SERVER) OK to server 1`\
`34  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 received`\
`35  DEBUG >> Dec 16 17:57:13 RISK) LOGIN with username: 4 REPLICATING `\
`36  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 to servers [ 0 3 ]`\
`37  DEBUG >> Dec 16 17:57:13 SERVER) OK received`\
`38  DEBUG >> Dec 16 17:57:13 SERVER) OK received`\
`39  DEBUG >> Dec 16 17:57:13 SERVER) OK to server 1`\
`40  DEBUG >> Dec 16 17:57:13 SERVER) token 1 received`\
`41  DEBUG >> Dec 16 17:57:13 SERVER) token 2 to server 0`\
`42  DEBUG >> Dec 16 17:57:13 SERVER) Scénario K`\
`43  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 received`\
`44  DEBUG >> Dec 16 17:57:28 RISK) LOGIN with username: 0 REPLICATING `\
`45  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 to servers [ 1 3 ]`\
`46  DEBUG >> Dec 16 17:57:28 SERVER) OK received`\
`47  DEBUG >> Dec 16 17:57:28 SERVER) OK received`\
`48  DEBUG >> Dec 16 17:57:28 SERVER) OK to server 0`\
`49  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 received`\
`50  DEBUG >> Dec 16 18:00:27 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`51  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 to servers [ 1 3 ]`\
`52  DEBUG >> Dec 16 18:00:27 SERVER) OK received`\
`53  DEBUG >> Dec 16 18:00:27 SERVER) OK received`\
`54  DEBUG >> Dec 16 18:00:27 SERVER) OK to server 0`

> Scénario I : Le serveur n'est pas la racine, reçoit une demande (ligne 7) et fait une demande (ligne 9)\
  Scénario K : Le serveur reçoit le token (ligne 10) et le renvoi (ligne 11)\
  Scénario L : Le serveur reçoit le token (ligne 24), le renvoi (ligne 25) et fait une demande (ligne 27)\

__Server 3__\
`1  Parent 2 connected `\
`2  DEBUG >> Dec 16 17:55:47 SERVER) 3 to server 2`\
`3  DEBUG >> Dec 16 17:55:49 SERVER) GO to servers [ ]`\
`4  Server 3 ready to handle clients`\
`5  DEBUG >> Dec 16 17:55:54 SAFE) To 127.0.0.1:28410: WELCOME { ... }`\
`6  DEBUG >> Dec 16 17:56:28 SAFE) From 127.0.0.1:28410: LOGIN 3`\
`7  DEBUG >> Dec 16 17:56:28 RISK) --------- Enter shared zone ---------`\
`8  DEBUG >> Dec 16 17:56:28 RISK) LOGIN with username: 3 HANDLING `\
`9  DEBUG >> Dec 16 17:56:28 MUTEX) --------- Asking ---------`\
`10  DEBUG >> Dec 16 17:56:28 MUTEX) --------- Waiting ---------`\
`11  DEBUG >> Dec 16 17:56:28 SERVER) Scénario B`\
`12  DEBUG >> Dec 16 17:56:28 SERVER) req 3 to server 2`\
`13  DEBUG >> Dec 16 17:56:28 SERVER) token 2 received`\
`14  DEBUG >> Dec 16 17:56:28 SERVER) Scénario J`\
`15  DEBUG >> Dec 16 17:56:28 MUTEX) --------- Entering ---------`\
`16  DEBUG >> Dec 16 17:56:31 SERVER) req 2 received`\
`17  DEBUG >> Dec 16 17:56:31 SERVER) Scénario H`\
`18  DEBUG >> Dec 16 17:56:43 RISK) RESULT_LOGIN SUCCESS `\
`19  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 to servers [ 2 ]`\
`20  DEBUG >> Dec 16 17:56:43 SERVER) OK received`\
`21  DEBUG >> Dec 16 17:56:43 MUTEX) --------- Leaving ---------`\
`22  DEBUG >> Dec 16 17:56:43 RISK) --------- Leave shared zone ---------`\
`23  DEBUG >> Dec 16 17:56:43 SERVER) Scénario E`\
`24  DEBUG >> Dec 16 17:56:43 SAFE) To 127.0.0.1:28410: RESULT_LOGIN`\
`25  DEBUG >> Dec 16 17:56:43 SERVER) token 3 to server 2`\
`26  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 received`\
`27  DEBUG >> Dec 16 17:56:58 RISK) LOGIN with username: 1 REPLICATING `\
`28  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 to servers [ ]`\
`29  DEBUG >> Dec 16 17:56:58 SERVER) OK to server 2`\
`30  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 received`\
`31  DEBUG >> Dec 16 17:57:13 RISK) LOGIN with username: 4 REPLICATING `\
`32  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 to servers [ ]`\
`33  DEBUG >> Dec 16 17:57:13 SERVER) OK to server 2`\
`34  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 received`\
`35  DEBUG >> Dec 16 17:57:28 RISK) LOGIN with username: 0 REPLICATING `\
`36  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 to servers [ ]`\
`37  DEBUG >> Dec 16 17:57:28 SERVER) OK to server 2`\
`38  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 received`\
`39  DEBUG >> Dec 16 18:00:27 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`40  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 to servers [ ]`\
`41  DEBUG >> Dec 16 18:00:27 SERVER) OK to server 2`

> Scénario E : Le serveur qui la SC (ligne 22) et rend le token (ligne 25).

__Server 4__\
`1  Parent 1 connected `\
`2  DEBUG >> Dec 16 17:55:49 SERVER) 4 to server 1`\
`3  DEBUG >> Dec 16 17:55:49 SERVER) GO to servers [ ]`\
`4  Server 4 ready to handle clients`\
`5  DEBUG >> Dec 16 17:56:00 SAFE) To 127.0.0.1:28427: WELCOME { ... }`\
`6  DEBUG >> Dec 16 17:56:36 SAFE) From 127.0.0.1:28427: LOGIN 4`\
`7  DEBUG >> Dec 16 17:56:36 RISK) --------- Enter shared zone ---------`\
`8  DEBUG >> Dec 16 17:56:36 RISK) LOGIN with username: 4 HANDLING `\
`9  DEBUG >> Dec 16 17:56:36 MUTEX) --------- Asking ---------`\
`10  DEBUG >> Dec 16 17:56:36 MUTEX) --------- Waiting ---------`\
`11  DEBUG >> Dec 16 17:56:36 SERVER) Scénario B`\
`12  DEBUG >> Dec 16 17:56:36 SERVER) req 4 to server 1`\
`13  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 received`\
`14  DEBUG >> Dec 16 17:56:43 RISK) LOGIN with username: 3 REPLICATING `\
`15  DEBUG >> Dec 16 17:56:43 SERVER) LOGIN 3 to servers [ ]`\
`16  DEBUG >> Dec 16 17:56:43 SERVER) OK to server 1`\
`17  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 received`\
`18  DEBUG >> Dec 16 17:56:58 RISK) LOGIN with username: 1 REPLICATING `\
`19  DEBUG >> Dec 16 17:56:58 SERVER) LOGIN 1 to servers [ ]`\
`20  DEBUG >> Dec 16 17:56:58 SERVER) OK to server 1`\
`21  DEBUG >> Dec 16 17:56:58 SERVER) token 1 received`\
`22  DEBUG >> Dec 16 17:56:58 SERVER) req 1 received`\
`23  DEBUG >> Dec 16 17:56:58 SERVER) Scénario J`\
`24  DEBUG >> Dec 16 17:56:58 SERVER) Scénario H`\
`25  DEBUG >> Dec 16 17:56:58 MUTEX) --------- Entering ---------`\
`26  DEBUG >> Dec 16 17:57:13 RISK) RESULT_LOGIN SUCCESS `\
`27  DEBUG >> Dec 16 17:57:13 SERVER) LOGIN 4 to servers [ 1 ]`\
`28  DEBUG >> Dec 16 17:57:13 SERVER) OK received`\
`29  DEBUG >> Dec 16 17:57:13 MUTEX) --------- Leaving ---------`\
`30  DEBUG >> Dec 16 17:57:13 RISK) --------- Leave shared zone ---------`\
`31  DEBUG >> Dec 16 17:57:13 SERVER) Scénario D`\
`32  DEBUG >> Dec 16 17:57:13 SERVER) Scénario E`\
`33  DEBUG >> Dec 16 17:57:13 SAFE) To 127.0.0.1:28427: RESULT_LOGIN`\
`34  DEBUG >> Dec 16 17:57:13 SERVER) token 4 to server 1`\
`35  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 received`\
`36  DEBUG >> Dec 16 17:57:28 RISK) LOGIN with username: 0 REPLICATING `\
`37  DEBUG >> Dec 16 17:57:28 SERVER) LOGIN 0 to servers [ ]`\
`38  DEBUG >> Dec 16 17:57:28 SERVER) OK to server 1`\
`39  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 received`\
`40  DEBUG >> Dec 16 18:00:27 RISK) BOOK room 1 from night 1 for 1 night(s) REPLICATING `\
`41  DEBUG >> Dec 16 18:00:27 SERVER) BOOK 1 1 1 0 to servers [ ]`\
`42  DEBUG >> Dec 16 18:00:27 SERVER) OK to server 1`\

__Attention__ : Le scénario A n'est jamais apparu. La suite de requêtes effectués ne le permettait tout simplement pas.

#### Scénario A

Avec l'arborescence des serveurs ainsi que les demandes suivantes :
![Sans titre](https://user-images.githubusercontent.com/61196479/146682304-8fed4a44-d277-4f3e-97ba-4196257defd6.png)

Le scénario A s'effectuera sur le serveur 2.


## go race
L'application a passé le test du go race.

La procédure précédente (Debug de la concurrence -> vérification manuelle) a été exécutée en démarrant les serveurs avec l'argument <i>-race</i>. 
> `$ go run -race . {no du serveur}`

Aucune concurrence n'a été détectée.
