import { usersApi } from "./api.js?v=11";

import {
  applySavedTheme,
  clearAuthMessage,
  getFormData,
  redirectToApp,
  showAuthMessage,
} from "./auth.js?v=11";

import { assertRequiredString } from "./validators.js?v=11";

applySavedTheme();

const form = document.querySelector("#loginForm");
const button = document.querySelector("#loginBtn");
const message = document.querySelector("#loginMessage");

function validateLoginPayload(data) {
  return {
    email: assertRequiredString(data.email, "Email"),
    password: assertRequiredString(data.password, "Пароль"),
  };
}

form?.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearAuthMessage(message);

  button.disabled = true;
  button.textContent = "Входим...";

  try {
    const payload = validateLoginPayload(getFormData(form));

    await usersApi.login(payload);

    showAuthMessage(message, "Вход выполнен. Открываю приложение.");

    redirectToApp();
  } catch (error) {
    console.error(error);
    showAuthMessage(
        message,
        error.message || "Не удалось войти",
        "error"
    );
  } finally {
    button.disabled = false;
    button.textContent = "Войти";
  }
});