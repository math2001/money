import { EM, qs, Alerts } from "../utils.js";

const HOST = "payments/list";

export default class List {
  section: HTMLElement;
  table: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.table = qs(this.section, "table");
  }

  setup() {
    this.table.innerHTML = "Loading payments from server...";
    this.load();
  }

  async load() {
    const resp = await fetch("/api/payments/list");
    const obj = await resp.json();
    if (obj.kind !== "success") {
      this.table.innerHTML = "Error occured. More details in the console.";
      Alerts.add({ ...Alerts.invalidResponse, host: HOST });
      console.error("response", obj);
      throw new Error("expected kind 'success'");
    }

    const payments = obj.payments;
    if (payments === null) {
      this.table.textContent = "No payments yet";
      return;
    }
    if (!Array.isArray(payments)) {
      console.error(payments);
      throw new Error("expected array of payments");
    }

    this.table.innerHTML = "";
    const head = document.createElement("tr");

    const fields = new Set<string>();
    for (let p of payments) {
      for (let field of Object.keys(p)) {
        fields.add(field);
      }
    }

    for (let field of fields) {
      const cell = document.createElement("th");
      cell.textContent = field;
      head.appendChild(cell);
    }

    this.table.appendChild(head);

    for (let p of payments) {
      const row = document.createElement("tr");

      for (let field of fields) {
        const cell = document.createElement("td");
        const value = p[field];

        if (typeof value === "number") {
          cell.textContent = String(value);
        } else if (typeof value === "string") {
          cell.textContent = value;
        } else if (typeof value === "undefined") {
          cell.textContent = "";
        } else {
          throw new Error(`unsupported field type ${field} (value: ${value})`);
        }

        row.appendChild(cell);
      }
      this.table.appendChild(row);
    }
  }
  teardown() {}
}
