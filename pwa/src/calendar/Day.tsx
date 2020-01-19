import React from "react";

interface Props {
  // note that these are all 0 based
  weekday: number;
  week: number;
  month: number;
  year: number;
}

class Day extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  // gives the day month number (eg. 1st, 22nd) based on the week number,
  // the day of the week number (0 = sunday, 1 = monday), and the current
  // month. It returns null if that day doesn't exist during that month.
  // note that the weeks *always* start on Sundays (0 is first sunday, 7
  // is second sunday, 9 is second tuesday, etc)
  getMonthNumber() {
    const reference = new Date();
    reference.setFullYear(this.props.year);
    reference.setMonth(this.props.month);
    reference.setDate(1);

    const firstDayOfWeek = reference.getDay();
    const dayMonthNumber =
      this.props.week * 7 + this.props.weekday - firstDayOfWeek + 1;

    // check that the day exists in that month (ie. remove -1, and 32s)
    reference.setDate(dayMonthNumber);
    if (reference.getMonth() !== this.props.month) {
      return null;
    }

    return dayMonthNumber;
  }

  render() {
    return <td className="day">{this.getMonthNumber()}</td>;
  }
}

export default Day;
