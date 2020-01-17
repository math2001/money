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
    const weekdays = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]
    const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]

    const today = new Date()
    this.monthname.textContent = months[today.getMonth()]

    let firstWeekdayIndex = today.getDay() - (today.getDate() % 7) + 1
    console.log(firstWeekdayIndex, "should be wed", 3)

    const row = document.createElement('tr')
    for (let weekday of weekdays) {
      const th = document.createElement('th')
      th.textContent = weekday
      row.appendChild(th)
    }
    this.month.appendChild(row)

    for (let i = 0; i < 5; i++) {
      const row = document.createElement("tr")
      for (let j = 0; j < 7; j++) {
        const day = document.createElement('td')
        day.innerHTML = '' + (i * 7 + j - firstWeekdayIndex + 1)
        row.appendChild(day)
      }
      this.month.appendChild(row)
    }
  }

  getDayOfMonth(weekIndex, dayIndex) {
    const today = new Date()
    today
  }

  setup() {
    this.generateHTML()
  }

  teardown() {}

}