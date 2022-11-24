import json
from pathlib import Path

from jupyter_server.serverapp import ServerApp

from .handlers import setup_handlers
from ._version import __version__

HERE = Path(__file__).parent.resolve()

with (HERE / "labextension" / "package.json").open() as fid:
    data = json.load(fid)


def _jupyter_labextension_paths():
    return [{"src": "labextension", "dest": data["name"]}]


def _jupyter_server_extension_points():
    return [{"module": "virtual_board_extension"}]


def _load_jupyter_server_extension(serverapp: ServerApp):
    """Registers the API handler to receive HTTP requests from the frontend extension.
    Parameters
    ----------
    server_app: jupyterlab.labapp.LabApp
        JupyterLab application instance
    """
    url_path = "virtual-board"
    setup_handlers(serverapp, url_path)
    serverapp.log.info(
        f"Registered jlab_ext_example extension at URL path /{url_path}"
    )


# For backward compatibility with the classical notebook
load_jupyter_server_extension = _load_jupyter_server_extension
