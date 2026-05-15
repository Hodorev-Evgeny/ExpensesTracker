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
import { validateUserPayload } from "./validators.js?v=11";

applySavedTheme();

const form = document.querySelector("#registerForm");
const button = document.querySelector("#registerBtn");
const message = document.querySelector("#registerMessage");

function validateRegisterPayload(data) {
  const payload = validateUserPayload(data);
  const passwordRepeat = String(data.password_repeat ?? "");

  if (payload.password !== passwordRepeat) {
    throw new Error("Пароли не совпадают");
  }

  return payload;
}

form?.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearAuthMessage(message);

  button.disabled = true;
  button.textContent = "Создаём...";

  try {
    const payload = validateRegisterPayload(getFormData(form));
    const createdUser = await usersApi.create(payload);
    saveCurrentUserId(getUserIdFromPayload(createdUser));

    try {
      const loginResponse = await usersApi.login({
        email: payload.email,
        password: payload.password,
        full_name: payload.full_name,
      });
      saveCurrentUserId(getUserIdFromPayload(loginResponse) ?? getUserIdFromPayload(createdUser));
      showAuthMessage(message, "Аккаунт создан. Открываю приложение.");
      redirectToApp();
    } catch (loginError) {
      console.warn("Аккаунт создан, но автоматический вход не выполнен", loginError);
      showAuthMessage(message, "Аккаунт создан. Теперь войдите через страницу входа.");
      window.setTimeout(() => {
        window.location.href = "/login.html";
      }, 1000);
    }
  } catch (error) {
    console.error(error);
    showAuthMessage(message, error.message || "Не удалось зарегистрироваться", "error");
  } finally {
    button.disabled = false;
    button.textContent = "Зарегистрироваться";
  }
});
