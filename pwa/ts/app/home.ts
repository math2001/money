export default class Home {
  section: Element;

  constructor() {
    const section = document.querySelector("#home");
    if (section === null) {
      throw new Error("no home found");
    }
    this.section = section;
  }

  show() {
    console.log("show home!");
    this.section.classList.add("active");
  }
}
