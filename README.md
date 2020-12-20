
# README concernant l'implémentation locale à l'UVSQ
L'original du code est accessible sur https://github.com/timewasted/go-check-certs.  
Le fichier README-ORIG.md contient le readme original ainsi que le copyrigth  ( Copyright (c) 2013, Ryan Rogers )

Le code a été légèrement modifié afin de changer :
 - le code de retour, la concurence, les options, 
 - les timeouts de connections et retry,
 - le logging sur la sortie standard.

## Détails de la solution en place sur mmonit.uvsq.fr  : 
``` 
cd /local/go-check-certs
[root@dhcp go-check-certs]# ls -l
total 5084
-rw-r--r-- 1 root root     272 26 oct.  11:57 exclude.txt
lrwxrwxrwx 1 root root      18  6 nov.  16:28 go-check-certs -> go-check-certs-0.2
-rwxr-xr-x 1 root root 5196147 27 oct.  18:11 go-check-certs-0.2
-rw-r--r-- 1 root root    1830 29 oct.  14:54 Services.txt

```
- go-check-certs-0.2 : le binaire versionné dans https://git.dsi.uvsq.fr/thiecail/go-check-certs
- Services.txt : les services SSL accessibles à vérifier au format host:port.
- exclude: certains services issues des exports mais non vérifiable en état.

## Exemple de sortie : 
```
[root@dhcp go-check-certs]# ./go-check-certs -days 90
2020/11/06 mailhost.dsi.uvsq.fr:993: x509: certificate has expired or is not yet valid: current time 2020-11-06T16:34:25+01:00 is after 2010-10-24T16:00:44Z
2020/11/06 sinchro.uvsq.fr:443: 'sinchro.uvsq.fr' (S/N 63F3076A419FA057B586AEE6B5DDD02) expires in roughly 60 days.
2020/11/06 scribo.rambouillet.iut-velizy.uvsq.fr:443: 'scribo.rambouillet.iut-velizy.uvsq.fr' (S/N 2DFCE9082FC90610BEF000DC56E82AE) expires in roughly 73 days.
2020/11/06 ldap-renater.uvsq.fr:636: 'ldap-renater.uvsq.fr' (S/N 9E84500E6C8AEA5F0278A2422745245) expires in roughly 25 days.

L'option -v force l'affichage des connections avec un niveau de securité faible inférieur à TLS1.1.
```



## Vérification via Crontabs :
```
# Verifier les certificats pour Dominique
0  9  1 * * /local/go-check-certs/go-check-certs -hosts /local/go-check-certs/Services.txt -days 45 2>&1 | mail -s "[CERTS] verification mensuelle" dominique.fiquet@uvsq.fr
10 10 * * 1 /local/go-check-certs/go-check-certs -hosts /local/go-check-certs/Services.txt -days 15 2>&1 | mail -s "[CERTS] verification hebdomadaire" thiecail@poste.uvsq.fr
```

## Vérification quotidienne via Monit : 
```
# cat  /etc/monit.d/go-check-certs

check program go-check-certs with path " /local/go-check-certs/go-check-certs -hosts /local/go-check-certs/Services.txt -days 31"
        every "30 7 * * 1-5"
        if status > 0 then alert

```
A noter que si le code de retour via monit ne change pas (comprendre pas de corrections), l'alerte ne sera pas mise à jour. Donc, pour un problème persistant, il vaut mieux déplacer le service dans le fichier exclude.txt avec un commentaire.

#### Dossier CA
 Contient les certificats d'autorités a rajouter sur les serveurs anciens ne disposant plus de ces autorités.
#### Contrib 
 Contient de quoi downloader / dedupliquer les certificats lors du premier import.