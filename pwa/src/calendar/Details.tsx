import React, { ReactElement } from "react";
import { DayDate, Entry } from "./data";
import { assert } from "utils";

interface Props {
  from: DayDate | null;
  to: DayDate;

  entries: Entry[];
}

class Details extends React.Component<Props> {
  render() {
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

    const selectedEntries = this.props.entries.filter(entry => {
      return (start === null || entry.date > start) && entry.date < end;
    });

    // compute a nice report
    // it'd be more efficient to do from the filter above, but it's going to allow
    // adding reports field much easier I think. (lol. early optimisation. damn)

    let expenses = 0.0;
    let incomes = 0.0;
    let matchedExpenses = 0.0;
    let matchedIncomes = 0.0;

    let nExpenses = 0;
    let nIncomes = 0;
    let nUnmatchedExpenses = 0;
    let nUnmatchedIncomes = 0;

    for (let entry of selectedEntries) {
      if (entry.amount > 0) {
        incomes += entry.amount;
        if (entry.matched) {
          matchedIncomes += entry.amount;
        } else {
          nUnmatchedIncomes++;
        }
        nIncomes++;
      } else {
        // minus to get a positive expense total
        expenses -= entry.amount;
        if (entry.matched) {
          matchedExpenses -= entry.amount;
        } else {
          nUnmatchedExpenses++;
        }
        nExpenses++;
      }
    }

    assert(
      nIncomes + nExpenses === selectedEntries.length,
      "# incomes + # expenses != # entries",
    );
    assert(nIncomes >= nUnmatchedIncomes, "# unmatched > # all");
    assert(nExpenses >= nUnmatchedExpenses, "# unmatched > # all");

    let from: ReactElement | null = null;
    if (this.props.from !== null) {
      from = this.props.from.render();
    }

    return (
      <article className="calendar-details">
        {from ? (
          <p>
            From {from} to {this.props.to.render()}
          </p>
        ) : (
          <p>{this.props.to.render()}</p>
        )}
        <p>
          -${expenses} from {nExpenses} expenses{" "}
          <small>
            ({nUnmatchedExpenses} weren't matched; ${expenses - matchedExpenses}
            )
          </small>
        </p>
        <p>
          +${incomes} from {nIncomes} incomes{" "}
          <small>
            ({nUnmatchedIncomes} weren't matched; ${incomes - matchedIncomes})
          </small>
        </p>
        <p>
          {" "}
          Balance: ${incomes - expenses}{" "}
          <small>${matchedIncomes - matchedExpenses} matched</small>
        </p>
      </article>
    );
  }
}

export default Details;
