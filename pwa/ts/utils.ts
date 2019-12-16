export class EM {
  // define them as constant so that we get compile time checking

  static browseto = "browseto";
  static loggedin = "loggedin";
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
}

export function qs(from: Element | Document, selector: string): HTMLElement {
  const element = from.querySelector(selector) as HTMLElement | null;
  if (element === null) {
    console.error(`${selector} not found in`, from);
    throw new Error(`element ${selector} not found`);
  }
  return element;
}
