import React, { ChangeEvent, FormEvent } from "react";
import { Entry } from "./data";
import { cast } from "utils";

interface Props {
  year: number;
  month: number;
  dayOfMonth: number;

  onClose: () => void;
  onDateChange: (change: { field: string; value: number }) => void;
  onSubmit: (entry: Entry) => void;
}

class AddEntry extends React.Component<Props> {
  constructor(props: Props) {
    super(props);

    this.onChange = this.onChange.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
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

    // abstract away the implementation details
    this.props.onDateChange({
      field: event.target.name,
      value: value,
    });
  }

  onSubmit(event: FormEvent) {
    const form = cast(event.target, HTMLFormElement);
    event.preventDefault();

    const time = cast(form.elements.namedItem("time"), HTMLInputElement).value;
    if (time === "") {
      throw new Error("invalid user input [validation not implemented]");
    }
    const splits = time.split(":");
    const hours = parseInt(splits[0], 10);
    const minutes = parseInt(splits[1], 10);

    this.props.onSubmit({
      id: -1,
      name: cast(form.elements.namedItem("name"), HTMLInputElement).value,
      description: cast(
        form.elements.namedItem("description"),
        HTMLTextAreaElement,
      ).value,
      amount: parseInt(
        cast(form.elements.namedItem("amount"), HTMLInputElement).value,
        10,
      ),
      matched: cast(form.elements.namedItem("matched"), HTMLInputElement)
        .checked,
      date: new Date(0, 0, 0, hours, minutes),
    });

    form.reset();
  }

  render() {
    return (
      <article className="add-entry">
        <h3>
          Add Entry <button onClick={this.props.onClose}>&times;</button>
        </h3>
        <form onSubmit={this.onSubmit}>
          <p>
            <label htmlFor="name">Name: </label> <input type="text" id="name" />
          </p>
          <p>
            <label htmlFor="description">Description</label>:
          </p>
          <textarea id="description" />
          <pre>FIXME: have a button "bound date to the calendar"</pre>
          <p>
            <label>Date: </label>
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
            <input type="time" name="time" defaultValue="11:15" />
          </p>
          <p>
            <label htmlFor="amount">Amount</label>:{" "}
            <input type="number" id="amount" />
          </p>
          <p>
            <label htmlFor="matched">Matched</label>:{" "}
            <input type="checkbox" id="matched" />
          </p>
          <p>
            <input type="submit" value="Add entry" />
          </p>
        </form>
      </article>
    );
  }
}

export default AddEntry;
