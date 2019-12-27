import { EM, qs, State, Alerts } from "../utils.js";

const HOST = "report list";

export default class {
  section: HTMLElement;
  state: HTMLElement;
  reports: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.state = qs(this.section, ".state");
    this.reports = qs(this.section, ".reports");
  }

  async setup() {
    if (State.admin !== true) {
      this.state.innerHTML = "You are not authorized to view this page";
    }

    let resp;
    try {
      resp = await fetch("/api/reports/list");
    } catch (e) {
      Alerts.add({ ...Alerts.networkError, host: HOST });
      return;
    }

    const obj = await resp.json();

    try {
      this.handleResponse(obj);
    } catch (e) {
      Alerts.add({ ...Alerts.invalidResponse, host: HOST });
      console.error(obj);
      throw e;
    }
  }
  async handleResponse(obj: any) {
    if (obj.kind === undefined) {
      throw new Error("no kind field provided");
    }

    if (obj.kind === "unauthorized") {
      if (obj.msg === undefined) {
        throw new Error("msg field expected");
      }
      this.state.innerHTML = obj.msg;
    } else if (obj.kind === "success") {
      if (!Array.isArray(obj.reports)) {
        throw new Error("reports field expected");
      }
      let html = "";
      for (let reportname of obj.reports) {
        html += `<li><a href="/reports/get?filename=${reportname}">${reportname}</a></li>`;
      }
      this.reports.innerHTML = html;
    }
  }
  teardown() {}
}
