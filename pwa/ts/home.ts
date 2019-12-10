namespace Home {
  const section = document.querySelector("#home");
  if (section === null) {
    throw new Error("no home found");
  }

  export function show() {
    console.log("show home!");
    section!.classList.add("active");
  }
}
