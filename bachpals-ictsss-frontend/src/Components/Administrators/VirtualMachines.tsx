import { useEffect, useState } from "react";
import {
  get,
  handleErrorResponse,
  handleJSONResponse,
} from "../../Lib/http/request-handler";
import { useAuthProviderContext } from "../Authentication/AuthProvider";
import { SERVER_INFO, VMS_ARRAY } from "../../@types/types";
import Spinner from "../Spinner";
import VmListAdmin from "../Lists/VmListAdmin";
import OrderVmModalAdmin from "../Modals/OrderVmModalAdmin";

function VirtualMachines() {
  const auth = useAuthProviderContext();
  const [servers, setServers] = useState<Array<SERVER_INFO> | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [open, setModalOpen] = useState<boolean>(false);

  useEffect(() => {
    if (!servers) {
      setIsLoading(true);
      get("/vms/all", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setServers(r.data);
          setIsLoading(false);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoading(false));
    }
  }, []);

  return (
    <div className="flex flex-col">
      <button
        type="button"
        className="flex -mt-8 mb-2 w-min bg-gray-900 text-white px-8 py-2 rounded-md text-sm font-medium"
        onClick={() => setModalOpen(true)}
      >
        Add
      </button>
      <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          {isLoading ? (
            <Spinner
              w={6}
              h={6}
              fillColor={"black"}
              textColor={"grey-500"}
              textColorDark={"gray-300"}
              label={"Loading..."}
            />
          ) : servers !== null ? (
            servers.map((vm, key) => (
              <VmListAdmin
                auth={auth}
                data={vm}
                showOwner={true}
                updateVms={(vm: SERVER_INFO) => {
                  setServers((prev: VMS_ARRAY) => {
                    if (vm === null) {
                      return [];
                    }

                    if (Array.isArray(vm)) {
                      return [...vm];
                    }

                    const newVms = prev.map((oldVm: SERVER_INFO) => {
                      if (oldVm.ServerId === vm.ServerId) {
                        oldVm = vm;
                      }

                      return oldVm;
                    });

                    return [...newVms];
                  });
                }}
                key={key}
              />
            ))
          ) : (
            <div className="flex mt-20">
              <div className="m-auto">
                <p>
                  There's no virtual machines in the system. Check your internet
                  connection, or add one now.
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
      {open && (
        <OrderVmModalAdmin
          open={open}
          setOpen={setModalOpen}
          setVms={setServers}
        />
      )}
    </div>
  );
}

export default VirtualMachines;
