import React from "react";
import { FiSearch } from "react-icons/fi";

interface SearchBarProps {
  searchValue: string;
  onSearchChange: (value: string) => void;
  onSearchSubmit: () => void;
}

const SearchBar: React.FC<SearchBarProps> = ({
  searchValue,
  onSearchChange,
  onSearchSubmit,
}) => {
  return (
    <div className="relative flex h-[61px] w-[355px] items-center rounded-full bg-white px-2 py-2 shadow-xl shadow-shadow-500 dark:bg-navy-800 dark:shadow-none">
      <button
        onClick={onSearchSubmit}
        className="flex h-full items-center justify-center rounded-l-full bg-lightPrimary text-navy-700 dark:bg-navy-900 dark:text-white px-4 text-xl"
      >
        <FiSearch className="h-4 w-4 text-gray-400 dark:text-white" />
      </button>
      <input
        type="text"
        value={searchValue}
        onChange={(e) => onSearchChange(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            onSearchSubmit();
          }
        }}
        placeholder="Search..."
        className="block h-full w-full rounded-r-full bg-lightPrimary text-sm font-medium text-navy-700 outline-none placeholder:!text-gray-400 dark:bg-navy-900 dark:text-white dark:placeholder:!text-white"
      />
    </div>
  );
};

export default SearchBar;
