export class EM {
  // define them as constant so that we get compile time checking

  static browseto = "browseto";
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
