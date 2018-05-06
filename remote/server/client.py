import json as json_lib
import requests

class Client():
    """
    A wrapper for RESTful requests made to a server.
    """

    def __init__(self, url, port):
        self._url = url
        self._port = port

    @property
    def url(self):
        return self._url

    @property
    def port(self):
        return self._port

    def get_full_url(self, resource):
        full_url = "{}:{}{}".format(self.url, str(self.port), resource)
        if "http://" not in full_url:
            full_url = "http://{}".format(full_url)
        return full_url

    def get(self, resource, default=None):
        full_url = self.get_full_url(resource)
        try:
            r = requests.get(url = full_url)
            data = json_lib.dumps(ast.literal_eval(r.text))
            return json_lib.loads(data)
        except Exception as e:
            return default

    def post(self, resource, json, default=None):
        full_url = self.get_full_url(resource)
        try:
            r = requests.post(
                url = full_url,
                json = json)
            data = json_lib.dumps(ast.literal_eval(r.text))
            return json_lib.loads(data)
        except Exception as e:
            return default
