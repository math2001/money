import { EM, State, qs } from "./utils.js";

export default class Logout {
  section: HTMLElement;
  logoutState: HTMLElement;
  task: Promise<any> | null;

  constructor(section: HTMLElement) {
    this.section = section;

    this.logoutState = qs(this.section, ".logout-state");
    this.task = null;
  }

  setup() {
    this.task = this.handleLogout();
  }

  async handleLogout() {
    this.logoutState.innerHTML = "Checking logged state...";
    if (State.useremail == null) {
      // FIXME: better error communication, redirect to error page
      alert("logging out of nothing (you are not logged in)");
      EM.emit(EM.browseto, "/");
      return;
    }

    this.logoutState.innerHTML = "Sending logout request to the servers...";

    // FIXME: handle offline
    const params: { [key: string]: string } = {
      useremail: State.useremail
    };

    const formData = new FormData();
    for (let key in params) {
      formData.append(key, params[key]);
    }

    const resp = await fetch("/api/logout", {
      method: "post",
      body: formData
    });

    const obj = await resp.json();
    // FIXME: better error communication
    if (obj.kind !== "success") {
      console.error(obj);
      throw new Error("expected kind 'success'");
    }

    if (obj.goto === undefined) {
      console.error(obj);
      throw new Error("expected 'goto' key");
    }

    State.useremail = null;
    EM.emit(EM.loggedout);
  }

  teardown() {
    if (this.task === null) {
      // FIXME: better error communication
      throw new Error("tearing down logout page that hasn't been setup");
    }

    this.task.catch((err: Error) => {
      console.error("tearing down logout task");
      console.error(err);
      throw err;
    });
  }
}
