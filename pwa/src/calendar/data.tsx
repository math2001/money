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

export interface Entry {
  id: number;
  name: string;
  description: string;
  amount: number;
  date: Date;
  matched: boolean;
}

export interface ServerEntry {
  id: number;
  name: string;
  description: string;
  amount: number;
  date: number; // timestamp
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
    // always return the dates at the same time
    return new Date(this.year, this.month, this.dayOfMonth, 0, 0, 0, 0);
  }

  equals(target: DayDate): boolean {
    return (
      this.year === target.year &&
      this.month === target.month &&
      this.dayOfMonth === target.dayOfMonth
    );
  }

  between(start: DayDate | Date, end: DayDate | Date): boolean {
    let start_: Date;
    let end_: Date;
    if (start instanceof DayDate) {
      start_ = start.obj();
    } else {
      start_ = start;
    }

    if (end instanceof DayDate) {
      end_ = end.obj();
    } else {
      end_ = end;
    }

    if (start_ > end_) {
      [start_, end_] = [end_, start_];
    }

    const obj = this.obj();
    return obj > start_ && obj < end;
  }

  static between(
    start: DayDate | Date,
    end: DayDate | Date,
    target: DayDate | Date,
  ): boolean {
    let start_: Date;
    let end_: Date;
    let target_: Date;

    if (start instanceof DayDate) {
      start_ = start.obj();
    } else {
      start_ = start;
    }

    if (end instanceof DayDate) {
      end_ = end.obj();
    } else {
      end_ = end;
    }

    if (target instanceof DayDate) {
      target_ = target.obj();
    } else {
      target_ = target;
    }

    if (start_ > end_) {
      [start_, end_] = [end_, start_];
    }

    return target_ > start_ && target_ < end;
  }

  toString() {
    return `${this.year}-${this.month}-${this.dayOfMonth}`;
  }

  copy(): DayDate {
    return new DayDate(this.year, this.month, this.dayOfMonth);
  }
}
