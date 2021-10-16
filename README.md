# PRR-01_Forestier_Herzig
Respository du laboratoire 01 pour le cours PRR

# Étudiants
- Forestier Quentin
- Herzig Melvyn

# Installation

* Cloner le répertoire.
> `$ git clone https://github.com/MelvynHerzig/PRR-01_Forestier_Herzig.git`

* Démarrer le serveur. Trois arguments sont nécessaires.
  * Nombre de chambres dans l'hotel (obligatoire).
  * Nombre de nuits dans l'hotel (obligatoire).
  * -debug pour lancer en mode debug (facultatif).

> Depuis le dossier <i>server</i>.
>
> Pour lancer un serveur avec 10 chambres sur 10 nuits sans debug: </br>
> `$ go run . 10 10`
>
> Pour lancer un serveur avec 5 chambres sur 20 nuits avec debug: </br>
> `$ go run . 5 20 -debug`

* Démarrer le(s) client(s). Une argument est nécessaire.
  * Adresse ip du serveur. 
> Depuis le dossier <i>client</i>.
>
> Pour lancer un client qui se connecte à un serveur sur la même machine.
> `$ go run . localhost`
>
> Pour lancer un client qui se connecte à un serveur à l'adresse 1.2.3.4
> `$ go run . 1.2.3.4`

# Utilisation
Toutes les fonctionnalités de la donnée du laboratoire ont été implémentées avec succès.

## Serveur
Une fois le serveur lancé, aucune action supplémentaire nécessaire.

## Client
Au démarrage les clients reçoivent la bienvenue du serveur sous cette forme:

` Welcome in the FH Hotel ! Nb rooms: 10, nb nights: 10 ` <br>
`Available commands:` <br>
`-  LOGIN userName` <br>
`-  LOGOUT` <br>
`-  BOOK roomNumber arrivalNight nbNights` <br>
`-  ROOMLIST night` <br>
`-  FREEROOM arrivalNight nbNights` <br>

> Remarquez:
>* le nombre de chambre supportées, ici de 1 à 10.
>* le nombre de nuits supportées, ici de 1 à 10.
>* la liste des commandes, ici LOGIN, LOGOUT, BOOK, ROOMLIST et FREEROOM.

### LOGIN
` LOGIN <userName>`</br>
Première commande à effectuer. Les autres commandes ne fonctionnement pas tant que login avec un nom d'utilisateur n'a pas été exécutée. Les noms d'utilisateur sont supposés unique. En conséquence, deux utilisateurs avec le même nom ne devraient pas se connecter simultanément.

En cas de succès, l'utilisateur reçoit:
> `Login success` </br>

sinon il reçoit un message d'erreur avec une explication.


### BOOK
` BOOK roomNumber arrivalNight nbNights` <br>
Cette commande sert à réserver une chambre de numéro <i>roomNumber</i> à partir de la nuit <i>arrivalNight</i> durant un nombre de nuit <i>nbNights</i>. Cette commande est disponible seulement après un <i>LOGIN</i> avec succès.

Si la commande `BOOK 1 2 3` est effectuée avec succès, le client reçoit:
>`You successfully booked room  1  for  3  night(s), starting night 2`

sinon il reçoit un message d'erreur avec une explication.

### ROOMLIST 
` ROOMLIST night` <br>
Cette commande permet de voir l'état des chambres dans l'hotel pour une nuit donnée <i>night</i>. Cette commande est disponible seulement après un <i>LOGIN</i> avec succès.

Si la commande `ROOMLIST 2` est effectuée avec succès, le client reçoit:
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

> Pour la nuit 2, nous apperçevons que toutes le chambres sont libres sauf la chambre 1, réservée par le client lui même, et la chambre 4, réservée par un autre client.

sinon il reçoit un message d'erreur avec une explication.

