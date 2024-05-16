import os

from django.contrib.auth import get_user_model
from django.core.management.base import BaseCommand

class Command(BaseCommand):
    help = 'Add first superuser'

    def handle(self, *args, **options):
        username = os.getenv('ADMIN_USER')
        password = os.getenv('ADMIN_PASSWORD')

        user = get_user_model().objects.filter(username=username).first()
        if not user:
            get_user_model().objects.create_superuser(username=username, password=password, email='')