import React from "react";
import Day from "./Day";
import { days, months } from "./data";

function range(n: number): number[] {
  if (n < 0) {
    throw new Error("expected n >= 0, " + n);
  }
  const arr = Array(n);
  for (let i = 0; i < n; i++) {
    arr[i] = i;
  }
  return arr;
}

function DayOfWeekHeader() {
  return (
    <tr>
      <th key="week-header">Week</th>
      {days.map(day => (
        <th key={"day-header-" + day}>{day}</th>
      ))}
    </tr>
  );
}

interface Props {
  // all 0 based
  month: number;
  year: number;
  dayOfMonth: number;
}

class Month extends React.Component<Props> {
  render() {
    return (
      <table>
        <thead>
          <DayOfWeekHeader />
        </thead>
        <tbody>
          {range(6).map(weekNumber => (
            <tr key={"week-" + weekNumber}>
              <th>{weekNumber + 1}</th>
              {range(7).map(dayNumber => (
                <Day
                  key={"day-" + dayNumber}
                  weekday={dayNumber}
                  week={weekNumber}
                  month={this.props.month}
                  year={this.props.year}
                />
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    );
  }
}

export default Month;
