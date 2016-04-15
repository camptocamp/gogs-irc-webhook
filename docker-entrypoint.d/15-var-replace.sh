#!/bin/bash

DIR="/etc/webhook/*.json"

grep -hoR "%{[A-Za-z0-9_]\+}" $DIR | sort | uniq | while read k; do
  trim=${k%\}}
  n=${trim#%{}       # calm down, vim }
  v=$(eval "echo \$${n}")
  find -L $DIR -type f -exec sed -i "s/${k}/${v}/g" {} \;
done
