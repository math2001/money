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

export class State {
  static useremail: string | null = null;
}

export function qs(from: Element | Document, selector: string): HTMLElement {
  const element = from.querySelector(selector) as HTMLElement | null;
  if (element === null) {
    console.error(`${selector} not found in`, from);
    throw new Error(`element ${selector} not found`);
  }
  return element;
}
