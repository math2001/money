import React from "react";
import Month from "./Month";
import Details from "./Details";
import Header, { MoveType } from "./Header";
import AddEntry from "./AddEntry";
import { DayDate, Entry } from "./data";
import "./Calendar.css";

interface Props {}

interface State {
  year: number;
  month: number;

  selectedFrom: DayDate | null;
  selectedTo: DayDate;

  entries: Entry[];
  addingNewEntry: boolean;
}

class Calendar extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const today = new Date();
    this.state = {
      year: today.getFullYear(),
      month: today.getMonth(),
      selectedFrom: null, // null means from the very beginning
      selectedTo: DayDate.from(today),
      entries: entries,
      addingNewEntry: false,
    };

    this.move = this.move.bind(this);
    this.onDayClick = this.onDayClick.bind(this);
    this.onDateChange = this.onDateChange.bind(this);
  }

  move(type: MoveType, amount: number) {
    this.setState(state => {
      const newDate = new Date(
        state.year + (type === MoveType.Year ? amount : 0),
        state.month + (type === MoveType.Month ? amount : 0),
      );
      // if the displayed months is the current month and nothing
      // is currently selected, select the current date otherwise,
      // select nothing
      let selectedTo: DayDate | null = state.selectedTo;
      if (selectedTo === null) {
        const today = new Date();
        if (newDate.getMonth() === today.getMonth()) {
          selectedTo = DayDate.from(newDate);
        }
        // just to make sure
        if (state.selectedFrom !== null) {
          console.error({
            selectedTo: state.selectedTo,
            selectedFrom: state.selectedFrom,
          });
          throw new Error("selectedTo === null but selectedFrom !== null");
        }
      }

      return {
        year: newDate.getFullYear(),
        month: newDate.getMonth(),
        selectedTo: selectedTo,
      };
    });
  }

  onDayClick(date: DayDate, extend: boolean) {
    // if only one day is highlighted, then selectedTo should be set
    // and selectedFrom is null.
    this.setState(state => {
      if (state.selectedTo === null || extend === false) {
        return {
          selectedTo: date,
          selectedFrom: null,
        };
      }
      return {
        selectedFrom: date,
        selectedTo: state.selectedTo,
      };
    });
  }

  onDateChange({ field, value }: { field: string; value: number }) {
    const selectedTo = this.state.selectedTo.copy();
    if (field !== "year" && field !== "month" && field !== "dayOfMonth") {
      console.error({ field });
      throw new Error("invalid date change field");
    }
    selectedTo[field] = value;
    this.setState({
      selectedTo: selectedTo,
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
          onDayClick={this.onDayClick}
          selectedFrom={this.state.selectedFrom}
          selectedTo={this.state.selectedTo}
          entries={this.state.entries}
        />

        <Details
          from={this.state.selectedFrom}
          to={this.state.selectedTo}
          entries={this.state.entries}
        />

        {this.state.addingNewEntry && (
          <AddEntry
            year={this.state.selectedTo.year}
            month={this.state.selectedTo.month}
            dayOfMonth={this.state.selectedTo.dayOfMonth}
            onClose={() => this.setState({ addingNewEntry: false })}
            onDateChange={this.onDateChange}
          />
        )}
        {!this.state.addingNewEntry && (
          <button onClick={() => this.setState({ addingNewEntry: true })}>
            Add Entry
          </button>
        )}
      </section>
    );
  }
}

export default Calendar;

// for debug purposes... Obviously, this will be fetched from the server later on

const entries: Entry[] = [
  {
    id: 0,
    name: "first",
    description: "",
    amount: 10,
    date: new Date(),
    matched: true,
  },
  {
    id: 1,
    name: "second",
    description: "Hello world",
    amount: -10,
    date: new Date(2020, 0, 10),
    matched: true,
  },
  {
    id: 3,
    name: "third",
    description: "Hello world, some long description...",
    amount: 50,
    date: new Date(2019, 11, 15),
    matched: false,
  },
  {
    id: 4,
    name: "fourth",
    description: "Hello world, some long description...",
    amount: 50,
    date: new Date(2019, 11, 15),
    matched: false,
  },
];
