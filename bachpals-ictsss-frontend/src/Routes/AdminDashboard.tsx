import Navigation from "../Components/Navigation/Navigation";
import { useAdminContext } from "../Components/AdminContext/AdminContextComponent";
import Spinner from "../Components/Spinner";
import { TabMenu, Tab, TabsContent, Tabs } from "../Components/Tabs/index";
import Administrators from "../Components/Administrators/Administrators";
import Images from "../Components/Administrators/Images";
import VirtualMachines from "../Components/Administrators/VirtualMachines";

function AdminDashboard() {
  const admin = useAdminContext();

  return (
    <>
      <Navigation />
      <div className="container mx-auto mt-10">
        {admin.loading ? (
          <Spinner
            w={6}
            h={6}
            fillColor={"black"}
            textColor={"grey-500"}
            textColorDark={"gray-300"}
            label={"Loading..."}
          />
        ) : (
          admin.admin && (
            <Tabs>
              <TabMenu defaultSelected={1}>
                <>
                  <Tab id={1} text={"Virtual machines"} />
                  <Tab id={2} text={"Administrators"} />
                  <Tab id={3} text={"Images"} />
                </>
              </TabMenu>
              <TabsContent
                tabs={[
                  {
                    id: 1,
                    component: <VirtualMachines />,
                  },
                  {
                    id: 2,
                    component: <Administrators />,
                  },
                  {
                    id: 3,
                    component: <Images />,
                  },
                ]}
              />
            </Tabs>
          )
        )}
      </div>
    </>
  );
}

export default AdminDashboard;
