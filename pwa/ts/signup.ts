// this is really similar to the login form... how could we reuse the code?

export default class SignUp {
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

  setup() {
    this.form.addEventListener("submit", this.submitForm.bind(this));
  }

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
    this.section.classList.remove("active");
    this.form.removeEventListener("submit", this.submitForm.bind(this));

    // FIXME: check that this unfocus any field in the current form
    this.form.blur();
  }
}
