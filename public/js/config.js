export const API_BASE = window.APP_CONFIG?.API_BASE
  ?? localStorage.getItem("expensesTrackerApiBase")
  ?? "http://127.0.0.1:8080/api/v1";

export const STORAGE_KEYS = {
  activeUserId: "expensesTrackerActiveUserId",
  apiBase: "expensesTrackerApiBase",
};

export const TRANSACTION_TYPES = {
  expense: "Expenditure",
  income: "Income",
};

export const DEFAULT_PAGE_SIZE = 20;
