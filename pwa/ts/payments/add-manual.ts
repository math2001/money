import { EM, qs } from "../utils.js";

export default class AddManual {
  section: HTMLElement;
  form: HTMLFormElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.form = qs(this.section, "form") as HTMLFormElement;

    this.form.addEventListener("submit", this.onsubmit.bind(this));

    this.form.addEventListener("click", (e: MouseEvent) => {
      if (!(e.target instanceof HTMLElement)) {
        return;
      }
      if (e.target.classList.contains("remove-field")) {
        const row = e.target.parentElement;
        if (row === null) {
          throw new Error("no parent element for button");
        }
        if (row.nodeName !== "P") {
          console.error("expected to remove p element, got", row);
        }
        if (row.parentElement === null) {
          throw new Error("element to remove (p) has no parent");
        }

        row.parentElement.removeChild(row);
      }
    });

    qs(this.section, "#add-field").addEventListener("click", (e: Event) => {
      e.preventDefault();
      this.addfield("text", false);
    });

    this.addfield("text", true, "name");
    const value = qs(
      this.addfield("number", true, "amount"),
      "input[name^='value']"
    ) as HTMLInputElement;
    value.step = "any";
    this.addfield("date", true, "date");
  }

  addfield(type: string, required: boolean, fieldname?: string): HTMLElement {
    const n = this.form.querySelectorAll("input[name^='field']").length;

    const p = document.createElement("p");

    const field = document.createElement("input");
    field.name = "field" + n;
    field.placeholder = "field name";
    if (fieldname !== undefined) {
      field.value = fieldname;
      field.readOnly = true;
    }

    const value = document.createElement("input");
    value.type = type;
    value.name = "value" + n;
    value.placeholder = "value";

    if (!required) {
      const button = document.createElement("button");
      button.type = "button";
      button.textContent = "Remove";
      button.classList.add("remove-field");
      p.appendChild(button);
    }

    p.appendChild(field);
    p.appendChild(value);

    this.form.appendChild(p);
    return p;
  }

  async onsubmit(e: Event) {
    e.preventDefault();

    const payment: { [key: string]: any } = {};

    for (let input of this.section.querySelectorAll<HTMLInputElement>(
      'input[name^="field"]'
    )) {
      const name = input.getAttribute("name") as string;
      if (name in payment) {
        throw new Error("duplicate keys in form (internal error)");
      }

      if (input.value === "") {
        // FIXME: better error reporting
        alert("empty fields not allowed (remove empty inputs)");
        return;
      }

      if (payment[input.value] !== undefined) {
        // FIXME: better error reporting
        alert("duplicate field name: " + input.value);
        return;
      }
      const corresponding = qs(
        this.section,
        `input[name="${name.replace("field", "value")}"]`
      ) as HTMLInputElement;
      if (corresponding.value === "") {
        // FIXME: better error reporting
        alert(
          "empty values not allowed (fill required inputs, and/or remove empty inputs)"
        );
        return;
      }
      payment[input.value] = corresponding.value;
    }

    payment["amount"] = parseFloat(payment["amount"]);

    const formdata = new FormData();
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
