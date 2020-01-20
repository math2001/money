import React, { ReactElement } from "react";
import { DayDate } from "./data";

interface Props {
  from: DayDate | null;
  to: DayDate | null;
}

class Details extends React.Component<Props> {
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

    return (
      <article className="calendar-details">
        <p>From: {from}</p>
        <p>To: {this.props.to.render()}</p>
        <p>Balance: $to compute (+total profit, -total expenses)</p>
      </article>
    );
  }
}

export default Details;
