#! /bin/bash

socket='.testdb/mysql.sock'
database='note'

mysqldump -S"$socket" -uroot --databases "$database" \
    --compact \
    --no-data \
    --no-create-db \
    --skip-comments \
| sed -r '/^USE/d' \
| sed -r '/^\/\*!/d' \
| sed -r 's/;$/;\n/g' \
> tables.sql
