export type PathKind = "api" | "auth";

export const apiEndpoint = process.env.REACT_APP_API_PATH;
export const authEndpoint = process.env.REACT_APP_AUTH_PATH;

export const getPath = async (type: PathKind): Promise<string> => {
  const endpoints: Record<PathKind, string> = {
    api: apiEndpoint,
    auth: authEndpoint,
  };

  return endpoints[type];
};
