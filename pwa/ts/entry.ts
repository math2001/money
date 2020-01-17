import { State, EM, qs, Alerts } from "./utils.js";
import Home from "./home.js";
import Login from "./login.js";
import SignUp from "./signup.js";
import Err404 from "./err404.js";
import Logout from "./logout.js";
import Calendar from "./calendar/calendar.js";
import payments from "./payments/index.js";

import ReportsGet from "./reports/get.js";
import ReportsList from "./reports/list.js";

interface Page {
  setup(): void;
  teardown(): void;
  section: HTMLElement;
}

class App {
  current: Page | null;

  main: HTMLElement;

  home: Page;
  login: Page;
  err404: Page;
  signup: Page;
  logout: Page;
  payments: {
    list: Page;
    addManual: Page;
    camera: Page;
  };
  reports: {
    Get: Page;
    List: Page;
  };

  calendar: Page;

  constructor() {
    this.current = null;

    Alerts.init(qs(document, "#alerts"));

    this.main = qs(document, "main");

    this.main.addEventListener("click", (e: MouseEvent) => {
      if (e.target instanceof HTMLAnchorElement) {
        if (this.router(e.target.href)) {
          e.preventDefault();
          e.stopImmediatePropagation();
          e.stopPropagation();
          EM.emit(EM.browseto, e.target.href);
        }
        // otherwise, we just let the user browse to that URL like any old a tag
        // would do
      }
    });

    EM.on(EM.browseto, (url: string) => {
      if (url === undefined) {
        console.trace();
        throw new Error("browsing to undefined URL");
      }

      console.info("browsing to", url);
      history.pushState({ url }, "", url);
      this.browseto(url);
    });

    window.addEventListener("popstate", (e: PopStateEvent) => {
      if (e.state === null) {
        throw new Error(
          "popstate event.state is null. Did you just reload the page?",
        );
      }
      console.info("browsing to", e.state.url);
      this.browseto(e.state.url);
    });

    EM.on(EM.loggedin, () => {
      if (State.user === null) {
        throw new Error("EM.loggedin event triggered, but State.user is null");
      }
      const useremail = State.user.email; // load from the localstorage once
      for (let node of document.querySelectorAll('[fill-with="useremail"]')) {
        node.textContent = useremail;
      }
    });

    EM.on(EM.loggedout, () => {
      for (let node of document.querySelectorAll('[fill-with="useremail"]')) {
        node.textContent = "[internal error]";
      }
    });

    this.home = new Home(this.getSection("home"));
    this.login = new Login(this.getSection("login"));
    this.err404 = new Err404(this.getSection("err404"));
    this.signup = new SignUp(this.getSection("signup"));
    this.logout = new Logout(this.getSection("logout"));

    this.calendar = new Calendar(this.getSection("calendar"));

    this.payments = {
      addManual: new payments.addManual(this.getSection("payments-add-manual")),
      list: new payments.list(this.getSection("payments-list")),
      camera: new payments.camera(this.getSection("camera")),
    };

    this.reports = {
      Get: new ReportsGet(this.getSection("reports-get")),
      List: new ReportsList(this.getSection("reports-list")),
    };

    if (State.user !== null) {
      EM.emit(EM.loggedin);
    }

    history.replaceState({ url: location.href }, "", location.href);
  }

  getSection(name: string): HTMLElement {
    return qs(document, "#" + name);
  }

  changeto(page: Page) {
    if (this.current !== null) {
      this.current.section.classList.remove("active");
      this.current.teardown();
    }
    this.current = page;
    this.current.setup();
    this.current.section.classList.add("active");
  }

  router(URI: string): Page | null {
    // FIXME: clean up pathname

    const pathname = new URL(URI, location.href).pathname;

    if (pathname === "/") {
      return this.home;
    } else if (pathname === "/login") {
      return this.login;
    } else if (pathname === "/signup") {
      return this.signup;
    } else if (pathname == "/logout") {
      return this.logout;
    } else if (pathname == "/payments/add-manual") {
      return this.payments.addManual;
    } else if (pathname == "/payments/list") {
      return this.payments.list;
    } else if (pathname == "/payments/camera") {
      return this.payments.camera;
    } else if (pathname == "/reports/get") {
      return this.reports.Get;
    } else if (pathname == "/reports/list") {
      return this.reports.List;
    } else if (pathname == "/calendar") {
      return this.calendar;
    } else {
      console.error("unknown page: ", pathname);
      return null;
    }
  }

  browseto(href: string) {
    this.changeto(this.router(href) || this.err404);
  }
}

window.addEventListener("load", () => {
  const app = new App();
  app.browseto(location.pathname);
});
