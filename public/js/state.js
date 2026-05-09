import { DEFAULT_PAGE_SIZE, STORAGE_KEYS } from "./config.js";

const savedUserId = Number(localStorage.getItem(STORAGE_KEYS.activeUserId));

export const state = {
  userId: Number.isInteger(savedUserId) && savedUserId > 0 ? savedUserId : 1,
  categories: [],
  transactions: [],
  limits: [],
  users: [],
  stats: null,
  filters: {
    transactions: {
      limit: DEFAULT_PAGE_SIZE,
      offset: 0,
    },
  },
};

export function setActiveUserId(userId) {
  const normalizedUserId = Number(userId);

  if (!Number.isInteger(normalizedUserId) || normalizedUserId <= 0) {
    throw new Error("ID пользователя должен быть положительным целым числом");
  }

  state.userId = normalizedUserId;
  localStorage.setItem(STORAGE_KEYS.activeUserId, String(normalizedUserId));
}

export function getCategoryId(category) {
  return category?.id ?? category?.ID ?? category?.Id;
}

export function getCategoryTitle(category) {
  return (
    category?.category_name ??
    category?.categoryName ??
    category?.title ??
    category?.name ??
    category?.CategoryName ??
    category?.Name
  );
}

export function getCategoryUserId(category) {
  return category?.user_id ?? category?.userId ?? category?.UserID ?? category?.UserId;
}

export function getCategoryName(categoryId) {
  const category = state.categories.find((item) => Number(getCategoryId(item)) === Number(categoryId));
  return getCategoryTitle(category) ?? `Категория #${categoryId ?? "—"}`;
}
