import Navigation from "../Components/Navigation/Navigation";

function PageNotFound() {
  return (
    <>
      <Navigation />
      <div className="flex items-center justify-center w-screen mt-60">
        <div className="px-4">
          <div className="lg:gap-4 lg:flex">
            <div className="flex flex-col items-center justify-center">
              <h1 className="font-bold text-blue-600 text-9xl">404</h1>
              <p className="mb-2 text-2xl font-bold text-center text-gray-800 md:text-3xl">
                <span className="text-red-500">Oops!</span> Page not found
              </p>
              <p className="mb-8 text-center text-gray-500 md:text-lg">
                The page you’re looking for does not exist.
              </p>
              <a
                href="/"
                className="bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium cursor-pointer"
              >
                Go home
              </a>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

export default PageNotFound;
