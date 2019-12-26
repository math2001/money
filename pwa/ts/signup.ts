import { EM, qs, State, Alerts } from "./utils.js";
// this is really similar to the login form... how could we reuse the code?

const HOST = "signup";

export default class SignUp {
  section: HTMLElement;
  form: HTMLFormElement;

  constructor(section: HTMLElement) {
    this.section = section;

    this.form = qs(this.section, "form.signup-form") as HTMLFormElement;

    this.form.addEventListener("submit", this.submitForm.bind(this));
  }

  setup() {}

  async submitForm(e: Event) {
    e.preventDefault();
    const formData = new FormData(this.form);

    for (let input of this.form.querySelectorAll("input")) {
      input.disabled = true;
    }
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

    if (obj.kind === "bad request") {
      Alerts.add({
        kind: Alerts.ERROR,
        html: `The request we sent was invalid. This is <em>our</em> fault, <strong>not yours</strong>. This problem has been reported. More details are available in the console`,
        host: HOST,
      });
      console.error(obj);
      throw new Error("bad request (it wasn't formatted properly)");
    } else if (obj.kind === "invalid input") {
      if (obj.msg === undefined) {
        console.log("response", obj);
        Alerts.add({
          kind: Alerts.ERROR,
          host: HOST,
          html: `invalid response from the server. This is <em>our</em> fault, <strong>not yours</strong>. This problem is being reported. More details available in the console.`,
        });
        throw new Error("invalid response from the server (signup)");
      }
      Alerts.add({
        kind: Alerts.ERROR,
        html: obj.msg,
        host: HOST,
      });
      return;
    } else if (obj.kind === "password too short") {
      if (obj.msg === undefined) {
        console.error("response", obj);
        Alerts.add({
          kind: Alerts.ERROR,
          host: HOST,
          html: `invalid response from the server. This is <em>our</em> fault, <strong>not yours</strong>. This problem is being reported. More details available in the console.`,
        });
        throw new Error(
          "invalid response from the server (signup), missing keys",
        );
      }
      Alerts.add({
        kind: Alerts.ERROR,
        html: obj.msg,
        host: HOST,
      });
      return;
    } else if (obj.kind === "password dismatch") {
      if (obj.msg === undefined) {
        console.error("response", obj);
        Alerts.add({
          kind: Alerts.ERROR,
          host: HOST,
          html: `invalid response from the server. This is <em>our</em> fault, <strong>not yours</strong>. This problem is being reported. More details available in the console.`,
        });
        throw new Error(
          "invalid response from the server (signup), missing keys",
        );
      }
      for (let input of this.form.querySelectorAll<HTMLInputElement>(
        "input[type='password']",
      )) {
        input.value = "";
      }
      Alerts.add({
        kind: Alerts.ERROR,
        host: HOST,
        html: obj.msg,
      });
      return;
    } else if (obj.kind !== "success") {
      // FIXME: send minimal error report automatically, and maybe show the
      // user. Don't wanna constantly interupt the users flow
      Alerts.add({
        kind: Alerts.ERROR,
        host: HOST,
        html: `invalid response from the server. This is <em>our</em> fault, <strong>not yours</strong>. This problem is being reported. More details available in the console.`,
      });
      console.error("response:", obj);
      throw new Error("invalid response kind");
    }

    if (obj.email === undefined || obj.goto === undefined) {
      Alerts.add({
        kind: Alerts.ERROR,
        host: HOST,
        html: `invalid response from the server. This is <em>our</em> fault, <strong>not yours</strong>. This problem is being reported. More details available in the console.`,
      });
      console.error("response:", obj);
      throw new Error("invalid server response");
    }

    State.useremail = obj.email;
    EM.emit(EM.loggedin);
    EM.emit(EM.browseto, obj.goto);
  }

  teardown() {}
}
