import ModalBase from "./ModalBase";
import React, { ChangeEvent, useState } from "react";

import { handleJSONResponse, post } from "../../Lib/http/request-handler";
import { Administrators as AdminType } from "../../@types/types";
import { useAuthProviderContext } from "../Authentication/AuthProvider";

interface IAddAdminModal {
  open: boolean;
  setModalOpen: React.Dispatch<React.SetStateAction<boolean>>;
  setAdminState: React.Dispatch<React.SetStateAction<Array<AdminType> | null>>;
}

function AddAdminModal(props: IAddAdminModal) {
  const { open, setModalOpen, setAdminState } = props;
  const auth = useAuthProviderContext();
  const [fields, setFields] = useState<{
    name: string;
    user_id: string;
  } | null>({
    name: "",
    user_id: "",
  });

  function handleFieldsEdit(e: ChangeEvent<HTMLInputElement>): void {
    e.stopPropagation();
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: e.target.value,
    }));
  }

  return (
    <ModalBase
      cancelCallback={() => {
        setModalOpen(false);
      }}
      acceptCallback={async (setIsLoading) => {
        await post("/admin/", auth.user, fields)
          .then(handleJSONResponse)
          .then((r: any) => {
            setIsLoading(false);
            r.data !== null && setAdminState(r.data);
            setFields({
              name: "",
              user_id: "",
            });
            setModalOpen(false);
          })
          .catch((e) => {
            setFields({
              name: "",
              user_id: "",
            });
            setIsLoading(false);
          });
      }}
      cancelButtonText={"Cancel"}
      acceptButtonText={"Add"}
      acceptButtonStyles={
        "border-none text-white bg-blue-700 hover:bg-blue-800 focus:ring-0 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
      }
      cancelButtonStyles={"focus:ring-0"}
      open={open}
    >
      <div className="sm:flex sm:items-start">
        <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
          <h3
            className="text-lg leading-6 font-medium text-gray-900"
            id="modal-title"
          >
            Add new administrator
          </h3>
          <div className="mt-2">
            <p className="text-sm text-gray-500">
              Keep in mind that the email-address has to be @uia.no.
            </p>
          </div>
          <div className="mt-4">
            <div className="grid grid-cols-6 gap-6">
              <div className="col-span-6 sm:col-span-3">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Name
                </label>
                <input
                  type="text"
                  name="Name"
                  id="name"
                  autoComplete="given-name"
                  className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                  onChange={handleFieldsEdit}
                  value={fields.name}
                  placeholder={"Name"}
                />
              </div>
              <div className="col-span-6 sm:col-span-3">
                <label
                  htmlFor="Email"
                  className="block text-sm font-medium text-gray-700"
                >
                  Email
                </label>
                <input
                  type="email"
                  name="Email"
                  id="user_id"
                  autoComplete="email"
                  className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                  onChange={handleFieldsEdit}
                  value={fields.user_id}
                  placeholder={"Email"}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </ModalBase>
  );
}

export default AddAdminModal;
