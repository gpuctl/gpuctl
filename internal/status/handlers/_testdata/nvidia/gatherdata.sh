#!/bin/sh
jumphost="shell2.doc.ic.ac.uk"
hosts="gpu01 gpu02 gpu03 gpu04 gpu05 gpu10 gpu30 gpu36 ash01 ash02 ash03 ash04 ash05"
username=$1

for h in $hosts
do
    ssh -J $username@$jumphost $username@$h "nvidia-smi -x -q" > dump_$h.xml
done
