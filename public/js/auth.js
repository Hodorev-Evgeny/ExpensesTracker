import { STORAGE_KEYS, THEMES } from "./config.js?v=11";

export function applySavedTheme() {
  const savedTheme = localStorage.getItem(STORAGE_KEYS.theme);
  document.documentElement.dataset.theme = Object.values(THEMES).includes(savedTheme) ? savedTheme : THEMES.light;
}

export function getFormData(formElement) {
  return Object.fromEntries(new FormData(formElement).entries());
}

export function showAuthMessage(element, text, type = "success") {
  if (!element) return;
  element.textContent = text;
  element.className = `auth-message is-${type}`;
}

export function clearAuthMessage(element) {
  if (!element) return;
  element.textContent = "";
  element.className = "auth-message";
}

export function getAuthErrorMessage(payload, fallback) {
  if (!payload) return fallback;
  return payload.message || payload.massage || payload.error || payload.detail || fallback;
}

export function redirectToApp(delay = 700) {
  window.setTimeout(() => {
    window.location.href = "/";
  }, delay);
}


export function getUserIdFromPayload(payload) {
  return payload?.id
    ?? payload?.ID
    ?? payload?.Id
    ?? payload?.user_id
    ?? payload?.userId
    ?? payload?.userID
    ?? payload?.UserID;
}

export function saveCurrentUserId(userId) {
  const numericId = Number(userId);
  if (Number.isFinite(numericId) && numericId > 0) {
    localStorage.setItem(STORAGE_KEYS.currentUserId, String(numericId));
  }
}

export function readCurrentUserId() {
  const raw = localStorage.getItem(STORAGE_KEYS.currentUserId);
  const numericId = Number(raw);
  return Number.isFinite(numericId) && numericId > 0 ? numericId : null;
}

export function clearCurrentUserId() {
  localStorage.removeItem(STORAGE_KEYS.currentUserId);
}
