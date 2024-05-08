from django import forms


class AuthForm(forms.Form):
    username = forms.CharField(max_length=30, label='Имя пользователя')
    password = forms.CharField(widget=forms.PasswordInput, label='Пароль')
