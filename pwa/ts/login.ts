export default class Login {
  section: HTMLElement;
  form: HTMLFormElement;
  formstatus: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;

    const form = this.section.querySelector(
      "form.login-form"
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

  teardown() {
    this.form.removeEventListener("submit", this.submitForm.bind(this));
  }
}
