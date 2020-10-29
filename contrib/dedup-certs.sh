#!/bin/sh

mkdir -p certs

for i in `cat ../Services.txt
do
  echo -n | openssl s_client -connect $i \
    | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > certs/$i
done

godedup -path certs
