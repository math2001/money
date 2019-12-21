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
    this.task = this.handleLogout().catch((err: Error) => {
      console.error(err);
      this.logoutState.innerHTML = "Error occured: " + err + " (check console)";
    });
  }

  async handleLogout() {
    // we need to promise resolve. It's like sleep(0). It allows othe
    // coroutines to take control. It's needed because otherwise, it is just
    // sync code, and we return if State.useremail === null, hence the caller
    // won't have a returned value (not even a promise), and it doesn't work
    // with the teardown function

    await Promise.resolve();

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
      email: State.useremail
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

    // remove user info as soon as possible. If the user is trying to manually
    // logout, then it's probably because he's already getting errors
    State.useremail = null;

    // FIXME: better error communication
    if (obj.kind === "error" && obj.id === "no user") {
      console.info("no user is currently logged in");
    } else if (obj.kind !== "success") {
      console.error(obj);
      throw new Error("expected kind 'success'");
    }

    if (obj.goto === undefined) {
      console.error(obj);
      throw new Error("expected 'goto' key");
    }

    EM.emit(EM.loggedout);
    EM.emit(EM.browseto, obj.goto);
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
