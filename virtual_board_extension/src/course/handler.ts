import { URLExt } from "@jupyterlab/coreutils";

import { ServerConnection } from "@jupyterlab/services";

export async function requestAPI<T>(
  endPoint = "",
  init: RequestInit = {}
): Promise<T> {
  // Make request to Jupyter API
  const settings = ServerConnection.makeSettings();
  const requestUrl = URLExt.join(settings.baseUrl, "virtual-board", endPoint);

  let response: Response;
  try {
    response = await ServerConnection.makeRequest(requestUrl, init, settings);
  } catch (error) {
    throw new ServerConnection.NetworkError(
      error as ServerConnection.NetworkError
    );
  }
  const data = await response.json();

  if (!response.ok) {
    throw new ServerConnection.ResponseError(response, data.message);
  }
  return data;
}
