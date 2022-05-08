import React, { useEffect } from "react";
import { useTabContext } from "./Tabs";

interface ITabMenu {
  children: Array<JSX.Element> | JSX.Element;
  defaultSelected?: number;
}

function TabMenu(props: ITabMenu) {
  const { defaultSelected, children } = props;
  const { setSelectedTab } = useTabContext();

  useEffect(() => {
    if (defaultSelected) {
      setSelectedTab(defaultSelected);
    }
  }, []);

  return (
    <div className="mb-4 border-b">
      <ul className="flex flex-wrap -mb-px" role="tablist">
        {children}
      </ul>
    </div>
  );
}

export default TabMenu;
