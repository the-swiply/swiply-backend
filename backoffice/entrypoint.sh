#!/bin/sh

python manage.py makemigrations --noinput
python manage.py migrate --noinput
python manage.py add_admin
python manage.py runserver 0.0.0.0:80

exec "$@"
