import { EM } from "./utils.js";
// this is really similar to the login form... how could we reuse the code?

export default class SignUp {
  section: HTMLElement;
  form: HTMLFormElement;
  formstatus: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;

    const form = this.section.querySelector(
      "form.signup-form"
    ) as HTMLFormElement | null;
    if (form === null) {
      throw new Error("no form element in login page");
    }
    this.form = form;

    const formstatus = this.form.querySelector(
      ".form-status"
    ) as HTMLElement | null;
    if (formstatus === null) {
      throw new Error("no .form-status element in login page");
    }
    this.formstatus = formstatus;
  }

  setup() {
    this.form.addEventListener("submit", this.submitForm.bind(this));
  }

  submitForm(e: Event) {
    e.preventDefault();
    this.formstatus.innerHTML = "Sending request...";

    console.log(this.form.action);
    fetch(this.form.action, {
      method: "post",
      body: new FormData(this.form),
      redirect: "error"
    })
      .then((r: Response) => r.json())
      .then(this.postsignup.bind(this));
  }

  postsignup(obj: any) {
    if (obj["kind"] != "goto") {
      // FIXME: send minimal error report automatically, and maybe show the
      // user. Don't wanna constantly interupt the users flow
    }

    EM.emit(EM.browseto, obj["goto"]);
  }

  teardown() {
    this.section.classList.remove("active");
    this.form.removeEventListener("submit", this.submitForm.bind(this));

    // FIXME: check that this unfocus any field in the current form
    this.form.blur();
  }
}
