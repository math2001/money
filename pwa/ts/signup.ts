import { EM, qs, State } from "./utils.js";
// this is really similar to the login form... how could we reuse the code?

export default class SignUp {
  section: HTMLElement;
  form: HTMLFormElement;
  formstatus: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;

    this.form = qs(this.section, "form.signup-form") as HTMLFormElement;
    this.formstatus = qs(this.section, ".form-status");

    this.form.addEventListener("submit", this.submitForm.bind(this));
  }

  setup() {}

  async submitForm(e: Event) {
    e.preventDefault();
    this.formstatus.innerHTML = "Sending request...";

    const resp = await fetch(this.form.action, {
      method: "post",
      body: new FormData(this.form)
    });

    this.formstatus.innerHTML = "Processing response...";
    const obj = await resp.json();

    if (obj["kind"] !== "success") {
      // FIXME: send minimal error report automatically, and maybe show the
      // user. Don't wanna constantly interupt the users flow
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
    EM.emit(EM.browseto, obj["goto"]);
  }

  teardown() {}
}
