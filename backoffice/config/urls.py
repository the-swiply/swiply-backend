from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    path('backoffice/admin/', admin.site.urls),
    path('backoffice/', include('apps.backoffice.urls')),
]
