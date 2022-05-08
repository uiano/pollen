import { useEffect, useState } from "react";
import { Administrators as AdminType } from "../../@types/types";

import { useAdminContext } from "../AdminContext/AdminContextComponent";
import {
  get,
  handleErrorResponse,
  handleJSONResponse,
} from "../../Lib/http/request-handler";
import Spinner from "../Spinner";
import AdminListItem from "../Lists/AdminListItem";
import AddAdminModal from "../Modals/AddAdminModal";
import { useAuthProviderContext } from "../Authentication/AuthProvider";

function Administrators() {
  const [adminState, setAdminState] = useState<Array<AdminType> | null>(null);
  const [isLoadingData, setIsLoadingData] = useState<boolean>(true);
  const [modalOpen, setModalOpen] = useState<boolean>(false);
  const auth = useAuthProviderContext();
  const admin = useAdminContext();

  useEffect(() => {
    if (auth.user && auth.user.token && admin.admin) {
      setIsLoadingData(true);
      get("/admin/", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setAdminState(r.data);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingData(false));
    }
  }, [auth.user && auth.user.token, admin.admin]);

  return (
    <div className="flex flex-col">
      <button
        type="button"
        className="flex -mt-8 mb-2 w-min bg-gray-900 text-white px-8 py-2 rounded-md text-sm font-medium"
        onClick={() => setModalOpen(true)}
      >
        Add
      </button>
      {isLoadingData ? (
        <Spinner
          w={6}
          h={6}
          fillColor={"black"}
          textColor={"grey-500"}
          textColorDark={"gray-300"}
          label={"Loading..."}
        />
      ) : adminState === null && isLoadingData === false ? (
        <div className="flex mt-20">
          <div className="m-auto">
            <p className="text-center">
              There's no administrators in the system. Check your internet
              connection, or add one now.
            </p>
          </div>
        </div>
      ) : (
        <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
            <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      Description
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {adminState.map((admin, index) => {
                    return (
                      <AdminListItem
                        key={index}
                        admin={admin}
                        updateList={(items: AdminType | Array<AdminType>) => {
                          setAdminState((prev: Array<AdminType>) => {
                            if (items === null) {
                              return [];
                            }

                            if (Array.isArray(items)) {
                              return [...items];
                            }

                            const newItems = prev.map((prev: AdminType) => {
                              if (prev.UserId === items.UserId) {
                                prev = items;
                              }

                              return prev;
                            });

                            return [...newItems];
                          });
                        }}
                      />
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
      <AddAdminModal
        open={modalOpen}
        setModalOpen={setModalOpen}
        setAdminState={setAdminState}
      />
    </div>
  );
}

export default Administrators;
