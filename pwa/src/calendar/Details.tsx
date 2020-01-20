import React, { ReactElement } from "react";
import { DayDate, Entry, assert } from "./data";

interface Props {
  from: DayDate | null;
  to: DayDate | null;

  entries: Entry[];
}

interface Report {
  incomes: number;
  expenses: number;
}

class Details extends React.Component<Props> {
  getReport(): Report {
    assert(
      this.props.to !== null,
      "should have custom message when to (and from) is null (no data to display)",
    );

    let start: Date | null;
    let end: Date;

    if (this.props.from !== null) {
      start = this.props.from.obj();
      end = this.props.to.obj();
      if (start > end) {
        [start, end] = [end, start];
      }
    } else {
      start = null;
      end = this.props.to.obj();
    }

    // move forward one day because the last day is *excluded* from the calculations
    // (like str[1:3] excludes the character at index 3) but it's highlighted as
    // selected.
    end.setDate(end.getDate() + 1);
    console.log("selecting from", start, "to", end);

    const selectedEntries = this.props.entries.filter(entry => {
      return (start === null || entry.date > start) && entry.date < end;
    });

    // compute a nice report
    // it'd be more efficient to do from the filter above, but it's going to allow
    // adding reports field much easier I think. (lol. early optimisation. damn)

    let expenses = 0.0;
    let incomes = 0.0;
    for (let entry of selectedEntries) {
      if (entry.amount > 0) {
        incomes += entry.amount;
      } else {
        // minus to get a positive expense total
        expenses -= entry.amount;
      }
    }

    return { expenses, incomes };
  }

  render() {
    if (this.props.to === null) {
      return (
        <article className="calendar-details calendar-details-hidden">
          Nothing to display
        </article>
      );
    }
    let from: ReactElement | string;
    if (this.props.from !== null) {
      from = this.props.from.render();
    } else {
      from = "beginning";
    }

    const { expenses, incomes } = this.getReport();
    return (
      <article className="calendar-details">
        <p>From: {from}</p>
        <p>To: {this.props.to.render()}</p>
        <p>
          Balance: {incomes - expenses} (+{expenses}, -{incomes})
        </p>
      </article>
    );
  }
}

export default Details;
