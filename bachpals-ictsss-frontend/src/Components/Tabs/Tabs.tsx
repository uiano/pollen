import React from "react";
import { createContext, useContext, useState } from "react";
import { ITabContext } from "../../@types/types";

const DEFAULT_SELECT_TAB = 0;

const TabContext = createContext<ITabContext>({
  selectedTab: DEFAULT_SELECT_TAB,
  setSelectedTab: null,
});

export function useTabContext() {
  return useContext(TabContext);
}

interface ITabs {
  children: Array<JSX.Element>;
}

function Tabs(props: ITabs) {
  const { children } = props;
  const [selectedTab, setSelectedTab] = useState<number>(0);

  return (
    <TabContext.Provider
      value={{ selectedTab: selectedTab, setSelectedTab: setSelectedTab }}
    >
      {children}
    </TabContext.Provider>
  );
}

export default Tabs;
