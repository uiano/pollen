import { useEffect, useState } from "react";
import { Image } from "../../@types/types";

import { useAdminContext } from "../AdminContext/AdminContextComponent";
import {
  get,
  handleErrorResponse,
  handleJSONResponse,
} from "../../Lib/http/request-handler";
import Spinner from "../Spinner";
import AddImageModal from "../Modals/AddImageModal";
import ImagesListItem from "../Lists/ImagesListItem";
import EditImageModal from "../Modals/EditImageModal";
import { useAuthProviderContext } from "../Authentication/AuthProvider";

function Images() {
  const [imageState, setImagesState] = useState<Array<Image> | null>(null);
  const [isLoadingData, setIsLoadingData] = useState<boolean>(true);
  const [modalOpen, setModalOpen] = useState<boolean>(false);
  const [editModalOpen, setEditModalOpen] = useState<boolean>(false);
  const [selectedImage, setSelectedImage] = useState<Image | null>(null);
  const auth = useAuthProviderContext();
  const admin = useAdminContext();

  useEffect(() => {
    if (auth.user && auth.user.token && admin.admin) {
      setIsLoadingData(true);
      get("/image/", auth.user)
        .then(handleJSONResponse)
        .then((r: any) => {
          r.data !== null && setImagesState(r.data);
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
      ) : imageState === null && isLoadingData === false ? (
        <div className="flex mt-20">
          <div className="m-auto">
            <p className="text-center">
              There's no images in the system. Check your internet connection,
              or add one now.
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
                      Display name
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      Image name
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      Description
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      Published
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {imageState.map((image, index) => {
                    return (
                      <ImagesListItem
                        key={index}
                        image={image}
                        setImage={(image: Image) => {
                          setSelectedImage(image);
                          setEditModalOpen(true);
                        }}
                        updateList={(image: Image | Array<Image>) => {
                          setImagesState((prev: Array<Image>) => {
                            if (image === null) {
                              return [];
                            }

                            if (Array.isArray(image)) {
                              return [...image];
                            }

                            const newItems = prev.map((prev: Image) => {
                              if (prev.ImageId === image.ImageId) {
                                prev = image;
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
      <AddImageModal
        open={modalOpen}
        setModalOpen={setModalOpen}
        setImagesState={setImagesState}
      />
      {selectedImage !== null && (
        <EditImageModal
          open={editModalOpen}
          setModalOpen={setEditModalOpen}
          setImagesState={(image: Image | Array<Image>) => {
            setImagesState((prev: Array<Image>) => {
              if (image === null) {
                return [];
              }

              if (Array.isArray(image)) {
                return [...image];
              }

              const newItems = prev.map((prev: Image) => {
                if (prev.Id === image.Id) {
                  prev = image;
                }

                return prev;
              });

              return [...newItems];
            });
          }}
          image={selectedImage}
          setSelectedImage={setSelectedImage}
        />
      )}
    </div>
  );
}

export default Images;
