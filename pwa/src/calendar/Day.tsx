import React, { MouseEvent } from "react";
import { Entry, DayDate } from "./data";

interface Props {
  // note that these are all 0 based
  // dayMonth is the day of the months (new Date().getDate())
  dayMonth: number;
  month: number;
  year: number;

  dim: boolean;
  selected: boolean;

  onClick: (date: DayDate, extend: boolean) => void;

  entries: Iterable<Entry>;
}

class Day extends React.Component<Props> {
  constructor(props: Props) {
    super(props);

    this.onClick = this.onClick.bind(this);
  }

  onClick(event: MouseEvent) {
    this.props.onClick(
      new DayDate(this.props.year, this.props.month, this.props.dayMonth),
      event.shiftKey,
    );
  }

  render() {
    return (
      <td
        className={
          "day" +
          (this.props.dim ? " day-dim" : "") +
          (this.props.selected ? " day-selected" : "")
        }
        onClick={this.onClick}
      >
        {this.props.dayMonth}
        <ul className="day-entries">
          {Array.from(this.props.entries, entry => (
            <li key={entry.id}>{entry.name}</li>
          ))}
        </ul>
      </td>
    );
  }
}

export default Day;
