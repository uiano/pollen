import React, { useState } from "react";
import { useAdminContext } from "../AdminContext/AdminContextComponent";
import { useNavigate } from "react-router-dom";
import OrderVmModal from "../Modals/OrderVmModal";
import { VMS_ARRAY } from "../../@types/types";
import { authEndpoint } from "../../Lib/http/variables";

interface INavigation {
  setVms?: React.Dispatch<React.SetStateAction<VMS_ARRAY>>;
}

function Navigation(props: INavigation) {
  const { setVms } = props;
  const [menuHamburgerOpen, setMenuHamburgerOpen] =
    React.useState<boolean>(false);
  const [open, setOpen] = useState<boolean>(false);
  const navigate = useNavigate();
  const admin = useAdminContext();

  const onMenuHamburgerClickHandler = (
    e: React.MouseEvent<HTMLButtonElement>
  ) => {
    setMenuHamburgerOpen(!menuHamburgerOpen);
  };

  function HamburgerButton() {
    return (
      <button
        type="button"
        className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
        aria-controls="mobile-menu"
        aria-expanded="false"
        onClick={onMenuHamburgerClickHandler}
      >
        <svg
          className="block h-6 w-6"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M4 6h16M4 12h16M4 18h16"
          />
        </svg>
        <svg
          className="hidden h-6 w-6"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    );
  }

  return (
    <nav className="bg-gray-800">
      <div className="max-w-7xl mx-auto px-2 sm:px-6 lg:px-8">
        <div className="relative flex items-center justify-between h-16">
          <div className="absolute inset-y-0 left-0 flex items-center sm:hidden">
            <HamburgerButton />
          </div>
          <div className="flex-1 flex items-center justify-center sm:items-stretch sm:justify-start">
            <a
              onClick={() => navigate("/")}
              className="flex text-white px-3 py-2 rounded-md text-2xl font-medium cursor-pointer"
            >
              <img
                className="h-8 w-auto mr-2"
                src={require("../../Assets/Images/pollen-logo-web.png")}
                alt="logo"
              />
            </a>
            <div className="hidden sm:block sm:ml-6">
              <div className="flex space-x-4">
                <button
                  className="mt-1.5 bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium"
                  onClick={() => setOpen(true)}
                >
                  Order
                </button>
              </div>
            </div>
            <div className="hidden sm:block sm:ml-6">
              <div className="flex space-x-4">
                <input
                  name="search"
                  type="text"
                  className="mt-1 appearance-none w-96 px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md rounded-b-md focus:outline-none sm:text-sm"
                  placeholder="Search, group name, username, ip"
                />
              </div>
            </div>

            {admin.admin && (
              <div className="mt-1.5 hidden sm:block sm:ml-6">
                <div className="flex space-x-4">
                  <a
                    onClick={() => navigate("/admin")}
                    className="bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium cursor-pointer"
                  >
                    Admin panel
                  </a>
                </div>
              </div>
            )}
          </div>

          <div className="ml-3 relative">
            <div>
              <button
                type="button"
                className="bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium cursor-pointer"
                onClick={() => {
                  window.sessionStorage.removeItem("session");
                  window.location.href = `${authEndpoint}/logout`;
                }}
              >
                Sign out
              </button>
            </div>
          </div>
        </div>
      </div>

      {menuHamburgerOpen && (
        <div className="sm:hidden" id="mobile-menu">
          <div className="px-2 pt-2 pb-3 space-y-1">
            <a
              href="/"
              className="bg-gray-900 text-white block px-3 py-2 rounded-md text-base font-medium"
              aria-current="page"
            >
              Home
            </a>
            <button
              className="bg-gray-900 text-white block px-3 py-2 rounded-md text-base font-medium cursor-pointer"
              onClick={() => setOpen(true)}
            >
              Order
            </button>
            {admin.admin && (
              <a
                onClick={() => navigate("/admin")}
                className="bg-gray-900 text-white block px-3 py-2 rounded-md text-base font-medium"
                aria-current="page"
              >
                Admin panel
              </a>
            )}
          </div>
        </div>
      )}
      <OrderVmModal open={open} setOpen={setOpen} setVms={setVms} />
    </nav>
  );
}

export default Navigation;
