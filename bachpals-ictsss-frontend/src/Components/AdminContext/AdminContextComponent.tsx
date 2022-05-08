import { createContext, useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { get, handleErrorResponse } from "../../Lib/http/request-handler";
import { Administrators } from "../../@types/types";
import { useAuthProviderContext } from "../Authentication/AuthProvider";

type AdminContext = {
  admin: Administrators;
  loading: boolean;
};

const AdminContext = createContext<AdminContext>({
  admin: null,
  loading: false,
});

export function useAdminContext() {
  return useContext(AdminContext);
}

interface IAdminContextComponent {
  children: JSX.Element;
}

function AdminContextComponent(props: IAdminContextComponent) {
  const { children } = props;
  const auth = useAuthProviderContext();
  const navigate = useNavigate();
  const [adminState, setAdminState] = useState<Administrators | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (auth.user && auth.user.token) {
      const mail =
        auth && auth.user.email.includes("@student.uia.no")
          ? auth.user.email.replaceAll("@student.uia.no", "@uia.no")
          : auth.user.email;

      setLoading(true);
      get(`/admin/${mail}`, auth.user)
        .then(async (r: Response) => {
          if (r.ok) {
            const data = await r.json();
            setAdminState(data.data);
            setLoading(false);
          } else {
            setLoading(false);
            navigate("/", { replace: true });
          }
        })
        .catch(handleErrorResponse)
        .finally(() => setLoading(false));
    }
  }, [auth.user]);

  return (
    <AdminContext.Provider value={{ admin: adminState, loading: loading }}>
      {children}
    </AdminContext.Provider>
  );
}

export default AdminContextComponent;
