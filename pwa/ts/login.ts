import { qs } from "./utils.js";

export default class Login {
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
      .then((resp: Response) => resp.json())
      .then(this.postlogin.bind(this));
  }

  postlogin(resp: Response) {
    console.log(resp);
  }

  teardown() {}
}
