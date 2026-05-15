import { usersApi } from "./api.js?v=11";
import {
  applySavedTheme,
  clearAuthMessage,
  getFormData,
  getUserIdFromPayload,
  redirectToApp,
  saveCurrentUserId,
  showAuthMessage,
} from "./auth.js?v=11";
import { assertRequiredString } from "./validators.js?v=11";

applySavedTheme();

const form = document.querySelector("#loginForm");
const button = document.querySelector("#loginBtn");
const message = document.querySelector("#loginMessage");

function validateLoginPayload(data) {
  const email = assertRequiredString(data.email, "Email");
  const password = assertRequiredString(data.password, "Пароль");
  const fullName = String(data.full_name ?? "").trim();

  return {
    email,
    password,
    ...(fullName ? { full_name: fullName } : {}),
  };
}

form?.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearAuthMessage(message);

  button.disabled = true;
  button.textContent = "Входим...";

  try {
    const payload = validateLoginPayload(getFormData(form));
    const loginResponse = await usersApi.login(payload);

    let userId = getUserIdFromPayload(loginResponse);

    if (!userId) {
      try {
        const users = await usersApi.list({ limit: 100, offset: 0 });
        const currentUser = users.find((user) => String(user.email ?? user.Email ?? "").toLowerCase() === payload.email.toLowerCase());
        userId = getUserIdFromPayload(currentUser);
      } catch (error) {
        console.warn("Не удалось получить user id после входа", error);
      }
    }

    saveCurrentUserId(userId);
    showAuthMessage(message, "Вход выполнен. Открываю приложение.");
    redirectToApp();
  } catch (error) {
    console.error(error);
    showAuthMessage(message, error.message || "Не удалось войти", "error");
  } finally {
    button.disabled = false;
    button.textContent = "Войти";
  }
});
