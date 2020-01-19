import React from "react";
import { Entry } from "./data";

interface Props {
  // note that these are all 0 based
  // dayMonth is the day of the months (new Date().getDate())
  dayMonth: number;
  month: number;
  year: number;

  dim: boolean;

  entries: Entry[];
}

class Day extends React.Component<Props> {
  render() {
    return (
      <td className={"day" + (this.props.dim ? " day-dim" : "")}>
        {this.props.dayMonth}
        <ul className="day-entries">
          {this.props.entries.map(entry => (
            <li key={entry.id}>{entry.name}</li>
          ))}
        </ul>
      </td>
    );
  }
}

export default Day;
