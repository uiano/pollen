import ModalBase from "./ModalBase";
import React, { ChangeEvent, useEffect, useState } from "react";
import {
  get,
  handleErrorResponse,
  handleJSONResponse,
  post,
} from "../../Lib/http/request-handler";
import { Image, ImageConfig, ServerImage } from "../../@types/types";
import Spinner from "../Spinner";
import { useAuthProviderContext } from "../Authentication/AuthProvider";

interface IAddImageModal {
  open: boolean;
  setModalOpen: React.Dispatch<React.SetStateAction<boolean>>;
  setImagesState: React.Dispatch<React.SetStateAction<Array<Image> | null>>;
}

function AddImageModal(props: IAddImageModal) {
  const { open, setModalOpen, setImagesState } = props;
  const [serverImages, setServerImages] = useState<Array<ServerImage> | null>(
    null
  );
  const [imageConfigs, setImageConfigs] = useState<Array<ImageConfig> | null>(
    null
  );
  const [isLoadingDataImages, setIsLoadingDataImages] = useState<boolean>(true);
  const [isLoadingDataConfigs, setIsLoadingDataConfigs] =
    useState<boolean>(true);
  const auth = useAuthProviderContext();
  const [fields, setFields] = useState<{
    published: string;
    image_id: string;
    image_name: string;
    image_description: string;
    image_display_name: string;
    image_config: string;
    image_read_root_password: boolean;
  } | null>({
    published: "false",
    image_id: "",
    image_name: "",
    image_description: "",
    image_display_name: "",
    image_config: "",
    image_read_root_password: false,
  });

  function handleFieldsEdit(
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ): void {
    e.stopPropagation();
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: e.target.value,
    }));
  }

  function handleCheckboxEdit(e: ChangeEvent<HTMLInputElement>): void {
    e.stopPropagation();
    setFields((prevState) => ({
      ...prevState,
      [e.target.id]: e.target.checked,
    }));
  }

  function handleOnSelectImage(e: ChangeEvent<HTMLSelectElement>): void {
    const selectedOption = e.currentTarget.options[e.currentTarget.value];

    setFields((prevState) => ({
      ...prevState,
      image_id: selectedOption.dataset.imageId,
      image_name: selectedOption.dataset.imageName,
    }));
  }

  function handleOnSelectConfig(e: ChangeEvent<HTMLSelectElement>): void {
    const selectedOption =
      e.currentTarget.options[e.currentTarget.selectedIndex];
    setFields((prevState) => ({
      ...prevState,
      image_config: selectedOption.dataset.imageConfig,
    }));
  }

  useEffect(() => {
    if (open && auth.user) {
      setIsLoadingDataImages(true);
      get("/image/server", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setServerImages(r.data);
          setIsLoadingDataImages(false);

          // Initialize select box default values for image selection
          setFields((prevState) => ({
            ...prevState,
            image_id: r.data[0].ImageId,
            image_name: r.data[0].Name,
          }));
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingDataImages(false));
    }
  }, [open, auth.user]);

  useEffect(() => {
    if (open && auth.user) {
      setIsLoadingDataConfigs(true);
      get("/image/config", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setImageConfigs(r.data);
          setIsLoadingDataConfigs(false);
        })
        .catch(handleErrorResponse)
        .finally(() => setIsLoadingDataConfigs(false));
    }
  }, [open, auth.user]);

  return (
    <ModalBase
      cancelCallback={() => {
        setModalOpen(false);
      }}
      acceptCallback={async (setIsLoading) => {
        await post("/image/", auth.user, fields)
          .then(handleJSONResponse)
          .then((r: any) => {
            setIsLoading(false);
            r.data !== null && setImagesState(r.data);
            setFields({
              published: "false",
              image_id: "",
              image_name: "",
              image_description: "",
              image_display_name: "",
              image_config: "",
              image_read_root_password: false,
            });
            setModalOpen(false);
          })
          .catch(handleErrorResponse)
          .finally(() => setIsLoading(false));
      }}
      cancelButtonText={"Cancel"}
      acceptButtonText={"Add"}
      acceptButtonStyles={
        "border-none text-white bg-blue-700 hover:bg-blue-800 focus:ring-0 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700"
      }
      cancelButtonStyles={"focus:ring-0"}
      open={open}
      disableAcceptButton={isLoadingDataImages || isLoadingDataConfigs}
    >
      <div className="sm:flex sm:items-start">
        <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
          <h3
            className="text-lg leading-6 font-medium text-gray-900"
            id="modal-title"
          >
            Add new image
          </h3>
          <div className="mt-2">
            <p className="text-sm text-gray-500">
              Adding images allows you to chose what operating systems users can
              use.
            </p>
          </div>
          <div className="mt-4">
            <div className="grid grid-cols-6 gap-6">
              <div className="col-span-6 sm:col-span-3">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Display name
                </label>
                <input
                  type="text"
                  name="displayName"
                  id="image_display_name"
                  className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                  onChange={handleFieldsEdit}
                  value={fields.image_display_name}
                  placeholder={"Display name"}
                  required={true}
                />
              </div>
              <div className="col-span-6 sm:col-span-3">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Image
                </label>

                {isLoadingDataImages ? (
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
                    id="image_id"
                    className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                    onChange={handleOnSelectImage}
                    defaultValue={0}
                    required={true}
                  >
                    {serverImages.map((image: ServerImage, key) => {
                      return (
                        <option
                          key={key}
                          data-image-id={image.ImageId}
                          data-image-name={image.Name}
                          value={key}
                        >
                          {image.Name}
                        </option>
                      );
                    })}
                  </select>
                ) : (
                  <p>Failed to load</p>
                )}
              </div>
              <div className="col-span-6 sm:col-span-3">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Config
                </label>

                {isLoadingDataConfigs ? (
                  <Spinner
                    w={4}
                    h={4}
                    fillColor={"black"}
                    textColor={"grey-500"}
                    textColorDark={"gray-300"}
                    parentClassNames={"mt-2"}
                  />
                ) : imageConfigs ? (
                  <select
                    id="image_config"
                    className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                    onChange={handleOnSelectConfig}
                    defaultValue={""}
                    required={true}
                  >
                    <option key={0} value={""}>
                      {"None"}
                    </option>
                    {imageConfigs.map((config: ImageConfig, key) => {
                      return (
                        <option
                          key={key + 1}
                          value={config}
                          data-image-config={config}
                        >
                          {config}
                        </option>
                      );
                    })}
                  </select>
                ) : (
                  <p>Failed to load</p>
                )}
              </div>
              <div className="flex col-span-4 sm:col-span-4">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Allow reading root password?:
                  <p className="block text-sm font-small text-gray-400">
                    (For images not using ldap)
                  </p>
                </label>
                <input
                  type="checkbox"
                  name="imageReadRootPassword"
                  id="image_read_root_password"
                  className="mt-1 ml-2 block shadow-sm sm:text-sm border-gray-300 rounded-md"
                  onChange={handleCheckboxEdit}
                  checked={fields.image_read_root_password}
                  required={true}
                />
              </div>
              <div className="col-span-6 sm:col-span-6">
                <label
                  htmlFor="Name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Description
                </label>
                <textarea
                  required={true}
                  name="displayDescription"
                  id="image_description"
                  className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md resize-none"
                  onChange={handleFieldsEdit}
                  value={fields.image_description}
                  placeholder={"Description"}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </ModalBase>
  );
}

export default AddImageModal;
