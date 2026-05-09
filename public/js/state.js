import { STORAGE_KEYS, THEMES } from "./config.js";

function readNumberFromStorage(key, fallback) {
  const value = Number(localStorage.getItem(key));
  return Number.isFinite(value) && value > 0 ? value : fallback;
}

function readThemeFromStorage() {
  const savedTheme = localStorage.getItem(STORAGE_KEYS.theme);
  return Object.values(THEMES).includes(savedTheme) ? savedTheme : THEMES.light;
}

function readCategoryLimitMap() {
  try {
    const parsed = JSON.parse(localStorage.getItem(STORAGE_KEYS.categoryLimitMap) || "{}");
    return parsed && typeof parsed === "object" ? parsed : {};
  } catch {
    return {};
  }
}

function writeCategoryLimitMap(map) {
  localStorage.setItem(STORAGE_KEYS.categoryLimitMap, JSON.stringify(map));
}

export const state = {
  userId: readNumberFromStorage(STORAGE_KEYS.activeUserId, 1),
  theme: readThemeFromStorage(),
  categories: [],
  transactions: [],
  limits: [],
  users: [],
  stats: null,
  categoryLimitMap: readCategoryLimitMap(),
  filters: {
    transactions: {},
    stats: {},
  },
};

export function setActiveUserId(value) {
  const nextId = Number(value);
  if (!Number.isFinite(nextId) || nextId <= 0) {
    throw new Error("ID пользователя должен быть положительным числом");
  }

  state.userId = nextId;
  localStorage.setItem(STORAGE_KEYS.activeUserId, String(nextId));
}

export function setTheme(theme) {
  state.theme = theme === THEMES.dark ? THEMES.dark : THEMES.light;
  localStorage.setItem(STORAGE_KEYS.theme, state.theme);
  document.documentElement.dataset.theme = state.theme;
}

export function toggleTheme() {
  setTheme(state.theme === THEMES.dark ? THEMES.light : THEMES.dark);
  return state.theme;
}

export function getCategoryId(category) {
  return category?.id ?? category?.ID ?? category?.Id;
}

export function getCategoryTitle(category) {
  return category?.category_name
    ?? category?.categoryName
    ?? category?.CategoryName
    ?? category?.title
    ?? category?.Title
    ?? category?.name
    ?? category?.Name;
}

export function getCategoryUserId(category) {
  return category?.user_id ?? category?.userId ?? category?.UserID ?? category?.UserId;
}

export function getCategoryLimitId(category) {
  const categoryId = getCategoryId(category);
  const backendValue = category?.limit_id
    ?? category?.limitId
    ?? category?.LimitID
    ?? category?.LimitId
    ?? category?.budget_limit_id
    ?? category?.budgetLimitId
    ?? category?.limit?.id
    ?? category?.Limit?.ID;

  return backendValue ?? state.categoryLimitMap[String(categoryId)];
}

export function setCategoryLimitId(categoryId, limitId) {
  if (!categoryId) return;
  const key = String(categoryId);

  if (limitId === undefined || limitId === null || limitId === "") {
    delete state.categoryLimitMap[key];
  } else {
    state.categoryLimitMap[key] = Number(limitId);
  }

  writeCategoryLimitMap(state.categoryLimitMap);
}

export function getCategoryByLimitId(limitId) {
  return state.categories.find((category) => Number(getCategoryLimitId(category)) === Number(limitId));
}

export function getCategoryName(categoryId) {
  const category = state.categories.find((item) => Number(getCategoryId(item)) === Number(categoryId));
  return getCategoryTitle(category) ?? (categoryId ? `Категория #${categoryId}` : "—");
}

export function getLimitId(limit) {
  return limit?.id ?? limit?.ID ?? limit?.Id;
}

export function getLimitAmount(limit) {
  return limit?.amount_limit ?? limit?.amountLimit ?? limit?.AmountLimit ?? limit?.Amount_Limit;
}

export function getLimitDuration(limit) {
  return limit?.duration ?? limit?.Duration;
}

export function getUserId(user) {
  return user?.id ?? user?.ID ?? user?.Id;
}

export function getUserFullName(user) {
  return user?.full_name ?? user?.fullName ?? user?.FullName ?? user?.name ?? user?.Name;
}

export function getUserEmail(user) {
  return user?.email ?? user?.Email;
}

export function getUserPhone(user) {
  return user?.phone ?? user?.Phone;
}

export function getUserTimeAdd(user) {
  return user?.time_add ?? user?.timeAdd ?? user?.TimeAdd;
}
