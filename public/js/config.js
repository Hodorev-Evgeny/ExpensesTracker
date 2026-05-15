export const API_BASE = window.APP_CONFIG?.API_BASE ?? "/api/v1";

export const STORAGE_KEYS = {
  theme: "expensesTrackerTheme",
  categoryLimitMap: "expensesTrackerCategoryLimitMap",
  currentUserId: "expensesTrackerCurrentUserId",
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
