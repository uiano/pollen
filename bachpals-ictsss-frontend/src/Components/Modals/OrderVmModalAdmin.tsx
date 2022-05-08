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
import { useAdminContext } from "../AdminContext/AdminContextComponent";

interface IOrderVmModal {
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  setVms: React.Dispatch<React.SetStateAction<VMS_ARRAY>>;
}

function OrderVmModal(props: IOrderVmModal) {
  const { open, setOpen, setVms } = props;
  const auth = useAuthProviderContext();
  const admin = useAdminContext();
  const [courses, setCourses] = useState(null);
  const [isLoadingCourses, setIsLoadingCourses] = useState<boolean>(false);

  const [courseGroups, setCourseGroups] = useState(null);
  const [isLoadingCourseGroups, setIsLoadingCourseGroups] =
    useState<boolean>(false);

  const [courseStudents, setCourseStudents] = useState(null);
  const [isLoadingCourseStudents, setIsLoadingCourseStudents] =
    useState<boolean>(false);

  const [isLoadingImages, setIsLoadingImages] = useState<boolean>(false);
  const [serverImages, setServerImages] =
    useState<Array<SelectServerImage> | null>(null);

  const [fields, setFields] = useState<{
    course: string;
    course_id: string;
    group_id: string;
    group: string;
    personal: string;
    group_name: string;
    group_students: Array<string>;
    server_image: string;
    everyone: string;
    include_ta: string;
    include_teacher: string;
  } | null>({
    course: "",
    course_id: "",
    group_id: "",
    group: "",
    personal: "",
    group_name: "",
    group_students: [],
    server_image: "",
    everyone: "",
    include_ta: "false",
    include_teacher: "false",
  });

  function handleRadioButtonsEdit(e: ChangeEvent<HTMLInputElement>): void {
    e.stopPropagation();

    if (e.target.id === "group") {
      setFields((prevState) => ({
        ...prevState,
        [e.target.id]: e.target.value,
        personal: "",
        group_students: [],
        group_name: "",
        group_id: "",
        everyone: "",
        include_ta: "false",
        include_teacher: "false",
        server_image: "",
      }));
    } else if (e.target.id === "personal") {
      setFields((prevState) => ({
        ...prevState,
        [e.target.id]: e.target.value,
        group: "",
        group_students: [],
        group_name: "",
        group_id: "",
        everyone: "",
        include_ta: "false",
        include_teacher: "false",
        server_image: "",
      }));
    } else if (e.target.id === "everyone") {
      setFields((prevState) => ({
        ...prevState,
        [e.target.id]: e.target.value,
        group: "",
        group_students: [],
        group_name: "",
        group_id: "",
        personal: "",
        include_ta: "false",
        include_teacher: "false",
        server_image: "",
      }));
    }
  }

  function handleCheckbox(e: ChangeEvent<HTMLInputElement>): void {
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: String(e.target.checked),
    }));
  }

  function handleCourseSelect(e: ChangeEvent<HTMLSelectElement>): void {
    if (e.target.value.includes(",")) {
      const values = e.target.value.split(",");

      setFields((prevState) => ({
        ...prevState,
        [e.target.id]: values[0].replaceAll(" ", "-"),
        course_id: values[1],
        group_students: [],
        group_name: "",
        group_id: "",
        server_image: "",
        everyone: "",
        include_ta: "false",
        include_teacher: "false",
      }));
      return;
    }

    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: "",
      course_id: "",
      group: "",
      group_students: [],
      course: "",
      group_name: "",
      group_id: "",
      personal: "",
      server_image: "",
      everyone: "",
      include_ta: "false",
      include_teacher: "false",
    }));
  }

  function handleSelect(e: ChangeEvent<HTMLSelectElement>): void {
    if (e.target.value.length === 0) {
      setFields((prevState) => ({
        ...prevState,
        [e.target.id]: "",
        group_students: [],
        group_id: "",
        everyone: "",
        include_ta: "false",
        include_teacher: "false",
      }));
      return;
    }

    if (e.target.id === "group_name") {
      const value = e.target.value.split(",");
      const groupName = `grp-${value[0]}`;

      setFields((prevState) => ({
        ...prevState,
        [e.target.id]: groupName,
        group_id: value[1],
        group_students: [],
        everyone: "",
        include_ta: "false",
        include_teacher: "false",
      }));

      return;
    } else if (e.target.id === "student_id") {
      setFields((prevState) => ({
        ...prevState,
        group_id: "",
        group_name: "",
        group_students: [e.target.value],
        everyone: "",
        include_ta: "false",
        include_teacher: "false",
      }));
      return;
    }
  }

  useEffect(() => {
    if (!courses) {
      setIsLoadingCourses(true);
      get("/courses/", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setCourses(r.data);
          setIsLoadingCourses(false);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingCourses(false));
    }
  }, []);

  useEffect(() => {
    if (!fields.course_id || !fields.course) {
      return;
    }

    if (fields.group === "1" && !fields.group_id) {
      setIsLoadingCourseGroups(true);
      get(`/courses/${fields.course_id}/groups`, auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setCourseGroups(r.data);
          setFields((prevState) => ({
            ...prevState,
            group_name: "",
          }));
          setIsLoadingCourseGroups(false);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingCourseGroups(false));
    }

    if (fields.personal === "1") {
      setIsLoadingCourseStudents(true);
      get(`/courses/${fields.course_id}/users`, auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setCourseStudents(r.data);
          setIsLoadingCourseStudents(false);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingCourseStudents(false));
    }

    if (fields.group.length > 0 && fields.group_id) {
      get(`/courses/groups/${fields.group_id}/users`, auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null &&
            !r.data.status &&
            setFields((prev) => {
              return {
                ...prev,
                group_students: r.data.map((u) => u.login_id).filter(String),
              };
            });
        })
        .catch(handleErrorResponse);
    }
  }, [fields.group, fields.personal, fields.course_id, fields.group_id]);

  useEffect(() => {
    if (fields.group_students.length > 0 || fields.everyone) {
      if (auth.user && auth.user.token && admin.admin) {
        setIsLoadingImages(true);
        get("/image/published", auth.user)
          .then(handleJSONResponse)
          .then((r: any) => {
            r.data !== null && setServerImages(r.data);
          })
          .catch(handleErrorResponse)
          .finally(() => setIsLoadingImages(false));
      }
    } else {
      setFields((prev) => {
        return {
          ...prev,
          server_image: "",
        };
      });
    }
  }, [fields.group_students.length, fields.everyone]);

  function handleServerImageSelect(e: ChangeEvent<HTMLSelectElement>): void {
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: e.target.value,
    }));
  }

  return (
    <ModalBase
      cancelCallback={() => {
        setFields((prevState) => ({
          ...prevState,
          course_id: "",
          group: "",
          group_students: [],
          course: "",
          group_name: "",
          group_id: "",
          personal: "",
          server_image: "",
          everyone: "",
          include_ta: "false",
          include_teacher: "false",
        }));
        setOpen(false);
      }}
      acceptCallback={(setIsLoading) => {
        post(`/vms/canvas${fields.everyone === "1" ? "/all" : ""}`, auth.user, {
          server_name: fields.course,
          group_name: fields.group_students.length > 1 ? fields.group_name : "",
          users: fields.group_students,
          server_image: fields.server_image,
          everyone: fields.everyone,
          include_ta: fields.include_ta,
          include_teacher: fields.include_teacher,
          course_code: fields.course_id,
        })
          .then(handleJSONResponse)
          .then((r: any) => {
            r.data !== null && setVms(r.data);
            setFields((prevState) => ({
              ...prevState,
              course_id: "",
              group: "",
              group_students: [],
              course: "",
              group_name: "",
              group_id: "",
              personal: "",
              server_image: "",
              everyone: "",
              include_ta: "false",
              include_teacher: "false",
            }));
            setOpen(false);
          })
          .catch(handleErrorResponse)
          .finally(() => setIsLoading(false));
      }}
      cancelButtonText={"Close"}
      acceptButtonText={"Order"}
      acceptButtonStyles={
        "border-none text-white bg-blue-700 hover:bg-blue-800 focus:ring-0 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
      }
      open={open}
      disableAcceptButton={
        fields.course.length === 0
          ? true
          : fields.course_id.length === 0
          ? true
          : fields.group_students.length === 0 && fields.everyone === ""
          ? true
          : fields.server_image.length === 0
          ? true
          : false
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
              <div className="mt-4 col-span-4 sm:col-span-4">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Select course
                </label>

                {isLoadingCourses ? (
                  <Spinner
                    w={4}
                    h={4}
                    fillColor={"black"}
                    textColor={"grey-500"}
                    textColorDark={"gray-300"}
                    parentClassNames={"mt-2"}
                  />
                ) : courses ? (
                  <select
                    id="course"
                    className="block w-full mt-1 shadow-sm sm:text-sm border-gray-300 rounded-md"
                    onChange={handleCourseSelect}
                    required={true}
                  >
                    <option key={0} value={""}>
                      {"None"}
                    </option>
                    {courses.map((course, key) => {
                      return (
                        <option
                          key={key}
                          value={[course.course_code, course.id]}
                        >
                          {course.course_code}
                        </option>
                      );
                    })}
                  </select>
                ) : (
                  <p>Failed to load</p>
                )}
              </div>
              <div className="col-span-4 sm:col-span-4">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Personal or group?
                </label>
                <div>
                  <div className="form-check">
                    <input
                      className="form-check-input appearance-none rounded-full h-4 w-4 border border-gray-300 bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer"
                      type="radio"
                      name="group"
                      id="group"
                      value="1"
                      onChange={handleRadioButtonsEdit}
                      disabled={!fields.course_id || !fields.course}
                    />
                    <label
                      className="form-check-label inline-block text-gray-800"
                      htmlFor="group"
                    >
                      Group
                    </label>
                  </div>
                  <div className="form-check">
                    <input
                      className="form-check-input appearance-none rounded-full h-4 w-4 border border-gray-300 bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer"
                      type="radio"
                      name="group"
                      id="personal"
                      value="1"
                      onChange={handleRadioButtonsEdit}
                      disabled={!fields.course_id || !fields.course}
                    />
                    <label
                      className="form-check-label inline-block text-gray-800"
                      htmlFor="personal"
                    >
                      User
                    </label>
                  </div>
                  <div className="form-check">
                    <input
                      className="form-check-input appearance-none rounded-full h-4 w-4 border border-gray-300 bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer"
                      type="radio"
                      name="group"
                      id="everyone"
                      value="1"
                      onChange={handleRadioButtonsEdit}
                      disabled={!fields.course_id || !fields.course}
                    />
                    <label
                      className="form-check-label inline-block text-gray-800"
                      htmlFor="everyone"
                    >
                      Everyone
                    </label>
                  </div>
                </div>
              </div>
              {fields.group === "1" && (
                <div className="mt-4 col-span-4 sm:col-span-4">
                  <label
                    htmlFor="Name"
                    className="block text-sm font-medium text-gray-700"
                  >
                    Select group
                  </label>

                  {isLoadingCourseGroups ? (
                    <Spinner
                      w={4}
                      h={4}
                      fillColor={"black"}
                      textColor={"grey-500"}
                      textColorDark={"gray-300"}
                      parentClassNames={"mt-2"}
                    />
                  ) : courseGroups ? (
                    <select
                      id="group_name"
                      className="block w-full mt-1 shadow-sm sm:text-sm border-gray-300 rounded-md"
                      onChange={handleSelect}
                      required={true}
                    >
                      <option key={0} value={""}>
                        {"None"}
                      </option>
                      {courseGroups.map((group, key) => {
                        return (
                          <option key={key} value={[key, group.id]}>
                            {group.name}
                          </option>
                        );
                      })}
                    </select>
                  ) : (
                    <p>Failed to load</p>
                  )}
                </div>
              )}
              {fields.personal === "1" && (
                <div className="mt-4 col-span-4 sm:col-span-4">
                  <label
                    htmlFor="Name"
                    className="block text-sm font-medium text-gray-700"
                  >
                    Select student
                  </label>

                  {isLoadingCourseStudents ? (
                    <Spinner
                      w={4}
                      h={4}
                      fillColor={"black"}
                      textColor={"grey-500"}
                      textColorDark={"gray-300"}
                      parentClassNames={"mt-2"}
                    />
                  ) : courseStudents ? (
                    <select
                      id="student_id"
                      className="block w-full mt-1 shadow-sm sm:text-sm border-gray-300 rounded-md"
                      onChange={handleSelect}
                      required={true}
                    >
                      <option key={0} value={""}>
                        {"None"}
                      </option>
                      {courseStudents.map((student, key) => {
                        return (
                          <option key={key} value={student.login_id}>
                            {student.name}
                          </option>
                        );
                      })}
                    </select>
                  ) : (
                    <p>Failed to load</p>
                  )}
                </div>
              )}
              {fields.everyone === "1" && (
                <div className="mt-4 col-span-4 sm:col-span-4">
                  <label
                    htmlFor="Name"
                    className="block text-sm font-medium text-gray-700"
                  >
                    Options
                  </label>
                  <div className="form-check">
                    <input
                      className="form-check-input appearance-none rounded-full h-4 w-4 border border-gray-300 bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer"
                      type="checkbox"
                      name="include_teacher"
                      id="include_teacher"
                      value={String(fields.include_teacher)}
                      onChange={handleCheckbox}
                      disabled={!fields.course_id || !fields.course}
                    />
                    <label
                      className="form-check-label inline-block text-gray-800"
                      htmlFor="include_teacher"
                    >
                      Include teachers
                    </label>
                  </div>
                  <div className="form-check">
                    <input
                      className="form-check-input appearance-none rounded-full h-4 w-4 border border-gray-300 bg-white checked:bg-blue-600 checked:border-blue-600 focus:outline-none transition duration-200 mt-1 align-top bg-no-repeat bg-center bg-contain float-left mr-2 cursor-pointer"
                      type="checkbox"
                      name="include_ta"
                      id="include_ta"
                      value={String(fields.include_ta)}
                      onChange={handleCheckbox}
                      disabled={!fields.course_id || !fields.course}
                    />
                    <label
                      className="form-check-label inline-block text-gray-800"
                      htmlFor="include_ta"
                    >
                      Include teaching assistants
                    </label>
                  </div>
                </div>
              )}
              {(fields.group_students.length > 0 ||
                fields.everyone === "1") && (
                <div className="mt-4 col-span-4 sm:col-span-4">
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
              )}
            </div>
          </div>
        </div>
      </div>
    </ModalBase>
  );
}

export default OrderVmModal;
