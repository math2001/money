import React, { ChangeEvent } from "react";
import InputDate from "InputDate";
import { DayDate } from "./data";

interface Props {
  year: number;
  month: number;
  dayOfMonth: number;

  onClose: () => void;
  onDateChange: (change: { field: string; value: number }) => void;
}

class AddEntry extends React.Component<Props> {
  constructor(props: Props) {
    super(props);

    this.onChange = this.onChange.bind(this);
  }

  onChange(event: ChangeEvent) {
    if (!(event.target instanceof HTMLInputElement)) {
      console.error(event.target);
      throw new Error("expected HTMLInputElement");
    }
    const value = parseInt(event.target.value, 10);
    if (isNaN(value)) {
      return;
    }

    this.props.onDateChange({
      field: event.target.name,
      value: value,
    });
  }

  render() {
    return (
      <article className="add-entry">
        <h3>
          Add Entry <button onClick={this.props.onClose}>&times;</button>
        </h3>
        <p>
          <input
            type="number"
            name="year"
            value={this.props.year}
            onChange={this.onChange}
          />
          <input
            type="number"
            name="month"
            value={this.props.month}
            onChange={this.onChange}
          />
          <input
            type="number"
            name="dayOfMonth"
            value={this.props.dayOfMonth}
            onChange={this.onChange}
          />
        </p>
      </article>
    );
  }
}

export default AddEntry;
