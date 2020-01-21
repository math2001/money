import React, { ChangeEvent } from "react";

interface Props {
  year: number;
  month: number;
  dayOfMonth: number;

  onChange: (date: Date) => void;
}

enum Field {
  year = "year",
  month = "month",
  dayOfMonth = "dayOfMonth",
}

interface State {
  year: number;
  month: number;
  dayOfMonth: number;

  invalid: Field | null;
}

class InputDate extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {
      year: props.year,
      month: props.month,
      dayOfMonth: props.dayOfMonth,
      invalid: null,
    };
    this.onChange = this.onChange.bind(this);
  }

  onChange(event: ChangeEvent) {
    if (!(event.target instanceof HTMLInputElement)) {
      console.log(event.target);
      throw new Error("expected html input element");
    }

    const attr = event.target.attributes.getNamedItem("name");
    if (attr === null) {
      console.error(event.target);
      throw new Error("change event with name=null");
    }
    const name = attr.value;
    if (name !== "year" && name !== "month" && name !== "dayOfMonth") {
      console.error(event.target);
      throw new Error("unexpected name");
    }

    console.log("change", event.target.value);
    const value = parseInt(event.target.value, 10);
    if (isNaN(value)) {
      // @ts-ignore
      this.setState({
        invalid: Field[name],
        [name]: event.target.value, // let the user type
      });
      return;
    }

    this.setState(state => {
      const newState = Object.assign({}, state, {
        [name]: value,
        invalid: null,
      });
      this.props.onChange(
        new Date(newState.year, newState.month, newState.dayOfMonth),
      );
      return newState;
    });
  }

  render() {
    return (
      <div className="inputdate">
        <input
          type="number"
          name="year"
          onChange={this.onChange}
          value={this.state.year}
          data-lpignore="true"
        />
        <input
          type="number"
          name="month"
          onChange={this.onChange}
          value={this.state.month}
          data-lpignore="true"
        />
        <input
          type="number"
          name="dayOfMonth"
          onChange={this.onChange}
          value={this.state.dayOfMonth}
          data-lpignore="true"
        />
      </div>
    );
  }
}

export default InputDate;
