import React from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import "./App.css";

import CalendarPage from "./calendar/CalendarPage";

const App: React.FC = () => {
  return (
    <Router>
      <Switch>
        <Route path="/calendar">
          <CalendarPage />
        </Route>
        <Route path="/tableview">
          <p>
            Nothing to see here, move along... <Link to="/">Home</Link>
          </p>
        </Route>
      </Switch>
      <nav>
        <ul>
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <Link to="/calendar">Calendar</Link>
          </li>
          <li>
            <Link to="/table">[TODO] Table View</Link>
          </li>
        </ul>
      </nav>
    </Router>
  );
};

export default App;
