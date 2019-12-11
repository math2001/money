export default class Home {
  section: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;
  }

  setup() {
    this.section.classList.add("active");
  }

  teardown() {
    this.section.classList.remove("active");
  }
}
