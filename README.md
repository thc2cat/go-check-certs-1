
<<<<<<< HEAD
Original code from https://github.com/timewasted/go-check-certs.  

# Modifications 
 - return codes,
 - timeouts,
 - verbose options, for insecure TLS version
 - standard output logging
=======
# README concernant l'implémentation locale à l'UVSQ
L'original du code est accessible sur https://github.com/timewasted/go-check-certs.  La version modifié est disponible via  https://git.dsi.uvsq.fr/thiecail/go-check-certs. 

Le fichier README-ORIG.md contient le readme original ainsi que le  Copyright (c) 2013, Ryan Rogers.

Le code a été légèrement modifié afin de changer :
 - le code de retour, la concurence, les options, 
 - les timeouts de connections et retry,
 - le logging sur la sortie standard.
 - la vérification du niveau de sécurité TLS
>>>>>>> maj doc

# Usage: 

Add TLS services in file Services.txt in the form of hostname:port, then 
```
<<<<<<< HEAD
go build
./go-check-certs -days 45 -v 
=======
- go-check-certs-x.y : le binaire versionné du `git tag` 
- Services.txt : les services SSL accessibles à vérifier au format host:port.
- exclude: certains services issues des exports mais non vérifiable en état.

## Exemple de sortie : 
```
[root@dhcp go-check-certs]# ./go-check-certs -days 90
2020/11/06 mailhost.dsi.uvsq.fr:993: x509: certificate has expired or is not yet valid: current time 2020-11-06T16:34:25+01:00 is after 2010-10-24T16:00:44Z
2020/11/06 sinchro.uvsq.fr:443: 'sinchro.uvsq.fr' (S/N 63F3076A419FA057B586AEE6B5DDD02) expires in roughly 60 days.
2020/11/06 scribo.rambouillet.iut-velizy.uvsq.fr:443: 'scribo.rambouillet.iut-velizy.uvsq.fr' (S/N 2DFCE9082FC90610BEF000DC56E82AE) expires in roughly 73 days.
2020/11/06 ldap-renater.uvsq.fr:636: 'ldap-renater.uvsq.fr' (S/N 9E84500E6C8AEA5F0278A2422745245) expires in roughly 25 days.
```
L'option -v force l'affichage des connections avec un niveau de securité faible inférieur à TLS1.1.
```
~/p/go-check-certs (master) $ ./go-check-certs -hosts Services.txt  -days 600 -v | grep -i insecure
2021/01/15 WARNING: insecure TLS 1.0 version with siham-prorp.si.uvsq.fr:443
2021/01/15 WARNING: insecure TLS 1.0 version with siham-rp1.dsi.uvsq.fr:443
2021/01/15 WARNING: insecure TLS 1.0 version with mailhost.dsi.uvsq.fr:993
2021/01/15 WARNING: insecure TLS 1.0 version with ent.uvsq.fr:443
2021/01/15 WARNING: insecure TLS 1.0 version with siham-prerp.si.uvsq.fr:443
2021/01/15 WARNING: insecure TLS 1.0 version with sifacweb.uvsq.fr:443

```

## Vérification via Crontabs :
>>>>>>> maj doc
```
