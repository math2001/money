import Home from "./home.js";
import Login from "./login.js";
import Err404 from "./err404.js";

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

  constructor() {
    this.current = null;

    const main = document.querySelector("main") as HTMLElement | null;
    if (main === null) {
      throw new Error("no main element");
    }
    this.main = main;

    this.main.addEventListener("click", this.proxyLinks.bind(this));

    this.home = new Home(this.getSection("home"));
    this.login = new Login(this.getSection("login"));
    this.err404 = new Err404(this.getSection("err404"));
  }

  getSection(name: string): HTMLElement {
    const section = document.querySelector("#" + name) as HTMLElement | null;
    if (section === null) {
      throw new Error(`Element (section) #${name} not found`);
    }
    return section;
  }

  proxyLinks(e: MouseEvent) {
    if ((e.target as Node).nodeName === "A") {
      if (this.router((e.target as HTMLHyperlinkElementUtils).pathname)) {
        e.preventDefault();
        e.stopImmediatePropagation();
        e.stopPropagation();
        this.browseto(e.target.pathname);
        history.pushState({}, "", e.target.pathname);
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

    if (pathname === "/login") {
      return this.login;
    } else if (pathname === "/") {
      return this.home;
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
