import { createContext, useContext, useEffect, useState } from "react";
import { handleErrorResponse } from "../../Lib/http/request-handler";
import { authEndpoint } from "../../Lib/http/variables";
import { User, AuthContextProps } from "../../@types/types";
import Spinner from "../Spinner";

const AuthProviderContext = createContext<AuthContextProps>({
  user: null,
  isLoading: false,
});

export function useAuthProviderContext() {
  return useContext(AuthProviderContext);
}

interface IAuthProviderContextComponent {
  children: JSX.Element;
}

function AuthContext(props: IAuthProviderContextComponent) {
  const { children } = props;
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const urlSearchParams = new URLSearchParams(window.location.search);
  const t = urlSearchParams.has("t")
    ? decodeURI(urlSearchParams.get("t").toString())
    : null;

  // Not signed in
  if (!t && user === null) {
    if (window.sessionStorage.getItem("session")) {
      setUser(JSON.parse(window.sessionStorage.getItem("session")));
    } else {
      window.location.href = `${authEndpoint}/provider`;
    }
  }

  useEffect(() => {
    if (t && t.length && user === null) {
      const decodedToken = atob(t);
      setLoading(true);
      fetch(`${authEndpoint}/userdata`, {
        method: "POST",
        credentials: "same-origin",
        mode: "cors",
        headers: {
          Authorization: `Bearer ${decodedToken}`,
        },
      })
        .then(async (r: Response) => {
          const data = await r.json();
          setUser({
            ...data.data,
            token: decodedToken,
          });
          setLoading(false);
          window.sessionStorage.setItem(
            "session",
            JSON.stringify({
              ...data.data,
              token: decodedToken,
            })
          );
          window.location.href = "/";
        })
        .catch(handleErrorResponse)
        .finally(() => setLoading(false));
    }
  }, [user, t]);

  return (
    <AuthProviderContext.Provider value={{ user: user, isLoading: loading }}>
      {!t && user === null ? (
        <Spinner
          w={6}
          h={6}
          fillColor={"black"}
          textColor={"grey-500"}
          textColorDark={"gray-300"}
          label={"Loading..."}
        />
      ) : (
        children
      )}
    </AuthProviderContext.Provider>
  );
}

export default AuthContext;
