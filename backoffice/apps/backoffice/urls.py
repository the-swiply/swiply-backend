from django.urls import path

from apps.backoffice import views
from django.conf.urls.static import static

from config import settings

urlpatterns = [
    path('profiles/', views.ProfileListView.as_view(), name='profiles'),
    path('profile/<str:slug>/', views.ProfileDetailView.as_view(), name='profile_detail'),
    path('graph/', views.graph),
    path('block-user/<str:slug>/', views.BlockUserView.as_view(), name='block-user'),
    path('unblock-user/<str:slug>/', views.UnblockUserView.as_view(), name='unblock-user'),
    path('', views.AuthorizationView.as_view(), name='authorization'),
    path('logoout/', views.logout_view, name='logout'),
]

urlpatterns += static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
