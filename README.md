
Original code from https://github.com/timewasted/go-check-certs.  

# Modifications 
 - return codes,
 - timeouts,
 - verbose options, for insecure TLS version
 - standard output logging

# Usage: 

Add TLS services in file Services.txt in the form of hostname:port, then 
```
go build
./go-check-certs -days 45 -v 
```
