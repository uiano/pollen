import { useEffect, useState } from "react";
import { SERVER_INFO, VMS_ARRAY } from "../@types/types";
import VmList from "../Components/Lists/VmList";
import Navigation from "../Components/Navigation/Navigation";
import Spinner from "../Components/Spinner";
import {
  handleErrorResponse,
  handleJSONResponse,
  get,
} from "../Lib/http/request-handler";
import OrderVmModal from "../Components/Modals/OrderVmModal";
import { useAuthProviderContext } from "../Components/Authentication/AuthProvider";

function Dashboard() {
  const [vms, setVms] = useState<VMS_ARRAY>([]);
  const [isLoadingData, setIsLoadingData] = useState<boolean>(false);
  const [open, setOpen] = useState<boolean>(false);
  const auth = useAuthProviderContext();

  useEffect(() => {
    if (auth.user && auth.user.token) {
      setIsLoadingData(true);
      get("/vms/", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setVms(r.data);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingData(false));
    }
  }, [auth.user]);

  return (
    <>
      <Navigation setVms={setVms} />
      <div className="container mx-auto mt-10">
        <div className="flex flex-col">
          <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
              {isLoadingData ? (
                <Spinner
                  w={6}
                  h={6}
                  fillColor={"black"}
                  textColor={"grey-500"}
                  textColorDark={"gray-300"}
                  label={"Loading..."}
                />
              ) : vms.length ? (
                vms.map((vm, key) => (
                  <VmList
                    auth={auth}
                    data={vm}
                    showOwner={false}
                    updateVms={(vm: SERVER_INFO) => {
                      setVms((prev: VMS_ARRAY) => {
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
                    <p>You don't have any virtual machines.</p>
                    <button
                      type="button"
                      className="block m-auto mt-4 bg-gray-900 text-white px-8 py-2 rounded-md text-sm font-medium"
                      onClick={() => setOpen(true)}
                    >
                      Order
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
      <OrderVmModal open={open} setOpen={setOpen} setVms={setVms} />
    </>
  );
}

export default Dashboard;
