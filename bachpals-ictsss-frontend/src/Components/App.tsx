import Dashboard from "../Routes/Dahboard";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import AdminDashboard from "../Routes/AdminDashboard";
import AdminContextComponent from "./AdminContext/AdminContextComponent";
import PageNotFound from "../Routes/PageNotFound";
import Spinner from "./Spinner";
import AuthProvider from "./Authentication/AuthProvider";

function App() {
  return (
    <AuthProvider>
      <BrowserRouter basename={"/"}>
        <AdminContextComponent>
          <Routes>
            <Route index element={<Dashboard />} />
            <Route path="admin" caseSensitive element={<AdminDashboard />} />
            <Route
              caseSensitive
              path="/oauth2/redirect"
              element={
                <Spinner
                  w={6}
                  h={6}
                  fillColor={"black"}
                  textColor={"grey-500"}
                  textColorDark={"gray-300"}
                  label={"Redirecting..."}
                />
              }
            />
            <Route path="*" element={<PageNotFound />} />
          </Routes>
        </AdminContextComponent>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
