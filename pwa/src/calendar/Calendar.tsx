import React from "react";
import Month from "./Month";
import { months, days } from "./data";
import "./Calendar.css";

interface Props {}

interface State {
  month: number;
  year: number;
  // dayOfMonth is the number of the day in the month (1st, 2nd, ... 31st)
  // just like Date::setDate
  dayOfMonth: number;
}

class Calendar extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const today = new Date();
    this.state = {
      year: today.getFullYear(),
      month: today.getMonth(),
      dayOfMonth: today.getDate(),
    };
  }

  render() {
    if (this.state.month >= 12) {
      console.error({ month: this.state.month });
      throw new Error("expect state.month < 12");
    }
    return (
      <Month
        year={this.state.year}
        month={this.state.month}
        dayOfMonth={this.state.dayOfMonth}
      />
    );
  }
}

export default Calendar;