### FREEROOM
`FREEROOM arrivalNight nbNights` <br>
Cette commande permet de chercher la première chambre libre à partir d'une nuit <i>arrivalNight</i> pendant un nombre de nuit <i>nbNights</i>. Cette commande est disponible seulement après un <i>LOGIN</i> avec succès.

Si la commande FREEROOM 2 1 est effectuée avec succès, le client reçoit:
>`Room  2  is free from night  2  during  1  night(s).`

> Effectuée au moment du résultat de l'exécution de <i>ROOMLIST</i> précédent.

si aucune chambre est disponible, le client reçoit:
> `No rooms free from night  2  for  1  night(s).`

sinon il reçoit un message d'erreur avec une explication.

### LOGOUT
` LOGOUT`</br>
Cette commande permet à un utilisateur de se déconnecter. Cette commande est disponible seulement après un <i>LOGIN</i> avec succès.

En cas de succès, l'utilisateur reçoit:
> `Logout success` </br>

sinon il reçoit un message d'erreur avec une explication.

# Compatibilité
L'application a été testée est validée sous Windows et Linux.

La compatibilité MacOS n'a pas pu être contrôlée mais devrait être compatible. Seules des fonctionnalités de base de Golang ont été utilisée.

# Protocole de communication TCP
## Comment le client trouve le serveur (adresses et ports)?
Le serveur utilise sont adresse localhost et le port 8000

## Qui parle en premier ? 
Le serveur parle en premier.</br>
Un message de bienvenue suivi d'une liste de commandes est envoyé au client lorsque le client parvient à se connecter au serveur.

## Qui ferme la connexion et quand?
Le client ferme la connexion lorsque il termine son exécution.

## Qu'est ce qui se passe quand un message est reçu?
### Serveur
Le serveur récupère le premier mot de la requête. Si le mot est syntaxiquement:
* Inconnu: renvoie une erreur.
* Connu: tente de former la suite de la requête.

Si la requête peut être formée:
* Execute et retourne le résultat.
* Sinon retourne une erreur.


### Client
Le client récupère le premier mot de la réponse. Il détermine si c'est:
* Un résultat et affiche les détails.
* Une erreure et affiche la raison.

