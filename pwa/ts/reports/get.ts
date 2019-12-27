import { EM, qs, Alerts, State } from "../utils.js";

const HOST = "report get";

export default class {
  section: HTMLElement;
  state: HTMLElement;
  content: HTMLElement;
  filename: HTMLAnchorElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.state = qs(this.section, ".state");
    this.content = qs(this.section, ".report-content");
    this.filename = qs(this.section, ".report-filename") as HTMLAnchorElement;
  }

  async setup() {
    if (State.admin !== true) {
      this.state.innerHTML = "You are not authorized to view this page";
    }
    const params = new URLSearchParams(location.search);
    const filename = params.get("filename");
    if (filename === null) {
      this.filename.textContent = "not specified";
      this.filename.href = "#";
      this.state.innerHTML = `Missing filename parameter. Go back to the
      <a href="/reports/list">report list<a>`;
      return;
    }
    this.filename.textContent = filename;
    this.filename.href = "";

    const url = new URL("/api/reports/get", location.href);
    console.log(url);
    url.search = params.toString();

    const resp = await fetch(url.toString());
    const obj = await resp.json();

    if (obj.kind === "not found" || obj.kind == "unauthorized") {
      if (obj.msg === undefined) {
        Alerts.add({ ...Alerts.invalidResponse, host: HOST });
        console.error("response", obj);
        throw new Error("obj.msg is undefined");
      }
      this.state.innerHTML = obj.msg;
      return;
    } else if (obj.kind !== "success" || obj.report === null) {
      Alerts.add({ ...Alerts.invalidResponse, host: HOST });
      console.error("response", obj);
      throw new Error("obj.kind is unexpected or obj.report is missing");
    }

    this.content.innerHTML = JSON.stringify(obj.report, null, 4);
    this.state.innerHTML = "";
  }
  teardown() {}
}
