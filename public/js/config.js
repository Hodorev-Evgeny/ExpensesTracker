export const API_BASE = window.APP_CONFIG?.API_BASE ?? "/api/v1";

export const STORAGE_KEYS = {
  activeUserId: "expensesTrackerActiveUserId",
  apiBase: "expensesTrackerApiBase",
  theme: "expensesTrackerTheme",
  categoryLimitMap: "expensesTrackerCategoryLimitMap",
};

export const TRANSACTION_TYPES = {
  expense: "Expenditure",
  income: "Income",
};

export const DEFAULT_PAGE_SIZE = 20;

export const THEMES = {
  light: "light",
  dark: "dark",
};
