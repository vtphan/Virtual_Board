import { JupyterFrontEnd } from "@jupyterlab/application";
import { DocumentRegistry } from "@jupyterlab/docregistry";
import { NotebookPanel, INotebookModel } from "@jupyterlab/notebook";
import { IDisposable, DisposableDelegate } from "@lumino/disposable";
import { ToolbarButton, Dialog, showDialog } from "@jupyterlab/apputils";
import { requestAPI } from "./handler";
export class CoursesExtension
  implements DocumentRegistry.IWidgetExtension<NotebookPanel, INotebookModel> {
  constructor(app: JupyterFrontEnd) {
    this.app = app;
  }
  readonly app: JupyterFrontEnd;
  /**
   * Create a new extension for the notebook panel widget.
   *
   * @param panel Notebook panel
   * @param context Notebook context
   * @returns Disposable on the added button
   */

  createNew(
    panel: NotebookPanel,
    context: DocumentRegistry.IContext<INotebookModel>
  ): IDisposable {
    const courses = async () => {
      if (context.sessionContext.propertyChanged) {
        let fileName = context.sessionContext.name;
        console.log("notebook context", context);
        showDialog({
          title: "",
          body: "Notebook was posted to Virtual Board.",
          buttons: [Dialog.okButton({ label: "Ok" })],
        });
        const dataToSend = { name: fileName };
        try {
          const response = await requestAPI<any>("courses", {
            body: JSON.stringify(dataToSend),
            method: "POST",
          });

          console.log("Backend response", response);
        } catch (reason) {
          console.error(
            `Error on POST /virtual-board/courses ${dataToSend}.\n${reason}`
          );
        }
      }
    };
    const button = new ToolbarButton({
      className: "view-courses-button",
      label: "Send to VBoard",
      onClick: courses,
      tooltip: "View all courses",
    });

    panel.toolbar.insertItem(10, "DisplayCourses", button);
    return new DisposableDelegate(() => {
      button.dispose();
    });
  }
}
