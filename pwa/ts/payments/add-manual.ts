import { EM, qs } from "../utils.js";

export default class AddManual {
  section: HTMLElement;
  form: HTMLFormElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.form = qs(this.section, "form") as HTMLFormElement;

    this.form.addEventListener("submit", this.onsubmit.bind(this));
  }

  async onsubmit(e: Event) {
    e.preventDefault();

    const resp = await fetch(this.form.action, {
      method: this.form.method,
      body: new FormData(this.form)
    });
    const obj = await resp.json();
    if (obj.kind === undefined) {
      console.error(obj);
      throw new Error("expected 'kind' field");
    }
    if (obj.kind !== "success") {
      console.error(obj);
      throw new Error("expected 'kind' 'success'");
    }
    if (obj.goto === undefined) {
      console.error(obj);
      throw new Error("expected 'goto' field");
    }
    EM.emit(EM.browseto, obj.goto);

    // FIXME: EM.emit(EM.notificiation, "success", "Payment added!")
  }

  setup() {}

  teardown() {}
}
