import { EM, qs } from "../utils.js";


function addDays(date: Date, n: number): Date {
  const next = new Date(date.getTime())
  next.setDate(date.getDate() + n)
  return next
}

function dateEqual(a: Date, b: Date) {
  return a.getUTCDate() == b.getUTCDate() && a.getUTCFullYear() == b.getUTCFullYear() && a.getUTCMonth() == b.getUTCMonth()
}

function dateNormalize(a: Date): Date {
  const copy = new Date(a.getTime())
  copy.setUTCHours(8, 0, 0, 0)
  return copy
}

const today = dateNormalize(new Date())
console.log(addDays(today, 1))

const changes = new Map()

changes.set(addDays(today, 1).getTime(), [
  {name: "first income", description: "", amount: 200},
  {name: "first expense", description: "", amount: -20}
])

changes.set(addDays(today, 2).getTime(), [
  {name: "hello", description: "", amount: -50}
])

for (let key of changes.keys()) {
  console.log(key)
}

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

    const today = dateNormalize(new Date())

    this.monthname.textContent = months[today.getMonth()]

    let firstWeekdayIndex = today.getDay() - (today.getDate() % 7) + 1

    const row = document.createElement('tr')
    for (let weekday of weekdays) {
      const th = document.createElement('td')
      th.textContent = weekday
      row.appendChild(th)
    }
    this.month.appendChild(row)
    console.log('-')

    for (let i = 0; i < 5; i++) {
      const row = document.createElement("tr")
      for (let j = 0; j < 7; j++) {
        const day = document.createElement('td')
        day.classList.add('day')

        const monthDay = (i * 7 + j - firstWeekdayIndex + 1)
        if (monthDay <= 0 || monthDay > 31) {
          continue
        } 

        const title = document.createTextNode('' + (i * 7 + j - firstWeekdayIndex + 1))
        day.appendChild(title)

        const date = dateNormalize(new Date())
        date.setDate(i * 7 + j - firstWeekdayIndex + 1)

        if (changes.get(date.getTime()) !== undefined) {
          const dayChanges = document.createElement("ul")
          for (let ch of changes.get(date.getTime())) {
            const changeEl = document.createElement('li')
            changeEl.textContent = ch.name
            dayChanges.append(changeEl)
          }
          day.appendChild(dayChanges)
        } else {
          console.log('nop', date, changes.get(date.getTime()))
        }

        row.appendChild(day)
      }
      this.month.appendChild(row)
    }
  }

  setup() {
    this.generateHTML()
  }

  teardown() {}

}

// a change is either an expense (negative) or an income (positive)
// always in integers (cents)
// { "amount": 10, date: "json date", "name": "...", "description": "..." }

