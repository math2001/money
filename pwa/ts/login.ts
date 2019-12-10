export default class Login {
  section: HTMLElement;
  form: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;
    const form = this.section.querySelector(
      "form.login-form"
    ) as HTMLElement | null;
    if (form === null) {
      throw new Error("no form element in login page");
    }
    this.form = form;
  }

  setup() {
    this.form.addEventListener("submit", this._submit);
    this.section.classList.add("active");
  }

  _submit(e: Event) {
    alert("submit form");
    e.preventDefault();
  }

  teardown() {
    this.section.classList.remove("active");
    this.form.removeEventListener("submit", this._submit);

    // FIXME: check that this unfocus any field in the current form
    this.form.blur();
  }
}
