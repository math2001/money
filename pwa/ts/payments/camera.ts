import { EM, qs } from "../utils.js";

export default class Camera {
  section: HTMLElement;
  video: HTMLVideoElement;
  stream: MediaStream | null;

  startstopbtn: HTMLButtonElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.video = qs(this.section, "video") as HTMLVideoElement;
    this.startstopbtn = qs(this.section, ".start-stop") as HTMLButtonElement;
    this.stream = null;

    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      console.error("no camera API available");
      return;
    }
    qs(this.section, ".no-camera").style.display = "none";

    qs(this.section, ".scan").addEventListener("click", this.scan.bind(this));
    this.startstopbtn.addEventListener("click", this.startstop.bind(this));
  }

  setup() {
    this.startstop();
  }

  async scan() {}

  async startstop() {
    if (this.stream === null) {
      this.startstopbtn.disabled = true;
      await this.startStream();
      this.startstopbtn.textContent = "Stop";
      this.startstopbtn.disabled = false;
    } else {
      this.startstopbtn.textContent = "Start";
      this.stopStream();
    }
  }

  async startStream() {
    if (this.stream !== null) {
      console.error("overwriting existing stream");
    }
    try {
      this.stream = await navigator.mediaDevices.getUserMedia({ video: true });
    } catch (e) {
      console.error(e);
      alert(`Couldn't get video (${e.name}, see console for more details)`);
    }
    console.log("got media");
    this.video.srcObject = this.stream;
  }

  async stopStream() {
    if (this.stream === null) {
      throw new Error("this.stream shouldn't be null");
    }
    this.stream.getTracks().forEach(track => track.stop());
    this.stream = null;
  }

  teardown() {}
}
