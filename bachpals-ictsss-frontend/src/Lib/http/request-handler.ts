import { getPath } from "./variables";
import { User } from "../../@types/types";

export async function post(
  path: string,
  auth: User,
  obj?: Record<string, unknown>
): Promise<Response> {
  const serverPath = await getPath("api");
  const body = JSON.stringify(obj);

  return fetch(serverPath + path, {
    method: "POST",
    credentials: "same-origin",
    mode: "cors",
    headers: {
      Authorization: `Bearer ${auth.token}`,
    },
    body,
  }).catch((err: Response) => {
    throw err;
  });
}

export async function get(
  path: string,
  auth: User,
  obj?: Record<string, string | number | boolean>
): Promise<Response> {
  const params = new URLSearchParams(JSON.stringify(obj));

  const apiPath = await getPath("api");
  const u =
    typeof obj !== "undefined"
      ? apiPath + path + `?${params.toString()}`
      : apiPath + path;

  return fetch(u, {
    method: "GET",
    credentials: "same-origin",
    mode: "cors",
    headers: {
      Authorization: `Bearer ${auth.token}`,
    },
  }).catch((err: Response) => {
    throw err;
  });
}

export async function put(
  path: string,
  obj: Record<string, unknown>,
  auth: User
): Promise<Response> {
  const body = JSON.stringify(obj);

  const apiPath = await getPath("api");

  return fetch(apiPath + path, {
    method: "PUT",
    credentials: "same-origin",
    mode: "cors",
    headers: {
      Authorization: `Bearer ${auth.token}`,
    },
    body,
  }).catch((err: Response) => {
    throw err;
  });
}

export async function del(
  path: string,
  auth: User,
  obj?: Record<string, unknown>
): Promise<Response> {
  const serverPath = await getPath("api");
  const body = JSON.stringify(obj);

  return fetch(serverPath + path, {
    method: "DELETE",
    credentials: "same-origin",
    mode: "cors",
    headers: {
      Authorization: `Bearer ${auth.token}`,
    },
    body,
  }).catch((err: Response) => {
    throw err;
  });
}

export function handleJSONResponse<R>(response: Response): Promise<R> {
  if (!response.ok) {
    throw response;
  }

  return response.json() as Promise<R>;
}

export function handleEmptyResponse(response: Response): boolean {
  if (!response.ok) {
    throw response;
  }

  return true;
}

export async function handleErrorResponse(
  error: Response | Error
): Promise<any> {
  const errorText = await getErrorText(error);
}

async function getErrorText(error: Response | Error): Promise<string> {
  let errorMessage: string;

  if ("message" in error) {
    errorMessage = getMessageFromError(error);
  } else {
    errorMessage = await getMessageFromResponse(error);
  }

  if (!errorMessage) errorMessage = "unknown_error";

  return errorMessage;
}

function getMessageFromError(error: Error): string {
  return error.message;
}

async function getMessageFromResponse(response: Response): Promise<string> {
  if (!response?.json) return undefined;
  return response
    .json()
    .then((obj) => obj?.error)
    .catch(() => "could_not_parse_json");
}
