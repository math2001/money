import Home from "home";

function initApp() {
  const home = new Home();
  const login = initLogin();
  const err404 = initErr404();

  return {
    show: function(pathname: string) {
      switch (pathname) {
        case "/":
          home.show();
          break;
        case "/login":
          login.show();
          break;
        default:
          err404.show();
      }
    }
  };
}

function initLogin() {
  const section = document.querySelector("#login");
  if (section === null) {
    throw new Error("no section found");
  }

  const form = section.querySelector("form.login-form");
  if (form === null) {
    throw new Error("no form found");
  }

  return {
    show: () => {
      form.addEventListener("submit", (e: Event) => {
        alert("submit form");
        e.preventDefault();
      });
      section.classList.add("active");
    }
  };
}

function initErr404() {
  return {
    show: () => {}
  };
}

document.addEventListener("DOMContentLoaded", () => {
  const app = initApp();
  app.show(location.pathname);
});
