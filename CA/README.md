# Trusting TERENA Certificates

## Ubuntu : 

```
 cd /usr/local/share/ca-certificates
 mkdir TERENA
 scp TERENASSLCA3.crt here.

# update-ca-certificates
Updating certificates in /etc/ssl/certs...
1 added, 0 removed; done.
```

## Centos:

```
cd /etc/pki/ca-trust/source/
scp TERENASSLCA3.crt here. 
# update-ca-trust
```


