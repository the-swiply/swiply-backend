from django.core.exceptions import PermissionDenied


def check_auth(request):
    """
    Check if user is authenticated.
    """

    if not request.user.is_authenticated:
        raise PermissionDenied()
