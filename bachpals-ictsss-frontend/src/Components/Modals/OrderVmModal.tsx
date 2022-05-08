import ModalBase from "./ModalBase";
import React, { ChangeEvent, SyntheticEvent, useEffect, useState } from "react";

import classNames from "classnames";
import Button from "../Button";
import {
  get,
  handleErrorResponse,
  handleJSONResponse,
  post,
} from "../../Lib/http/request-handler";
import Spinner from "../Spinner";
import { SelectServerImage, VMS_ARRAY } from "../../@types/types";
import { useAuthProviderContext } from "../Authentication/AuthProvider";

interface IOrderVmModal {
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  setVms: React.Dispatch<React.SetStateAction<VMS_ARRAY>>;
}

function OrderVmModal(props: IOrderVmModal) {
  const { open, setOpen, setVms } = props;
  const auth = useAuthProviderContext();
  const [isLoadingImages, setIsLoadingImages] = useState<boolean>(false);
  const [serverImages, setServerImages] =
    useState<Array<SelectServerImage> | null>(null);
  const [fields, setFields] = useState<{
    course_code: string;
    group: boolean;
    for_course: boolean;
    group_name: string;
    group_members: Array<string>;
    server_image: string;
    member_name_text: string;
  } | null>({
    course_code: "",
    for_course: true,
    group: false,
    group_name: "",
    group_members: [],
    server_image: "",
    member_name_text: "",
  });

  function handleFieldsEdit(e: ChangeEvent<HTMLInputElement>): void {
    e.stopPropagation();
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]:
        e.target.type === "checkbox" ? e.target.checked : e.target.value,
    }));
  }

  function handleServerImageSelect(e: ChangeEvent<HTMLSelectElement>): void {
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: e.target.value,
    }));
  }

  function handleAddGroupMember(e: SyntheticEvent<HTMLButtonElement>): void {
    if (fields.member_name_text.length === 0) {
      return;
    }

    if (!fields.member_name_text.includes("@uia.no")) {
      return;
    }

    if (fields.group_members.includes(fields.member_name_text)) {
      return;
    }

    const sub = fields.member_name_text.split("@");

    if (auth.user.email.includes(sub[0])) {
      return;
    }

    fields.group_members.push(fields.member_name_text);
    setFields((prevState) => ({
      ...prevState,
      group_members: fields.group_members,
      member_name_text: "",
    }));
  }

  useEffect(() => {
    if (open) {
      setIsLoadingImages(true);
      get("/image/published", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setServerImages(r.data);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingImages(false));
    }
  }, [open]);

  async function orderVm(
    setIsLoading: React.Dispatch<React.SetStateAction<boolean>>
  ) {
    await post("/vms/", auth.user, {
      server_name: fields.course_code,
      group_name: fields.group ? fields.group_name : "",
      users: fields.group ? fields.group_members : [],
      server_image: fields.server_image,
    })
      .then(handleJSONResponse)
      .then((r: any) => {
        r.data !== null && setVms(r.data);
        setFields({
          course_code: "",
          for_course: true,
          group: false,
          group_name: "",
          group_members: [],
          server_image: "",
          member_name_text: "",
        });
        setOpen(false);
      })
      .catch(handleErrorResponse)
      .finally(() => setIsLoading(false));
  }

  return (
    <ModalBase
      cancelCallback={() => {
        setFields({
          course_code: "",
          for_course: true,
          group: false,
          group_name: "",
          group_members: [],
          server_image: "",
          member_name_text: "",
        });
        setOpen(false);
      }}
      acceptCallback={(setIsLoading) => orderVm(setIsLoading)}
      cancelButtonText={"Close"}
      acceptButtonText={"Order"}
      acceptButtonStyles={
        "border-none text-white bg-blue-700 hover:bg-blue-800 focus:ring-0 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
      }
      open={open}
      disableAcceptButton={
        fields.course_code.length === 0 ||
        fields.server_image.length === 0 ||
        (fields.group === true && fields.group_members.length === 0) ||
        (fields.for_course === true &&
          fields.group === true &&
          fields.group_name.length === 0)
      }
    >
      <div className="sm:flex sm:items-start">
        <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
          <h3
            className="text-lg leading-6 font-medium text-gray-900"
            id="modal-title"
          >
            Order virtual machine
          </h3>
          <div className="mt-2">
            <p className="text-sm text-gray-500"></p>
          </div>
          <div className="mt-4">
            <div className="grid grid-cols-4 gap-6">
              <div className="flex col-span-4 sm:col-span-4">
                <label
                  htmlFor="group"
                  className="block text-sm font-medium text-gray-700"
                >
                  Order for course:
                </label>
                <input
                  type="checkbox"
                  name="for_course"
                  id="for_course"
                  className="mt-0.5 ml-2 block shadow-sm sm:text-sm border-gray-300 rounded-md"
                  onChange={handleFieldsEdit}
                  value={String(fields.for_course)}
                  checked={fields.for_course}
                />
              </div>
              {!fields.for_course ? (
                <div className="col-span-4 sm:col-span-4">
                  <label
                    htmlFor="Name"
                    className="block text-sm font-medium text-gray-700"
                  >
                    Project name
                  </label>
                  <input
                    type="text"
                    name="courseCode"
                    id="course_code"
                    className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                    onChange={handleFieldsEdit}
                    value={fields.course_code}
                    placeholder={"Example: FunProject"}
                  />
                </div>
              ) : (
                <div className="col-span-4 sm:col-span-4">
                  <label
                    htmlFor="Name"
                    className="block text-sm font-medium text-gray-700"
                  >
                    Course code
                  </label>
                  <input
                    type="text"
                    name="courseCode"
                    id="course_code"
                    className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                    onChange={handleFieldsEdit}
                    value={fields.course_code}
                    placeholder={"Example: IKT201"}
                  />
                </div>
              )}
              <div className="flex col-span-4 sm:col-span-4">
                <label
                  htmlFor="group"
                  className="block text-sm font-medium text-gray-700"
                >
                  Order for a group:
                </label>
                <input
                  type="checkbox"
                  name="group"
                  id="group"
                  className="mt-1 ml-2 block shadow-sm sm:text-sm border-gray-300 rounded-md"
                  onChange={handleFieldsEdit}
                  value={String(fields.group)}
                  checked={fields.group}
                />
              </div>
              {fields.group ? (
                <>
                  {fields.for_course && (
                    <div className="col-span-4 sm:col-span-4">
                      <label
                        htmlFor="Name"
                        className="block text-sm font-medium text-gray-700"
                      >
                        Group name
                      </label>
                      <div className="flex">
                        <input
                          type="text"
                          name="group_name"
                          id="group_name"
                          className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                          onChange={handleFieldsEdit}
                          placeholder={"Example: grp-1"}
                        />
                      </div>
                    </div>
                  )}
                  <div className="col-span-4 sm:col-span-4">
                    <label
                      htmlFor="Name"
                      className="block text-sm font-medium text-gray-700"
                    >
                      Add group members
                    </label>
                    <div className="flex">
                      <input
                        type="email"
                        name="member_name_text"
                        id="member_name_text"
                        className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                        onChange={handleFieldsEdit}
                        value={fields.member_name_text}
                        placeholder={"Example: user@uia.no"}
                      />
                      <Button
                        className={classNames(
                          "mt-1 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                        )}
                        onClick={(e, _) => handleAddGroupMember(e)}
                        disableSpinner={true}
                      >
                        <p>{"Add"}</p>
                      </Button>
                    </div>
                    <div className="w-full bg-white rounded-lg">
                      <ul className="divide-y-2 divide-gray-100">
                        {fields.group_members &&
                        fields.group_members.length > 0 ? (
                          fields.group_members.map(
                            (member: string, key: number) => (
                              <li
                                className="flex justify-between p-3 sm:text-sm"
                                key={key}
                              >
                                {member}
                                <button
                                  onClick={() => {
                                    setFields((prevState) => ({
                                      ...prevState,
                                      group_members:
                                        fields.group_members.filter(
                                          (item) => item !== member
                                        ),
                                    }));
                                  }}
                                >
                                  <svg
                                    className="h-4 w-4 text-red-500"
                                    width="24"
                                    height="24"
                                    viewBox="0 0 24 24"
                                    strokeWidth="2"
                                    stroke="currentColor"
                                    fill="none"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                  >
                                    <path stroke="none" d="M0 0h24v24H0z" />
                                    <line x1="18" y1="6" x2="6" y2="18" />
                                    <line x1="6" y1="6" x2="18" y2="18" />
                                  </svg>
                                </button>
                              </li>
                            )
                          )
                        ) : (
                          <p className="text-sm font-medium text-gray-700">
                            {"No group members"}
                          </p>
                        )}
                      </ul>
                    </div>
                  </div>
                </>
              ) : (
                <></>
              )}
            </div>
            <div className="mt-4 col-span-6 sm:col-span-6">
              <label
                htmlFor="Name"
                className="block text-sm font-medium text-gray-700"
              >
                Image
              </label>

              {isLoadingImages ? (
                <Spinner
                  w={4}
                  h={4}
                  fillColor={"black"}
                  textColor={"grey-500"}
                  textColorDark={"gray-300"}
                  parentClassNames={"mt-2"}
                />
              ) : serverImages ? (
                <select
                  id="server_image"
                  className="block w-full mt-1 shadow-sm sm:text-sm border-gray-300 rounded-md"
                  onChange={handleServerImageSelect}
                  required={true}
                >
                  <option key={0} value={""}>
                    {"None"}
                  </option>
                  {serverImages.map((image: SelectServerImage, key) => {
                    return (
                      <option key={key} value={image.ImageId}>
                        {image.ImageDisplayName}
                      </option>
                    );
                  })}
                </select>
              ) : (
                <p>Failed to load</p>
              )}
            </div>
          </div>
        </div>
      </div>
    </ModalBase>
  );
}

export default OrderVmModal;
