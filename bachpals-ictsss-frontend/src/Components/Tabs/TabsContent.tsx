import { useTabContext } from "./Tabs";

interface ITabsContent {
  tabs: Array<{
    id: number;
    component: JSX.Element;
  }>;
}

function TabsContent(props: ITabsContent) {
  const { tabs } = props;
  const { selectedTab } = useTabContext();

  function RenderTab() {
    if (!tabs.length) {
      return <></>;
    }

    const tab = tabs
      .map((value) => value.id === selectedTab && value.component)
      .filter(Boolean);

    return <>{tab[0]}</>;
  }

  return (
    <div className="container mx-auto mt-10">
      <RenderTab />
    </div>
  );
}

export default TabsContent;
