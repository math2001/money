import { qs, EM, State, Alerts } from "./utils.js";

const HOST = "login";

export default class Login {
  section: HTMLElement;
  form: HTMLFormElement;

  constructor(section: HTMLElement) {
    this.section = section;

    this.form = qs(this.section, "form.login-form") as HTMLFormElement;
    this.form.addEventListener("submit", this.submitForm.bind(this));
  }

  setup() {}

  async submitForm(e: Event) {
    e.preventDefault();
    // make the form data before disabling every input, because otherwise
    // they aren't added to the object
    const formData = new FormData(this.form);

    for (let input of this.form.querySelectorAll("input")) {
      input.disabled = true;
    }
    (qs(this.form, "input[type='submit']") as HTMLInputElement).value =
      "Login in";
    Alerts.removeAll(HOST);

    let resp;
    try {
      resp = await fetch(this.form.action, {
        method: this.form.method,
        body: formData,
      });
    } catch (e) {
      console.error(e);
      Alerts.add({
        html:
          "Network error. Make sure you are connected to the internet " +
          "(check console for more details)",
        kind: Alerts.ERROR,
        host: HOST,
      });
      return;
    } finally {
      for (let input of this.form.querySelectorAll("input")) {
        input.disabled = false;
      }
    }
    const obj = await resp.json();
    this.postlogin(obj);
  }

  postlogin(obj: any) {
    console.info("login response", obj.kind);
    if (obj.kind === "wrong identifiers") {
      Alerts.add({
        html: "Wrong identifiers. Please try again",
        kind: Alerts.ERROR,
        host: HOST,
      });
      return;
    } else if (obj.kind === "internal error") {
      Alerts.add({
        kind: Alerts.ERROR,
        host: HOST,
        html:
          "Oops... Server encoutered an internal error. Please try again " +
          "or report if it keeps on occuring",
      });
    }

    if (obj.kind !== "success") {
      console.error("response:", obj);
      throw new Error("expected 'success' response");
    }

    if (obj.email === undefined) {
      console.error("response:", obj);
      throw new Error("no email field in 'success' response");
    }

    if (obj.goto === undefined) {
      console.error("response:", obj);
      throw new Error("no email field in 'goto' response");
    }

    State.useremail = obj.email;
    EM.emit(EM.loggedin);
    EM.emit(EM.browseto, obj.goto);
  }

  teardown() {}
}
