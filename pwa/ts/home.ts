import { State, qs } from "./utils.js";

export default class Home {
  section: HTMLElement;
  foruser: HTMLElement;
  forvisitor: HTMLElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.foruser = qs(this.section, ".foruser");
    this.forvisitor = qs(this.section, ".forvisitor");
  }

  setup() {
    if (State.user && State.user.email !== null) {
      this.forvisitor.style.display = "none";
      this.foruser.style.display = "block";
    } else {
      this.forvisitor.style.display = "block";
      this.foruser.style.display = "none";
    }
  }

  teardown() {}
}
