import React from "react";
import Day from "./Day";
import { days, Entry, DayDate } from "./data";
import { assert } from "utils";

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

  // these dates can be outside of the range of the current month. But it doesn't matter,
  // they are just used as boundaries.
  selectedFrom: DayDate | null;
  selectedTo: DayDate | null;

  entries: Entry[];

  onDayClick: (date: DayDate, extend: boolean) => void;
}

class Month extends React.Component<Props> {
  *entriesOf(target: Date): Generator<Entry, void, void> {
    // FIXME: please be a bit smarter
    for (let entry of this.props.entries) {
      if (
        entry.date.getUTCFullYear() === target.getUTCFullYear() &&
        entry.date.getUTCMonth() === target.getUTCMonth() &&
        entry.date.getUTCDate() === target.getUTCDate()
      ) {
        yield entry;
      }
    }
  }

  // instead of having one big messy condition, have a chain of
  // small meaningful check
  isSelected(target: DayDate): boolean {
    if (this.props.selectedTo === null) {
      assert(
        this.props.selectedFrom === null,
        `selectedFrom should be null, got ${this.props.selectedFrom}`,
      );
      return false;
    }

    if (this.props.selectedFrom === null) {
      return target.equals(this.props.selectedTo);
    }

    let end: DayDate = this.props.selectedTo;
    let start: DayDate = this.props.selectedFrom;
    if (start.obj() > end.obj()) {
      [start, end] = [end, start];
    }

    return (
      target.obj() <= end.obj() &&
      (this.props.selectedFrom === null || target.obj() >= start.obj())
    );
  }

  render() {
    // FIXME: use a constant time, this could cause hard to reproduce bugs
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
                    month={copy.getMonth()}
                    year={copy.getFullYear()}
                    entries={this.entriesOf(copy)}
                    onClick={this.props.onDayClick}
                    dim={copy.getMonth() !== this.props.month}
                    selected={this.isSelected(DayDate.from(copy))}
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
