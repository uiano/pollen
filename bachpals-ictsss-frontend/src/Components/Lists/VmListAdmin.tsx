import classNames from "classnames";
import React from "react";
import {
  AuthContextProps,
  SERVER_INFO,
  VMS_RESPONSE,
} from "../../@types/types";
import { useTimeout } from "../../Lib/hooks/useTimeout";
import {
  del,
  get,
  handleErrorResponse,
  handleJSONResponse,
  post,
} from "../../Lib/http/request-handler";
import Button from "../Button";

interface IListItem {
  auth: AuthContextProps;
  data: SERVER_INFO;
  updateVms: (vm: SERVER_INFO | Array<SERVER_INFO>) => void;
  showOwner: boolean;
}

function ListItem(props: IListItem) {
  const { data, auth, updateVms, showOwner } = props;

  const [expandListItem, setExpandListItem] = React.useState<boolean>(false);
  const [hasReceivedPassword, setHasReceivedPassword] =
    React.useState<string>("");

  const handleOnClick = (e: React.MouseEvent<HTMLDivElement>) => {
    setExpandListItem(!expandListItem);
  };

  useTimeout(async () => {
    if (auth.user) {
      get(`/vms/${data.ServerId}/status`, auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && updateVms(r.data);
        })
        .catch(handleErrorResponse);
    }
  }, 10000);

  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-2 shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
      <div
        className="px-4 py-4 grid grid-cols-6 gap-0 cursor-pointer"
        onClick={handleOnClick}
      >
        <p className="mt-2 px-4 text-sm text-gray-500">
          {data.ServerStatus === "ACTIVE" ? (
            <span className="inline-flex mt-0.5">
              <span className="flex relative h-3 w-3 top-0 right-0 mt-1 mr-1">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-300 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-green-400"></span>
              </span>
            </span>
          ) : (
            <span className="inline-flex mt-0.5">
              <span className="flex relative h-3 w-3 top-0 right-0 mt-1 mr-1">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-300 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-red-400"></span>
              </span>
            </span>
          )}
        </p>
        <p className="truncate pt-2 px-4 max-w-2xl text-sm text-gray-500">
          {"Owner: " + (data.UserId ? data.UserId : "Unknown")}
        </p>
        <p className="truncate pt-2 px-4 max-w-2xl text-sm text-gray-500">
          {"Name: " + (data.ServerName ? data.ServerName : "Unknown")}
        </p>
        <p className="truncate pt-2 px-4 max-w-2xl text-sm text-gray-500">
          {"IP: " + (data.ServerIp ? data.ServerIp : "Unknown")}
        </p>
        <p className="truncate pt-2 px-4 max-w-2xl text-sm text-gray-500">
          {"System: " +
            (data.ImageDisplayName ? data.ImageDisplayName : "Unknown")}
        </p>
        <span className="px-4 max-w-2xl text-sm text-gray-500">
          <Button
            className={classNames(
              "w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium focus:outline-none m:mt-0 sm:ml-3 sm:w-auto sm:text-sm",
              "border-red-300 bg-white text-white-700 hover:bg-red-50"
            )}
            onClick={(e, setIsLoading) => {
              del(`/vms/${data.ServerId}`, auth.user)
                .then(handleJSONResponse)
                .then((r: VMS_RESPONSE) => {
                  setIsLoading(false);
                  updateVms(r.data);
                })
                .catch(() => setIsLoading(false));
            }}
          >
            <p>Delete</p>
          </Button>
        </span>
      </div>
      {expandListItem && (
        <div className="border-t border-gray-200">
          <dl className="bg-gray-50 px-4 py-5 grid grid-cols-2 gap-4">
            <div className="row-span-3">
              {data.GroupMembers.length && (
                <>
                  <p>Group members</p>
                  <ol className="list-disc ml-4">
                    {data.GroupMembers.map((member: string, key: number) => {
                      return <li key={key}>{member}</li>;
                    })}
                  </ol>
                </>
              )}
            </div>

            <div className="col-span-0">
              <p>Actions</p>
              <div className="flex flex-row mt-2 -ml-3 flex-wrap">
                <Button
                  className={
                    data.ServerStatus === "ACTIVE"
                      ? "z-0 border-red-300 bg-white text-white-700 hover:bg-red-50"
                      : "z-0 border-green-300 bg-white text-gray-700 hover:bg-green-50"
                  }
                  onClick={(e, setIsLoading) => {
                    if (data.ServerStatus === "ACTIVE") {
                      post(`/vms/${data.ServerId}/stop`, auth.user)
                        .then(handleJSONResponse)
                        .then((r: VMS_RESPONSE) => {
                          updateVms(r.data);
                          setIsLoading(false);
                        })
                        .catch(handleErrorResponse)
                        .finally(() => setIsLoading(false));
                    } else if (data.ServerStatus === "SHUTOFF") {
                      post(`/vms/${data.ServerId}/start`, auth.user)
                        .then(handleJSONResponse)
                        .then((r: VMS_RESPONSE) => {
                          updateVms(r.data);
                          setIsLoading(false);
                        })
                        .catch(handleErrorResponse)
                        .finally(() => setIsLoading(false));
                    }
                  }}
                >
                  <p>{data.ServerStatus === "ACTIVE" ? "Stop" : "Start"}</p>
                </Button>
                <Button
                  className={classNames(
                    "z-0 w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium focus:outline-none m:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-75",
                    data.ServerStatus === "ACTIVE" && "hover:bg-gray-200"
                  )}
                  disabled={data.ServerStatus !== "ACTIVE"}
                  onClick={(
                    e: React.MouseEvent<HTMLButtonElement>,
                    setIsLoading: (s: boolean) => void
                  ) => {
                    post(`/vms/${data.ServerId}/respawn`, auth.user)
                      .then(handleJSONResponse)
                      .then((r: VMS_RESPONSE) => {
                        updateVms(r.data);
                        setIsLoading(false);
                      })
                      .catch(handleErrorResponse)
                      .finally(() => setIsLoading(false));
                  }}
                >
                  <p>Respawn</p>
                </Button>
                <Button
                  className={classNames(
                    "z-0 w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium focus:outline-none m:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-75",
                    data.ServerStatus === "ACTIVE" && "hover:bg-gray-200"
                  )}
                  disabled={data.ServerStatus !== "ACTIVE"}
                  onClick={(
                    e: React.MouseEvent<HTMLButtonElement>,
                    setIsLoading: (s: boolean) => void
                  ) => {
                    post(`/vms/${data.ServerId}/reboot`, auth.user)
                      .then(handleJSONResponse)
                      .then((r: VMS_RESPONSE) => {
                        updateVms(r.data);
                        setIsLoading(false);
                      })
                      .catch(handleErrorResponse)
                      .finally(() => setIsLoading(false));
                  }}
                >
                  <p>Reboot</p>
                </Button>
                <Button
                  className={classNames(
                    "z-0 w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium focus:outline-none m:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-75",
                    data.ServerStatus === "ACTIVE" && "hover:bg-gray-200"
                  )}
                  disabled={data.ServerStatus !== "ACTIVE"}
                  onClick={(e, setIsLoading) => {
                    get(`/vms/${data.ServerId}/console`, auth.user)
                      .then(handleJSONResponse)
                      .then((r: any) => {
                        setIsLoading(false);

                        r.data !== null &&
                          r.data.url &&
                          window.open(r.data.url, "_blank");
                      })
                      .catch(handleErrorResponse)
                      .finally(() => setIsLoading(false));
                  }}
                >
                  <p>Terminal</p>
                </Button>
                {data.ImageReadRootPassword && (
                  <>
                    <Button
                      className={classNames(
                        "z-0 w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium focus:outline-none m:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-75",
                        data.ServerStatus === "ACTIVE" && "hover:bg-gray-200"
                      )}
                      disabled={data.ServerStatus !== "ACTIVE"}
                      onClick={(e, setIsLoading) => {
                        get(`/vms/${data.ServerId}/password`, auth.user)
                          .then(handleJSONResponse)
                          .then((r: any) => {
                            setIsLoading(false);
                            r.data !== null && setHasReceivedPassword(r.data);
                          })
                          .catch(handleErrorResponse)
                          .finally(() => setIsLoading(false));
                      }}
                    >
                      <p>Get password</p>
                    </Button>
                    {hasReceivedPassword.length > 0 && (
                      <p className="py-2 ml-3 bg-gray-100 block w-1/4 shadow-sm sm:text-sm rounded-md p-2">
                        {hasReceivedPassword}
                      </p>
                    )}
                  </>
                )}
              </div>
            </div>
          </dl>
        </div>
      )}
    </div>
  );
}

export default ListItem;
