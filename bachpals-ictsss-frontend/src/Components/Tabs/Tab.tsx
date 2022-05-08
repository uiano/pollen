import classNames from "classnames";
import { useTabContext } from "./Tabs";

interface ITab {
  id: number;
  text: string;
}

function Tab(props: ITab) {
  const { id, text } = props;
  const { selectedTab, setSelectedTab } = useTabContext();

  function selectTab() {
    if (id === selectedTab) {
      return false;
    }

    setSelectedTab(id);
  }

  function isTabSelected() {
    if (id === selectedTab) {
      return true;
    }
    return false;
  }

  const style = {
    selected:
      "inline-block py-4 px-4 text-sm font-medium text-center text-blue-800 rounded-t-lg border-b-2 border-blue-800 active dark:text-blue-800 dark:border-blue-800",
    unselected:
      "inline-block py-4 px-4 text-sm font-medium text-center text-gray-500 rounded-t-lg border-b-2 border-transparent hover:text-gray-600 dark:text-gray-400",
  };

  return (
    <li className="mr-2" role="presentation" onClick={() => selectTab()}>
      <button
        className={classNames(
          "inline-block py-4 px-4 text-sm font-medium text-center",
          isTabSelected() ? style.selected : style.unselected
        )}
      >
        {text}
      </button>
    </li>
  );
}

export default Tab;
