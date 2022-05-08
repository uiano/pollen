/* This example requires Tailwind CSS v2.0+ */
import React, { useState } from "react";
import classNames from "classnames";
import Button from "../Button";

interface IModal {
  cancelCallback: (
    setIsLoading?: React.Dispatch<React.SetStateAction<boolean>>
  ) => void;
  acceptCallback: (
    setIsLoading?: React.Dispatch<React.SetStateAction<boolean>>
  ) => void;
  cancelButtonText: string;
  acceptButtonText: string;
  cancelButtonStyles?: string;
  acceptButtonStyles?: string;
  open: boolean;
  disableAcceptButton?: boolean;
  disableCancelButton?: boolean;
  children: JSX.Element | Array<JSX.Element>;
}

function ModalBase(props: IModal) {
  const {
    acceptButtonText,
    cancelButtonText,
    cancelCallback,
    children,
    acceptCallback,
    acceptButtonStyles,
    cancelButtonStyles,
    open,
    disableAcceptButton = false,
    disableCancelButton = false,
  } = props;

  return (
    <>
      {open && (
        <div className="fixed z-10 inset-0 overflow-y-auto" role="dialog">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

            <span
              className="hidden sm:inline-block sm:align-middle sm:h-screen"
              aria-hidden="true"
            >
              &#8203;
            </span>

            <div className="relative inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                {children}
              </div>
              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <Button
                  className={classNames(
                    "w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 text-base font-medium text-white sm:ml-3 sm:w-auto sm:text-sm",
                    acceptButtonStyles && acceptButtonStyles
                  )}
                  onClick={(e, setIsLoading) => {
                    e.stopPropagation();
                    e.preventDefault();
                    acceptCallback(setIsLoading);
                  }}
                  disabled={disableAcceptButton}
                >
                  <p>{acceptButtonText}</p>
                </Button>
                <Button
                  className={classNames(
                    "mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm",
                    cancelButtonStyles && cancelButtonStyles
                  )}
                  onClick={(e, setIsLoading) => {
                    e.stopPropagation();
                    e.preventDefault();
                    cancelCallback(setIsLoading);
                  }}
                  disabled={disableCancelButton}
                >
                  <p>{cancelButtonText}</p>
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

export default ModalBase;
