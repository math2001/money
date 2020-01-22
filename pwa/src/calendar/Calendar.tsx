import React from "react";
import Month from "./Month";
import Details from "./Details";
import Header, { MoveType } from "./Header";
import AddEntry from "./AddEntry";
import { DayDate, Entry } from "./data";
import { assert } from "utils";
import { TabSet, Tab } from "mp";
import "./Calendar.css";

enum tab {
  NewEntry = "new entry",
  Details = "details",
}

interface Props {}

interface State {
  year: number;
  month: number;

  selectedFrom: DayDate | null;
  selectedTo: DayDate;

  entries: Entry[] | null;
  activeTab: string;
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
      entries: null,
      activeTab: tab.Details,
    };

    this.move = this.move.bind(this);
    this.onDayClick = this.onDayClick.bind(this);
    this.onDateChange = this.onDateChange.bind(this);
    this.onNewEntrySubmit = this.onNewEntrySubmit.bind(this);
  }

  async componentDidMount() {
    const response = await fetch("/api/load");
    if (response.status !== 200) {
      console.error(response);
      throw new Error("invalid response code");
    }
    const text = await response.text();
    let body;
    try {
      body = JSON.parse(text);
    } catch (e) {
      console.error(body);
      throw e;
    }
    assert(
      body.kind === "success",
      `200 should mean success, got ${body.kind}`,
    );
    this.setState({
      entries: body.entries,
    });
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

  // FIXME: return the error so that the AddEntry component can display it
  async onNewEntrySubmit(entry: Entry) {
    if (this.state.entries === null) {
      throw new Error("state entries are still loading (null)");
    }

    // overwrite the day
    entry.date.setFullYear(this.state.selectedTo.year);
    entry.date.setMonth(this.state.selectedTo.month);
    entry.date.setDate(this.state.selectedTo.dayOfMonth);

    assert(
      entry.id < 0,
      "entry id should be < 0, because Calendar will overwrite it",
    );

    const response = await fetch("/api/add", {
      method: "POST",
      body: JSON.stringify({
        ...entry,
        // only talk in terms of seconds
        date: Math.round(entry.date.getTime() / 1000),
      }),
    });

    if (response.status !== 200) {
      console.error(response);
      throw new Error("expected response 200");
    }

    let responseEntry;
    try {
      responseEntry = await response.json();
    } catch (e) {
      console.error(response.text());
      throw e;
    }
    // the server only talks in terms of seconds, but javascripts understands
    // milliseconds
    responseEntry.date = new Date(responseEntry.date * 1000);

    assert(
      entry.id >= 0,
      "entry id should have been set to different id than -1",
    );

    this.setState(state => {
      if (state.entries === null) {
        throw new Error("state entries are still loading (null)");
      }
      return {
        entries: [...state.entries, entry],
        activeTab: tab.Details,
      };
    });
  }

  render() {
    if (this.state.month >= 12) {
      console.error({ month: this.state.month });
      throw new Error("expect state.month < 12");
    }
    if (this.state.entries === null) {
      return <p>Loading...</p>;
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

        <TabSet
          active={this.state.activeTab}
          onChange={(tabname: string) => this.setState({ activeTab: tabname })}
        >
          <Tab id={tab.Details} title="Details">
            <Details
              from={this.state.selectedFrom}
              to={this.state.selectedTo}
              entries={this.state.entries}
            />
          </Tab>
          <Tab id={tab.NewEntry} title="Add New Entry">
            <AddEntry
              year={this.state.selectedTo.year}
              month={this.state.selectedTo.month}
              dayOfMonth={this.state.selectedTo.dayOfMonth}
              onDateChange={this.onDateChange}
              onSubmit={this.onNewEntrySubmit}
            />
          </Tab>
        </TabSet>
      </section>
    );
  }
}

export default Calendar;
