import os
import json
from jupyter_server.base.handlers import APIHandler
from jupyter_server.base.handlers import JupyterHandler
from jupyter_server.utils import url_path_join
import tornado
import requests
from urllib import request
from IPython.display import HTML
import webbrowser
import re
from bs4 import BeautifulSoup
import urllib.request
import glob
import nbformat
from nbconvert.preprocessors import ExecutePreprocessor
from nbconvert import HTMLExporter


class CoursesHandler(APIHandler):
    @tornado.web.authenticated
    def get(self):
        config_data = config_file()
        try:
            courses_url = config_data['url'] + "/courses"
            response = requests.get(courses_url)
            print('response from the server', response['coursename'])
            self.finish(json.dumps(response))
        except requests.exceptions.RequestException as e:
            self.set_status(500)
            print('RequestException from the server', e)
            self.finish(json.dumps(
                {'message': "Courses get errorf. {}".format(e)}))
            return

    @tornado.web.authenticated
    def post(self):
        data = self.get_json_body()
        config_data = config_file()
        course_name = config_data['coursename']
        notebook_name = data['name']
        courses_url = config_data['url'] + "/courses"
        course_request_body = {}
        nb_request_body = {}

        with open(notebook_name) as f:
            nb = nbformat.read(f, as_version=4)
            # execute notebook
            ep = ExecutePreprocessor(timeout=-1, kernel_name='python3')
            ep.preprocess(nb)
            # export to html
            html_exporter = HTMLExporter()
            html_exporter.exclude_input = False
            html_data, resources = html_exporter.from_notebook_node(nb)

        # Get API call to fetch courses
        course_response = getCourse(
            course_name, courses_url)
        if course_response == None or not course_response.ok:
            # Post call to course
            course_request_body['coursename'] = course_name
            headers = {'Content-type': 'application/json',
                       'Accept': 'application/json'}
            try:
                course_response = requests.post(courses_url, data=json.dumps(
                    course_request_body), headers=headers, timeout=100)
                print('course_response', course_response)

            except requests.exceptions.RequestException as e:
                print('course_response error', e)
                self.set_status(500)
                self.finish(json.dumps(
                    {'message': "Courses post error. {}".format(e)}))
                return

        # Get API call to fetch Notebook by notebookName
        notebooks_endpoint = courses_url + "/"+course_name + "/notebooks"
        nb_response = getNotebookByCourse(
            notebooks_endpoint, notebook_name)

        print("notebooks_endpoint response", nb_response)
        nb_request_body['content'] = html_data
        if not nb_response.ok:
            # Post call to Notebook
            nb_request_body['notebookname'] = notebook_name
            headers = {'Content-type': 'application/json',
                       'Accept': 'application/json'}
            try:
                nb_post_response = requests.post(notebooks_endpoint, data=json.dumps(
                    nb_request_body), headers=headers, timeout=100).json()
                print('nb_post_response', nb_post_response)
                self.set_status(200)
                self.finish(tornado.escape.json_encode(nb_post_response))
                return
            except requests.exceptions.RequestException as e:
                self.set_status(500)
                self.finish(json.dumps(
                    {'message': "Notebook content post error. {}".format(e)}))
                return

        headers = {'Content-type': 'application/json',
                   'Accept': 'application/json'}

        try:
            notebook_url = notebooks_endpoint + '/' + notebook_name
            nb_put_response = requests.put(notebook_url, data=json.dumps(
                nb_request_body), headers=headers, timeout=100).json()
            print('nb_put_response')
            self.set_status(200)
            self.finish()
            return
        except requests.exceptions.RequestException as e:
            self.set_status(500)
            self.finish(json.dumps(
                {'message': "Notebook content put error. {}".format(e)}))
            return


def setup_handlers(server_app,  url_path):
    host_pattern = ".*$"
    base_url = server_app.web_app.settings["base_url"]
    route_pattern_code = url_path_join(base_url, url_path, "courses")
    handlers = [(route_pattern_code, CoursesHandler)]
    server_app.web_app.add_handlers(host_pattern, handlers)


def getCourse(course_name, base_url):
    course_url = base_url+"/" + course_name
    try:
        course_response = requests.get(course_url)
        print('getCourse', course_response)
        print('getCourse', type(course_response))
        return course_response
    except requests.exceptions.RequestException as e:
        print('getCourse RequestException', e)
        return None


def getNotebookByCourse(notebooks_endpoint, notebook_name):
    notebook_endpoint = notebooks_endpoint + "/"+notebook_name
    try:
        nb_response = requests.get(notebook_endpoint)
        print('getNotebookByCourse', nb_response)
        print('getNotebookByCourse', type(nb_response))
        return nb_response
    except requests.exceptions.RequestException as e:
        print('getNotebookByCourse RequestException', e)
        return None


def config_file():
    config_file = os.path.join(os.getcwd(), "virtual_board_config.json")
    if os.path.exists(config_file):
        f = open(config_file)
        return json.load(f)
    print("Can't Read config file")
    return None
