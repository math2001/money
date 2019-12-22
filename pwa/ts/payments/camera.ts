import { EM, qs } from "../utils.js";

export default class Camera {
  section: HTMLElement;
  video: HTMLVideoElement;
  stream: MediaStream | null;
  canvas: HTMLCanvasElement;

  startstopbtn: HTMLButtonElement;

  constructor(section: HTMLElement) {
    this.section = section;
    this.video = qs(this.section, "video") as HTMLVideoElement;
    this.startstopbtn = qs(this.section, ".start-stop") as HTMLButtonElement;
    this.canvas = qs(this.section, "canvas") as HTMLCanvasElement;
    this.stream = null;

    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      console.error("no camera API available");
      return;
    }
    qs(this.section, ".no-camera-error").style.display = "none";
    qs(this.section, ".no-camera-api").style.display = "none";

    qs(this.section, ".scan").addEventListener("click", this.scan.bind(this));
    this.startstopbtn.addEventListener("click", this.startstop.bind(this));
  }

  setup() {
    this.startstop();
  }

  async scan() {
    this.canvas.width = this.video.videoWidth;
    this.canvas.height = this.video.videoHeight;
    const context = this.canvas.getContext("2d");

    if (context === null) {
      throw new Error("shouldn't have null context (canvas)");
    }
    context.drawImage(this.video, 0, 0);

    // send the image the server
    const formdata = new FormData();
    formdata.append("img", await toBlob(this.canvas, "image/png", 0));

    console.info("scanning image");
    const resp = await fetch("/api/payments/scan", {
      method: "post",
      body: formdata,
    });
    const obj = await resp.json();
    console.log(obj.payment.result);
  }

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
      qs(this.section, ".no-camera-error").style.display = "block";
    }
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

function toBlob(
  canvas: HTMLCanvasElement,
  mimeType: string,
  qualityArgument: number,
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob(resolve as BlobCallback, mimeType, qualityArgument);
  });
}
