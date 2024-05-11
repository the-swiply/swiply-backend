import requests
from posixpath import join as urljoin


class ProfileClient:
    def __init__(self, url, s2s_token):
        self.url = url
        self.base_headers = {
            's2s-authorization': 'Bearer ' + s2s_token
        }

    def get_profiles(self):
        resp = requests.get(urljoin(self.url, 'v1/profile'), params={'updatedAfter': '2024-01-01T00:00:00Z'},
                            headers=self.base_headers)
        resp.raise_for_status()

        return resp.json()['profiles']

    def get_profile_by_id(self, id):
        resp = requests.get(urljoin(self.url, 'v1/profile', id), headers=self.base_headers)
        resp.raise_for_status()

        return resp.json()['user_profile']

    def get_interests(self):
        resp = requests.get(urljoin(self.url, 'v1/interests'), headers=self.base_headers)
        resp.raise_for_status()

        return resp.json()['interests']

    def get_photos_by_profile_id(self, id):
        resp = requests.get(urljoin(self.url, 'v1/photo', id), headers=self.base_headers)

        if resp.status_code == 404:
            return []

        resp.raise_for_status()

        return resp.json()['photos']

    def change_availability(self, id, is_blocked):
        resp = requests.post(urljoin(self.url, 'v1/profile/change-availability'), json={
            'id': id,
            'is_blocked': is_blocked
        }, headers=self.base_headers)

        resp.raise_for_status()
