export class EM {
  // define them as constant so that we get compile time checking

  static browseto = "browseto";
  static loggedin = "loggedin";
  static loggedout = "loggedout";
  static events: { [key: string]: Function[] } = {};

  static on(eventName: string, cb: Function): void {
    if (EM.events[eventName] === undefined) {
      EM.events[eventName] = [];
    }
    EM.events[eventName].push(cb);
  }
  static emit(eventName: string, ...args: any): void {
    for (let cb of EM.events[eventName]) {
      cb(...args);
    }
  }
}

// state just provides the keys for the local storage (to get compile time
// checking). For eg. localStorage.setItem(State.useremail, JSON.stringify("<useremail>"))
// (items stored in the local storage should always be json encoded)
export class State {
  private static setItem(key: string, value: any) {
    localStorage.setItem(key, JSON.stringify(value));
  }
  private static getItem(key: string): any {
    const item = localStorage.getItem(key);
    if (item === null) {
      return null;
    }
    return JSON.parse(item);
  }

  static get useremail(): string | null {
    return State.getItem("useremail");
  }

  static set useremail(value: string | null) {
    State.setItem("useremail", value);
  }

  static get admin(): boolean | null {
    return State.getItem("admin");
  }

  static set admin(value: boolean | null) {
    State.setItem("admin", value);
  }
}

export function qs(from: Element | Document, selector: string): HTMLElement {
  const element = from.querySelector(selector) as HTMLElement | null;
  if (element === null) {
    console.error(`${selector} not found in`, from);
    throw new Error(`element ${selector} not found`);
  }
  return element;
}

interface AlertParams {
  kind: number;
  html: string;
  host: string;
}

export class Alerts {
  static root: HTMLElement;
  static alerts: HTMLElement[];

  static ERROR = 0;
  static WARNING = 1;
  static SUCCESS = 2;

  static invalidResponse = {
    kind: Alerts.ERROR,
    html: `The response from the server was invalid.
           This is <em>our</em> fault, <strong>not yours</strong>.
           This problem has been reported. More details are
           available in the console.`,
  };

  static badRequest = {
    kind: Alerts.ERROR,
    html: `The request sent was invalid.
           This is <em>our</em> fault, <strong>not yours</strong>.
           This problem has been reported. More details are
           available in the console.`,
  };

  static serverInternalError = {
    kind: Alerts.ERROR,
    html: `An error occured on the server.
           This is <em>our</em> fault, <strong>not yours</strong>.
           This problem has been reported. More details are
           available in the console.`,
  };

  static init(root: HTMLElement) {
    this.alerts = [];
    this.root = root;

    this.root.addEventListener("click", (e: Event) => {
      if (e.target === null || !(e.target instanceof HTMLButtonElement)) {
        return;
      }
      if (!e.target.classList.contains("alert-close-btn")) {
        return;
      }
      const alert = e.target.parentElement!.parentElement!;
      this.root.removeChild(alert);
    });
  }

  private static typeToString(type: number): string {
    switch (type) {
      case this.ERROR:
        return "error";
        break;
      case this.SUCCESS:
        return "success";
        break;
      case this.WARNING:
        return "warning";
        break;
      default:
        throw new Error(`unkonwn alert type "${type}"`);
    }
  }

  // addAlert creates an alert and adds it to the html. It returns its id
  // which can then be used
  static add(params: AlertParams) {
    const alert = document.createElement("article");
    alert.innerHTML = `
    <div class="alert-content">
      ${params.html}
    </div>
    <div class="alert-close">
      <button class="alert-close-btn">&times</button>
    </div>
    `;
    alert.classList.add("alert");
    alert.classList.add("alert-" + this.typeToString(params.kind));
    alert.setAttribute("data-host", params.host);
    this.alerts.push(alert);
    this.root.appendChild(alert);
  }

  static removeAll(host: string) {
    for (let alert of this.root.querySelectorAll(
      `.alert[data-host="${host}"]`,
    )) {
      const parent = alert.parentElement;
      if (parent === null) {
        console.error(alert);
        throw new Error("can't remove alert, no parent");
      }
      parent.removeChild(alert);
    }
  }
}
