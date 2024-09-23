#! /bin/bash

host='mysql.comame.dev'
database='note'

mysqldump -h"$host" -uroot -p --databases "$database" \
    --compact \
    --no-data \
    --no-create-db \
    --skip-comments \
| sed -r '/^USE/d' \
| sed -r '/^\/\*!/d' \
| sed -r 's/;$/;\n/g' \
> tables.sql
