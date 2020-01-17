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

const changes = new Map()

interface Change {
  name: string;
  description: string;
  amount: number;
}

changes.set(addDays(today, -10).getTime(), [
  { name: "the past", description: "", amount: 50 },
])


changes.set(addDays(today, 1).getTime(), [
  { name: "first income", description: "", amount: 200 },
  { name: "first expense", description: "", amount: -20 }
])

changes.set(addDays(today, 2).getTime(), [
  { name: "hello", description: "", amount: -50 }
])

export default class Calendar {
  section: HTMLElement;
  month: HTMLTableElement;
  monthname: HTMLElement;
  report: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.month = qs(this.section, ".month") as HTMLTableElement;
    this.monthname = qs(this.section, ".month-name");
    this.report = qs(this.section, ".report");
    this.generateHTML()

    this.month.addEventListener('click', e => {
      if (!(e.target instanceof HTMLElement)) {
        return
      }
      const closest = e.target.closest('.day')
      if (!closest) {
        return
      }
      const datetime = closest.getAttribute('date')
      if (datetime === null) {
        console.error(closest)
        throw new Error("date attribute is null on .day element")
      }

      const end = dateNormalize(new Date(parseInt(datetime, 10)))
      const start = new Date(end.getTime())
      start.setDate(1)
      const report = this.makeReport(start, end)
      
      this.report.innerHTML = `From ${report.from} to ${report.to}<br><br>Balance: $${report.balance}`
    })
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

    for (let i = 0; i < 5; i++) {
      const row = document.createElement("tr")
      for (let j = 0; j < 7; j++) {
        const day = document.createElement('td')
        day.classList.add('day')

        const monthDay = (i * 7 + j - firstWeekdayIndex + 1)

        if (monthDay >= 1 && monthDay <= 31) {
          const title = document.createTextNode('' + (i * 7 + j - firstWeekdayIndex + 1))
          day.appendChild(title)

          const date = dateNormalize(new Date())
          date.setDate(i * 7 + j - firstWeekdayIndex + 1)
          day.setAttribute('date', "" + date.getTime())

          if (changes.get(date.getTime()) !== undefined) {
            const dayChanges = document.createElement("ul")
            for (let ch of changes.get(date.getTime())) {
              const changeEl = document.createElement('li')
              changeEl.textContent = ch.name
              dayChanges.append(changeEl)
            }
            day.appendChild(dayChanges)
          }
        }

        row.appendChild(day)
      }
      this.month.appendChild(row)
    }
  }

  makeReport(from: Date, to: Date) {
    let selectedchanges: Change[] = []
    for (let [date, daychanges] of changes.entries()) {
      // console.log(date, daychanges)
      if (date < from.getTime() || date > to.getTime()) {
        continue
      }
      selectedchanges.push(...daychanges)
    }

    console.log(selectedchanges)
    let balance = 0;
    for (let change of selectedchanges) {
      balance += change.amount;
    }

    return {
      from: from,
      to: to,
      nchanges: selectedchanges.length,
      balance: balance,
    }
  }

  setup() {
  }

  teardown() { }

}

// a change is either an expense (negative) or an income (positive)
// always in integers (cents)
// { "amount": 10, date: "json date", "name": "...", "description": "..." }

