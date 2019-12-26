import { State, EM, qs, Alerts } from "./utils.js";
import Home from "./home.js";
import Login from "./login.js";
import SignUp from "./signup.js";
import Err404 from "./err404.js";
import Logout from "./logout.js";
import payments from "./payments/index.js";

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

  constructor() {
    this.current = null;

    Alerts.init(qs(document, "#alerts"));

    this.main = qs(document, "main");

    this.main.addEventListener("click", this.proxyLinks.bind(this));

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
      console.info("browsing to", e.state.url);
      this.browseto(e.state.url);
    });

    EM.on(EM.loggedin, () => {
      const useremail = State.useremail; // load of the localstorage once
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

    this.payments = {
      addManual: new payments.addManual(this.getSection("payments-add-manual")),
      list: new payments.list(this.getSection("payments-list")),
      camera: new payments.camera(this.getSection("camera")),
    };

    if (State.useremail !== null) {
      EM.emit(EM.loggedin);
    }

    history.pushState({ url: location.pathname }, "", location.pathname);
  }

  getSection(name: string): HTMLElement {
    return qs(document, "#" + name);
  }

  proxyLinks(e: MouseEvent) {
    if (e.target !== null && (e.target as Node).nodeName === "A") {
      const target = e.target as HTMLAnchorElement;
      if (this.router((target as HTMLHyperlinkElementUtils).pathname)) {
        e.preventDefault();
        e.stopImmediatePropagation();
        e.stopPropagation();
        EM.emit(EM.browseto, target.pathname);
      }
      // otherwise, we just let the user browse to that URL like any old a tag
      // would do
    }
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

  router(pathname: string): Page | null {
    // FIXME: clean up pathname

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
    } else {
      return null;
    }
  }

  browseto(pathname: string) {
    this.changeto(this.router(pathname) || this.err404);
  }
}

window.addEventListener("load", () => {
  const app = new App();
  app.browseto(location.pathname);
});
