import {
  Administrators as AdminType,
  Administrators,
} from "../../@types/types";
import React, { ChangeEvent, useState } from "react";
import {
  del,
  handleErrorResponse,
  handleJSONResponse,
  put,
} from "../../Lib/http/request-handler";

import Button from "../Button";
import classNames from "classnames";
import { useAuthProviderContext } from "../Authentication/AuthProvider";

interface IAdminList {
  admin: Administrators;
  updateList: (items: AdminType | Array<AdminType>) => void;
}

function AdminListItem(props: IAdminList) {
  const { admin, updateList } = props;
  const [isEditable, setIsEditable] = useState<boolean>(false);
  const auth = useAuthProviderContext();

  const [fields, setFields] = useState<{
    name: string;
    user_id: string;
  } | null>({
    name: admin.Name,
    user_id: admin.UserId,
  });

  let shouldDisableSaveButton =
    fields.user_id.length === 0 || fields.name.length === 0;

  function handleItemEdit(e: React.MouseEvent<HTMLButtonElement>): void {
    e.stopPropagation();
    e.preventDefault();

    setIsEditable(!isEditable);
  }

  async function handleItemSave(
    e: React.MouseEvent<HTMLButtonElement>,
    setIsLoading?: React.Dispatch<React.SetStateAction<boolean>>
  ): Promise<void> {
    setIsLoading(true);
    await put(
      "/admin/",
      {
        name: fields.name,
        user_id: admin.UserId,
        updated_id: fields.user_id,
      },
      auth.user
    )
      .then(handleJSONResponse)
      .then((r: any) => {
        r.data !== null && updateList(r.data);
        setIsEditable(false);
      })
      .catch(handleErrorResponse)
      .finally(() => setIsLoading(false));
  }

  function handleFieldsEdit(e: ChangeEvent<HTMLInputElement>): void {
    e.stopPropagation();
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: e.target.value,
    }));
  }

  async function handleItemDelete(
    e: React.MouseEvent<HTMLButtonElement>,
    setIsLoading?: React.Dispatch<React.SetStateAction<boolean>>
  ): Promise<void> {
    e.stopPropagation();
    e.preventDefault();
    await del("/admin/", auth.user, {
      user_id: admin.UserId,
    })
      .then(handleJSONResponse)
      .then((r: any) => {
        setIsLoading(false);
        r.data !== null && updateList(r.data);
        setIsEditable(false);
      })
      .catch((e) => {
        setIsLoading(false);
      });
  }

  return (
    <tr>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center">
          <div className="ml-4">
            {isEditable ? (
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
            ) : (
              <>
                <div className="text-sm font-medium text-gray-900">
                  <p className="truncate">{admin.Name}</p>
                </div>
                <div className="text-sm text-gray-500">
                  <p className="truncate">{admin.UserId}</p>
                </div>
              </>
            )}
          </div>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
        {isEditable ? (
          <Button
            className={classNames(
              "text-indigo-600 hover:text-indigo-900",
              shouldDisableSaveButton && "text-red-600 hover:text-red-900"
            )}
            onClick={(e, setIsLoading) => handleItemSave(e, setIsLoading)}
            disabled={shouldDisableSaveButton}
          >
            <p>Save</p>
          </Button>
        ) : (
          <Button
            className="text-indigo-600 hover:text-indigo-900 cursor-pointer"
            onClick={(e) => handleItemEdit(e)}
            disableSpinner={true}
          >
            <p>Edit</p>
          </Button>
        )}
        <Button
          className="text-red-600 hover:text-red-900 cursor-pointer"
          onClick={(e, setIsLoading) => handleItemDelete(e, setIsLoading)}
          disabled={isEditable}
        >
          <p>Delete</p>
        </Button>
      </td>
    </tr>
  );
}

export default AdminListItem;
