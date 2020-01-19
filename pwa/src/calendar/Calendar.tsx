import React from "react";
import Month from "./Month";
import Header, { MoveType } from "./Header";
import { months, days } from "./data";
import "./Calendar.css";

interface Props {}

interface State {
  year: number;
  month: number;
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

    this.move = this.move.bind(this);
  }

  move(type: MoveType, amount: number) {
    this.setState(state => {
      const newDate = new Date(
        state.year + (type === MoveType.Year ? amount : 0),
        state.month + (type === MoveType.Month ? amount : 0),
        state.dayOfMonth,
      );
      return {
        year: newDate.getFullYear(),
        month: newDate.getMonth(),
        dayOfMonth: newDate.getDate(),
      };
    });
  }

  render() {
    if (this.state.month >= 12) {
      console.error({ month: this.state.month });
      throw new Error("expect state.month < 12");
    }
    return (
      <section className="calendar">
        <Header
          year={this.state.year}
          month={this.state.month}
          move={this.move}
        />

        <Month
          year={this.state.year}
          month={this.state.month}
          dayOfMonth={this.state.dayOfMonth}
        />
      </section>
    );
  }
}

export default Calendar;
