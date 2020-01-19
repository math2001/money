import React from "react";
import { months } from "./data";

const Header: React.FC<{ month: number; year: number }> = ({ month, year }) => {
  return (
    <header className="calendar-header">
      {months[month]} {year}
    </header>
  );
};

export default Header;
