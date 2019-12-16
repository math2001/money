import { qs, EM, State } from "./utils.js";

export default class Login {
  section: HTMLElement;
  form: HTMLFormElement;
  formstatus: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;

    this.form = qs(this.section, "form.login-form") as HTMLFormElement;
    this.form.addEventListener("submit", this.submitForm.bind(this));
    this.formstatus = qs(this.section, ".form-status");
  }

  setup() {}

  submitForm(e: Event) {
    e.preventDefault();
    this.formstatus.innerHTML = "Sending request...";

    fetch(this.form.action, {
      method: "post",
      body: new FormData(this.form)
    })
      .then((resp: Response) => resp.json())
      .then(this.postlogin.bind(this));
  }

  postlogin(obj: any) {
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
