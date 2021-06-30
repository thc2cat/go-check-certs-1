# go-check-certs

Original code from [github.com/timewasted/go-check-certs](https://github.com/timewasted/go-check-certs)

## Modifications

* return codes,
* timeouts,
* verbose options,
* standard output logging

## Usage

```shell
Add TLS services in file Services.txt in the form of hosts.fqdndomain:port
go build
./go-check-certs -days 45 -v 
```

### M/monit integration

```shell
#/etc/monit.d/go-check-certs

check program go-check-certs with path "/local/go-check-certs/go-check-certs -hosts /local/go-check-certs/Services.txt -days 31"
        every "51-59 8 * * 1-5"
        if status > 0 then alert


## https://mmonit.com/monit/documentation/monit.html
## See special remarks on cron :
# ==>  Therefor we strongly recommend to use an asterix 
#      in the minute field or at minimum a range, e..g. 0-15. 
#      Never use a specific minute as Monit may not run on that minute.
#
```
