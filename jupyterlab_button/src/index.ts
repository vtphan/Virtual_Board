import { IDisposable, DisposableDelegate } from '@lumino/disposable';

import {ICommandPalette, MainAreaWidget } from '@jupyterlab/apputils';
import { Widget } from '@lumino/widgets';
import { requestAPI } from './handler'
import {
  JupyterFrontEnd,
  JupyterFrontEndPlugin
} from '@jupyterlab/application';
import { ILauncher } from '@jupyterlab/launcher';

/**
 * Initialization data for the jupyterlab_button extension.
 */

 import { ToolbarButton } from '@jupyterlab/apputils';

 import { DocumentRegistry } from '@jupyterlab/docregistry';
 
 import {
   NotebookActions,
   NotebookPanel,
   INotebookModel,
 } from '@jupyterlab/notebook';


const plugin: JupyterFrontEndPlugin<void> = {
  requires: [ICommandPalette],
  id: 'jupyterlab_wid:plugin',
  optional: [ILauncher],
  autoStart: true,
  activate: async(  app: JupyterFrontEnd, palette: ICommandPalette) => {
    app.docRegistry.addWidgetExtension('Notebook', new ButtonExtension());

  
    const content = new Widget();
    const widget = new MainAreaWidget({ content });
    widget.id = 'course-jupyterlab';
    widget.title.label = 'Courses';
    widget.title.closable = true;
    
    try {
      const data = await requestAPI<any>('/courses',{  
        method: 'GET'
    });
      console.log(data);
    } catch (reason) {
      console.error(`Error on GET /courses.\n${reason}`);
    }


const command: string = 'course:open';
app.commands.addCommand(command, {
  label: 'Display courses',
  execute: () => {
    if (!widget.isAttached) {
      // Attach the widget to the main work area if it's not there
      app.shell.add(widget, 'main');
    }
    // Activate the widget
    app.shell.activateById(widget.id);
  }
});

 // Add the command to the palette.
 palette.addItem({ command, category: 'Courses' });

    

  }
};



export class ButtonExtension
  implements DocumentRegistry.IWidgetExtension<NotebookPanel, INotebookModel>
{
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
    const clearOutput = () => {
      NotebookActions.clearAllOutputs(panel.content);
    };
    const button = new ToolbarButton({
      className: 'clear-output-button',
      label: 'VBoard',
      onClick: clearOutput,
      tooltip: 'Clear All Outputs',
    });

    panel.toolbar.insertItem(10, 'clearOutputs', button);
    return new DisposableDelegate(() => {
      button.dispose();
    });
  }
}



export default plugin;
