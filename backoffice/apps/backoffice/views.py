import os

from django.shortcuts import render, redirect
from django.views import View
from django.contrib.auth import authenticate, login, logout

from apps.backoffice.clients.profile import ProfileClient
from apps.backoffice.forms import AuthForm
from apps.backoffice.utils.auth import check_auth
from config import settings


class ProfileListView(View):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)

        self.__profile_client = ProfileClient(settings.PROFILE_URL, os.getenv('PROFILE_S2S_TOKEN'))

    def get(self, request):
        check_auth(request)

        profiles = self.__profile_client.get_profiles()

        return render(request, 'profile_list.html', context={'profiles': profiles})


class ProfileDetailView(View):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)

        self.__profile_client = ProfileClient(settings.PROFILE_URL, os.getenv('PROFILE_S2S_TOKEN'))

    def get(self, request, slug):
        check_auth(request)

        profile = self.__profile_client.get_profile_by_id(slug)
        photos = self.__profile_client.get_photos_by_profile_id(slug)

        profile['birth_day'] = profile['birth_day'][:-10]

        return render(request, 'profile_detail.html', {'user': profile, 'photos': photos})


def logout_view(request):
    logout(request)
    return redirect('/')


class AuthorizationView(View):
    def get(self, request):
        form = AuthForm()

        return render(request, 'authorization.html', context={'form': form})

    def post(self, request):
        form = AuthForm(request.POST)

        if form.is_valid():
            username = form.cleaned_data['username']
            password = form.cleaned_data['password']
            user = authenticate(username=username, password=password)

            if user:
                login(request, user)
                return redirect('profiles/')
            else:
                form.add_error('__all__', 'Неверные данные!')

        return render(request, 'authorization.html', context={'form': form})


def graph(request):
    return render(request, 'graph.html')
