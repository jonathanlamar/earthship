from typing import Any
import requests


def sendEmptyPost(url: str) -> dict[str, Any]:
    response = requests.post(
        url=url,
        data={},
        headers={"Content-type": "application/json"},
    )

    if response.status_code != 200:
        raise requests.HTTPError(f"Empty POST failed with status code {response.status_code}.")
    else:
        return response.json()


def sendGetRequestWithAccessToken(url: str, accessToken: str) -> dict[str, Any]:
    response = requests.get(
        url=url,
        headers={
            "Content-type": "application/json",
            "Authorization": f"Bearer {accessToken}",
        },
    )

    if response.status_code != 200:
        raise requests.HTTPError(f"GET failed with status code {response.status_code}")
    else:
        return response.json()
