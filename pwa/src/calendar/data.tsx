import React from "react";

export const months = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

export const days = [
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];

export function assert(condition: boolean, message: string): asserts condition {
  if (condition === false) {
    throw new Error(`Assertion Error: ${message}`);
  }
}

export interface Entry {
  id: number;
  name: string;
  description: string;
  amount: number;
  date: Date;
  matched: boolean;
}

export class DayDate {
  year: number;
  month: number;
  dayOfMonth: number;

  static from(target: Date): DayDate {
    return new DayDate(
      target.getFullYear(),
      target.getMonth(),
      target.getDate(),
    );
  }

  constructor(year: number, month: number, dayOfMonth: number) {
    this.year = year;
    this.month = month;
    this.dayOfMonth = dayOfMonth;
  }

  render() {
    return (
      <span>
        {this.year} {months[this.month]} {this.dayOfMonth}
      </span>
    );
  }

  obj(): Date {
    // always return the dates at the same time (12 o'clock)
    return new Date(this.year, this.month, this.dayOfMonth, 12, 0, 0, 0);
  }

  equals(target: DayDate): boolean {
    return (
      this.year === target.year &&
      this.month === target.month &&
      this.dayOfMonth === target.dayOfMonth
    );
  }
}