## Syntaxe des messages envoyé par le client au serveur
| Utilité | Syntaxe |
|---|----|
| S'identifier à l'hôtel | LOGIN {nom de l'utilisateur} CRLF |
| Réserver une chambre | BOOK {numéro de chambre} {nuit d'arrivée} {nombre de nuits} CRLF  |
| Récupérer la liste des disponnibilités pour une nuit. | ROOMLIST {numéro de nuit}  CRLF  |
| Recevoir un numéro de chambre libre pour un nombre de nuits à partir d'une nuit d'arrivée. | FREEROOM {nuit d'arrivée} {nombre de nuit} CRLF  |
| Se déconnecter. | LOGOUT CRLF |

## Syntaxe des messages envoyé par le serveur au client
| Utilité | Syntaxe |
|---|----|
| Réponse positive à LOGIN | RESULT_LOGIN CRLF |
| Réponse positive à BOOK | RESULT_BOOK {numéro de chambre} {nuit d'arrivée} {nombre de nuits} CRLF  |
| Réponse positive à ROOMLIST | RESULT_ROOMLIST {état chambre1}, {état chambre2} ...  CRLF  |
| Réponse positive à FREEROOM | RESULT_FREEROOM {no chambre libre ou 0} CRLF |
| Réponse positive à LOGOUT | RESULT_LOGOUT {numéro de chambre/0} {nuit d'arrivée} {nombre de nuits} CRLF |
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

# Debug de la concurrence
Comme présenté dans la rubrique "Installation", le serveur peut être lancé en mode debug grâce au paraètre "-debug". Ce mode de fonctionnement affiche les événements dans la console. De plus, à chaque fois que deux utilisateurs sont connectés(login) avec succès, la goroutine qui gère les ressources partagées se met en pause pendant 20 secondes dans le but de laisser suffisament de temps pour créer une situation de concurrence. 

## Distinctions
Il existe deux types de log:
* RISK: log une requête effectuée dans la zone partagée (goroutine gérant l'hotel, hostelManager).
* SAFE: log une réception ou un envoi depuis/vers le client (goroutine gérant la communication avec les clients, clientHandler). Ce type de log peut apparaître au milieu d'un passage en zone concurrente sans souci.

Théoriquement, pour une bonne gestion de la concurrence, les logs qui indiquent un passage (entrée puis sortie) en zone partagé ne doivent pas se chevaucher.

__Correct__ \
`1 DEBUG >>  RISK)  --------- Enter shared zone ---------`\
`2 DEBUG >>  SAFE) From 127.0.0.1:5155 : BOOK 1 1 1`\
`3 DEBUG >>  RISK) From 127.0.0.1:5155 BOOK room 1 from night1 for 1 night(s) HANDLING`\
`4 DEBUG >>  RISK) From 127.0.0.1:5155 BOOK room 1 from night1 for 1 night(s) SUCCESS`\
`5 DEBUG >>  RISK) --------- Leave shared zone ---------`


> Ligne 1, nous voyons que la goroutine gérant les accès concurrants est entrée en zone partagée, prête à traiter la prochaine demande concurrente. Ensuite en ligne 2, par le prefix SAFE nous voyons que le client 127.0.0.1:5155 a envoyé la requête BOOK 1 1 1 et que sa goroutine dédiée a reçu sa requête. Ligne 3, nous voyons que la requête a été transmise et que la goroutine qui gère l'hotel la traite.
Finalement ligne 4 et 5, le traitement est terminé et la zone concurrente est quittée.
Cet exemple montre un cas d'exécution correct. Il n'y a qu'une exécution au sein du même passage en zone critique. En d'autre terme aucune nouvelle entrée en section critique est effectuée tant que la première n'est pas sortie.

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

> Contrairement à l'exemple précédent, celui-ci montre un cas qui ne doit pas avoir lieu. Nous pouvons voir que deux traitements critiques ont été executé en même temps. En effet, une entrée en zone partagée a été effectuée alors que la précédente entrée n'a pas terminé. Les lignes 2,7,8 et 10 devraient être exécutées après les lignes 1,5,6 et 9.

## Vérification manuelle

Pour vérifier manuellement la concurrence suivez les étapes suivantes:

* Démarrer le serveur en mode debug\
`$go run . 10 10 -debug`

* Démarrer deux clients (A et B) \
`$go run . localhost` (A) \
`$go run . localhost` (B)

* Identifier les clients\
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


> __Lignes 1-18)__ Ces lignes montrent l'enchaînement des appelles lorsque les deux clients se connectent au serveur et que les deux utilisateurs s'authentifient. Comme annoncé, lorsque deux utilisateurs sont authentifiés, la go routine de traîtement des requêtes se met en pause, c'est ce que nous voyons en ligne 12. Ensuite, les deux clients envoient la même requête. Les lignes 13 et 14 indiquent que les goroutines qui communiquent avec les clients ont bien reçu leur requête. La ligne 16 indique la reprise de la goroutine de traitement. Finalement, les lignes 17 et 18 terminent le traitement du login du client 127.0.0.1:13124. 

>__Lignes 19-22)__ La requête BOOK du client 127.0.0.1:13120 est traitée.

>__Lignes 23-26)__ La requête BOOK du client From 127.0.0.1:13124 est traitée. Ligne 25, confirmation de l'échec. Cet échec est dû au fait que la chambre a déjà été réservée par le client précédent.

>__Ligne 27)__ Attente d'un nouveau traitement concurrent.

>__Ligne 28-29)__ Envoi des réponse aux clients.

Comme nous pouvons le voir, aucune entrée en zone partagée n'est effectuée tant que l'entrée précédente n'est pas terminée. Chaque entrée partagée contient le traitement d'une seule requête. En conclusion, la gestion des accès concurrents est correcte.