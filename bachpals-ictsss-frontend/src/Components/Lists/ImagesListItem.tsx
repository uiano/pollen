import { Image } from "../../@types/types";
import React from "react";
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
  image: Image;
  updateList: (items: Image | Array<Image>) => void;
  setImage: (image: Image) => void;
}

function ImagesListItem(props: IAdminList) {
  const { image, setImage, updateList } = props;
  const auth = useAuthProviderContext();

  function handleItemEdit(e: React.MouseEvent<HTMLButtonElement>): void {
    e.stopPropagation();
    e.preventDefault();
    setImage(image);
  }

  async function handleImagePublish(
    image: Image,
    setIsLoading: React.Dispatch<React.SetStateAction<boolean>>
  ) {
    await put(
      "/image/",
      {
        id: image.Id,
        published: image.Published === "true" ? "false" : "true",
        image_id: image.ImageId,
        image_name: image.ImageName,
        image_description: image.ImageDescription,
        image_display_name: image.ImageDisplayName,
        image_config: image.ImageConfig,
        image_read_root_password: image.ImageReadRootPassword,
      },
      auth.user
    )
      .then(handleJSONResponse)
      .then((r: any) => {
        r.data !== null && updateList(r.data);
      })
      .catch(handleErrorResponse)
      .finally(() => setIsLoading(false));
  }

  async function handleItemDelete(
    e: React.MouseEvent<HTMLButtonElement>,
    setIsLoading?: React.Dispatch<React.SetStateAction<boolean>>
  ): Promise<void> {
    e.stopPropagation();
    e.preventDefault();

    await del("/image/", auth.user, {
      id: image.Id,
    })
      .then(handleJSONResponse)
      .then((r: any) => {
        r.data !== null && updateList(r.data);
      })
      .catch(handleErrorResponse)
      .finally(() => setIsLoading(false));
  }

  return (
    <tr>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center">
          <div className="ml-4">
            <div className="text-sm font-medium text-gray-900">
              <p className="truncate">{image.ImageDisplayName}</p>
            </div>
          </div>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center">
          <div className="ml-4">
            <div className="text-sm text-gray-500">
              <p className="truncate">{image.ImageName}</p>
            </div>
          </div>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center">
          <div className="ml-4">
            <div className="text-sm text-gray-500">
              <p className="truncate">{image.ImageDescription}</p>
            </div>
          </div>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center">
          <div className="ml-4">
            <div className="text-sm text-gray-500">
              <Button
                className={classNames(
                  image.Published === "true"
                    ? "border-red-300 bg-white text-white-700 hover:bg-red-50"
                    : "border-green-300 bg-white text-gray-700 hover:bg-green-50"
                )}
                onClick={(e, setIsLoading) =>
                  handleImagePublish(image, setIsLoading)
                }
              >
                <p>{image.Published === "true" ? "Disable" : "Enable"}</p>
              </Button>
            </div>
          </div>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
        <Button
          className="text-indigo-600 hover:text-indigo-900 cursor-pointer"
          onClick={(e) => handleItemEdit(e)}
          disableSpinner={true}
        >
          <p>Edit</p>
        </Button>
        <Button
          className="text-red-600 hover:text-red-900 cursor-pointer"
          onClick={(e, setIsLoading) => handleItemDelete(e, setIsLoading)}
        >
          <p>Delete</p>
        </Button>
      </td>
    </tr>
  );
}

export default ImagesListItem;
