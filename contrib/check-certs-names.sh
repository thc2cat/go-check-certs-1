
How to view certificates names with alternate host names 

./go-check-certs -years 4 | rg -v TERENA | awk -F '[ \t:,\']' '{ if($2!=$6) print $0 }'
