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

    let formdata = new FormData(this.form);
    const payment: { [key: string]: any } = {};
    for (let key of formdata.keys()) {
      if (key.startsWith("field")) {
        const strkey = formdata.get(key);
        if (typeof strkey !== "string") {
          throw new Error("type of 'field#' input should be string");
        }
        if (strkey in payment) {
          throw new Error("duplicate keys");
        }
        payment[strkey] = formdata.get(key.replace("field", "value"));
      }
    }

    formdata = new FormData();
    formdata.append("payment", JSON.stringify(payment));
    const resp = await fetch(this.form.action, {
      method: this.form.method,
      body: formdata
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
