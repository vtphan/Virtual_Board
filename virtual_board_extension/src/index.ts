import {
  JupyterFrontEnd,
  JupyterFrontEndPlugin,
} from "@jupyterlab/application";
import { ILauncher } from "@jupyterlab/launcher";
import { DocumentRegistry } from "@jupyterlab/docregistry";
import { INotebookModel } from "@jupyterlab/notebook";
import { CoursesExtension } from "./course/course";

/**
 * Initialization data for the virtual_board extension.
 */
const plugin: JupyterFrontEndPlugin<void> = {
  id: "virtual_board:plugin",
  optional: [ILauncher],
  autoStart: true,
  activate: async (
    app: JupyterFrontEnd,
    context: DocumentRegistry.IContext<INotebookModel>
  ) => {
    console.log("Virtual_board is activated!");
    const virtual_board = new CoursesExtension(app);
    app.docRegistry.addWidgetExtension("Notebook", virtual_board);
  },
};

export default plugin;
