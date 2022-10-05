from distutils.command.config import config
import os
import json

from notebook.base.handlers import APIHandler
from notebook.utils import url_path_join

import tornado
from tornado.web import StaticFileHandler
import requests

server_url = "http://localhost:8080/api/v1"

class CoursesHandler(APIHandler):
    # The following decorator should be present on all verb methods (head, get, post,
    # patch, put, delete, options) to ensure only authorized user can request the
    # Jupyter server

    @tornado.web.authenticated
    def get(self):
        url =  server_url + "/courses"
        try:
            resp = requests.get(url,timeout=5).json()
        except requests.exceptions.RequestException as e:
            self.set_status(500)
            self.finish(json.dumps({'message': "Go server Error. {}".format(e)}))
            return

        self.finish(json.dumps({'response': resp}))


def setup_handlers(web_app):
    host_pattern = ".*$"
    base_url = web_app.settings["base_url"]

    # Prepend the base_url so that it works in a JupyterHub setting
    route_pattern = url_path_join(base_url, "virtual_board", "courses")
    handlers = [(route_pattern, CoursesHandler)]
    web_app.add_handlers(host_pattern, handlers)
