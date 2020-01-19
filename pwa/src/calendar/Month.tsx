import React from "react";
import Day from "./Day";
import { days, months, Entry } from "./data";

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
  entriesOf(date: Date): Entry[] {
    return [];
  }

  render() {
    const reference = new Date();
    reference.setFullYear(this.props.year);
    reference.setMonth(this.props.month);
    reference.setDate(1);

    const firstDayOfWeek = reference.getDay();

    return (
      <table>
        <thead>
          <DayOfWeekHeader />
        </thead>
        <tbody>
          {range(6).map((week: number) => (
            <tr key={"week-" + week}>
              <th>{week + 1}</th>
              {range(7).map((weekDay: number) => {
                const dayMonth = week * 7 + weekDay - firstDayOfWeek + 1;
                const copy = new Date(reference.getTime());
                copy.setDate(dayMonth);
                // this will manage the changes in months/year
                // (if dayMonth is negative, or 32 for example)
                return (
                  <Day
                    key={"week-" + week + "-day-" + weekDay}
                    dayMonth={copy.getDate()}
                    month={copy.getMonth() + 1}
                    year={copy.getFullYear()}
                    dim={copy.getMonth() !== this.props.month}
                    entries={this.entriesOf(copy)}
                  />
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    );
  }
}

export default Month;
