import { EM, qs } from "./utils.js";
// this is really similar to the login form... how could we reuse the code?

export default class SignUp {
  section: HTMLElement;
  form: HTMLFormElement;
  formstatus: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;

    this.form = qs(this.section, "form.login-form") as HTMLFormElement;
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
      .then((resp: Response) => resp.text())
      .then((text: string) => {
        try {
          return JSON.parse(text);
        } catch (e) {
          console.info(text);
          throw e;
        }
      })
      .then(this.postsignup.bind(this));
  }

  postsignup(obj: any) {
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

    EM.emit(EM.loggedin, obj["email"]);
    EM.emit(EM.browseto, obj["goto"]);
  }

  teardown() {}
}
