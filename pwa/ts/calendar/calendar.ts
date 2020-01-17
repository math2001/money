import { EM, qs } from "../utils.js";

export default class Calendar {
  section: HTMLElement;
  month: HTMLTableElement;
  monthname: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.month = qs(this.section, ".month") as HTMLTableElement;
    this.monthname = qs(this.section, ".month-name");
  }

  generateHTML() {
    const weekdays = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
    const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]

    const today = new Date()
    this.monthname.textContent = months[today.getMonth()]


    for (let i = 0; i < 5; i++) {
      const row = document.createElement("tr")
      for (let j = 0; j < 7; j++) {
        const day = document.createElement('td')
        day.innerHTML = weekdays[j] + ' ' + 
        row.appendChild(day)
      }
    }
  }

  setup() {
    this.generateHTML()
  }

  teardown() {}

}