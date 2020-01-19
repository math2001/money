import React, { MouseEvent } from "react";
import { months } from "./data";

const Header: React.FC<{
  month: number;
  year: number;
  move: (movetype: MoveType, amount: number) => void;
}> = ({ month, year, move }) => {
  return (
    <header className="calendar-header">
      <button onClick={() => move(MoveType.Year, -1)}>Previous Year</button>
      <button onClick={() => move(MoveType.Month, -1)}>Previous Month</button>
      {months[month]} {year}
      <button onClick={() => move(MoveType.Month, 1)}>Next Month</button>
      <button onClick={() => move(MoveType.Year, 1)}>Next Year</button>
    </header>
  );
};

export default Header;

export enum MoveType {
  Month,
  Year,
}
