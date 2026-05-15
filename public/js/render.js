import { TRANSACTION_TYPES } from "./config.js?v=11";
import {
  getCategoryByLimitId,
  getCategoryId,
  getCategoryLimitId,
  getCategoryName,
  getCategoryTitle,
  getLimitAmount,
  getLimitDuration,
  getLimitId,
  getUserEmail,
  getUserFullName,
  getUserId,
  getUserPhone,
  getUserTimeAdd,
  state,
} from "./state.js?v=11";
import { nowToLocalDateTimeInputValue, toLocalDateTimeInputValue } from "./validators.js?v=11";

const formatter = new Intl.NumberFormat("ru-RU", {
  style: "currency",
  currency: "RUB",
  maximumFractionDigits: 0,
});

export function $(selector, root = document) {
  return root.querySelector(selector);
}

export function $$(selector, root = document) {
  return [...root.querySelectorAll(selector)];
}

export function formatMoney(value) {
  return formatter.format(Number(value) || 0);
}

export function formatDate(value) {
  if (!value) return "—";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "—";

  return date.toLocaleString("ru-RU", {
    dateStyle: "short",
    timeStyle: "short",
  });
}

export function escapeHtml(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

export function getTransactionTypeLabel(type) {
  if (type === TRANSACTION_TYPES.expense || String(type).toLowerCase().includes("expenditure")) {
    return "Расход";
  }

  if (type === TRANSACTION_TYPES.income || String(type).toLowerCase().includes("income")) {
    return "Доход";
  }

  return type || "—";
}

export function getTransactionTypeBadge(type) {
  const label = getTransactionTypeLabel(type);
  const className = label === "Расход"
    ? "badge badge-expense"
    : label === "Доход"
      ? "badge badge-income"
      : "badge badge-neutral";

  return `<span class="${className}">${escapeHtml(label)}</span>`;
}

export function showToast(message, type = "success") {
  const root = $("#toastRoot");
  if (!root) return;

  const toast = document.createElement("div");
  toast.className = `toast toast-${type}`;
  toast.textContent = message;

  root.append(toast);

  window.setTimeout(() => {
    toast.remove();
  }, 4200);
}

export function setLoading(element, isLoading) {
  if (!element) return;
  element.classList.toggle("loading", isLoading);
}

export function updateThemeButton() {
  const button = $("#themeToggleBtn");
  if (!button) return;

  button.textContent = state.theme === "dark" ? "Светлая тема" : "Тёмная тема";
  button.setAttribute("aria-pressed", String(state.theme === "dark"));
}

export function renderCategorySelects() {
  const selects = [
    $("#transactionCategoryId"),
    $("#filterCategoryId"),
    $("#statsCategoryId"),
    $("#limitCategoryId"),
  ].filter(Boolean);

  selects.forEach((select) => {
    const currentValue = select.value;
    const firstOptionText = select.id === "transactionCategoryId" || select.id === "limitCategoryId"
      ? "Выберите категорию"
      : "Все категории";

    select.innerHTML = `<option value="">${firstOptionText}</option>`;

    state.categories.forEach((category) => {
      const categoryId = getCategoryId(category);
      if (!categoryId) return;

      const option = document.createElement("option");
      option.value = categoryId;
      option.textContent = getCategoryTitle(category) ?? `Категория #${categoryId}`;
      select.append(option);
    });

    if ([...select.options].some((option) => option.value === currentValue)) {
      select.value = currentValue;
    }
  });
}

export function renderCategories() {
  const tbody = $("#categoriesTbody");
  const count = $("#categoriesCount");
  const template = $("#categoryRowTemplate");

  const categories = Array.isArray(state.categories) ? state.categories : [];

  if (count) count.textContent = `${categories.length} записей`;
  if (!tbody || !template) return;

  tbody.innerHTML = "";

  if (!categories.length) {
    tbody.innerHTML = `<tr><td colspan="4" class="empty-cell">Категории пока не созданы.</td></tr>`;
    return;
  }

  categories.forEach((category) => {
    const categoryId = getCategoryId(category);
    const limitId = getCategoryLimitId(category);
    const row = template.content.firstElementChild.cloneNode(true);
    row.dataset.id = categoryId ?? "";
    row.querySelector('[data-cell="id"]').textContent = categoryId ?? "—";
    row.querySelector('[data-cell="category_name"]').textContent = getCategoryTitle(category) ?? "—";
    row.querySelector('[data-cell="limit_id"]').innerHTML = limitId
      ? `<span class="badge badge-positive">Лимит #${escapeHtml(limitId)}</span>`
      : `<span class="badge badge-neutral">Без лимита</span>`;
    tbody.append(row);
  });
}

export function renderTransactions(transactions = state.transactions) {
  const tbody = $("#transactionsTbody");
  const count = $("#transactionsCount");
  const template = $("#transactionRowTemplate");

  const list = Array.isArray(transactions)
    ? transactions
    : transactions
      ? [transactions]
      : [];

  if (count) count.textContent = `${list.length} записей`;
  if (!tbody || !template) return;

  tbody.innerHTML = "";

  if (!list.length) {
    tbody.innerHTML = `<tr><td colspan="8" class="empty-cell">Операции не найдены.</td></tr>`;
    return;
  }

  list.forEach((transaction) => {
    const transactionId = transaction.id ?? transaction.ID ?? transaction.Id;
    const categoryId = transaction.category_id ?? transaction.categoryId ?? transaction.CategoryID ?? transaction.CategoryId;
    const timeChange = transaction.time_change ?? transaction.timeChange ?? transaction.TimeChange;

    const row = template.content.firstElementChild.cloneNode(true);
    row.dataset.id = transactionId ?? "";
    row.querySelector('[data-cell="id"]').textContent = transactionId ?? "—";
    row.querySelector('[data-cell="date"]').textContent = formatDate(transaction.date ?? transaction.Date);
    row.querySelector('[data-cell="type"]').innerHTML = getTransactionTypeBadge(transaction.type ?? transaction.Type);
    row.querySelector('[data-cell="category"]').textContent = getCategoryName(categoryId);
    row.querySelector('[data-cell="sum"]').textContent = formatMoney(transaction.sum ?? transaction.Sum);
    row.querySelector('[data-cell="comments"]').textContent = transaction.comments ?? transaction.Comments ?? "—";
    row.querySelector('[data-cell="time_change"]').textContent = formatDate(timeChange);
    tbody.append(row);
  });
}

function getExpenseSumByCategory(categoryId, duration) {
  const durationTime = duration ? new Date(duration).getTime() : null;

  return state.transactions.reduce((sum, transaction) => {
    const transactionCategoryId = transaction.category_id ?? transaction.categoryId ?? transaction.CategoryID ?? transaction.CategoryId;
    const transactionType = transaction.type ?? transaction.Type;

    if (Number(transactionCategoryId) !== Number(categoryId)) return sum;
    if (!String(transactionType).toLowerCase().includes("expenditure")) return sum;

    if (durationTime) {
      const transactionTime = new Date(transaction.date ?? transaction.Date).getTime();
      if (!Number.isNaN(transactionTime) && transactionTime > durationTime) return sum;
    }

    return sum + Number(transaction.sum ?? transaction.Sum ?? 0);
  }, 0);
}

function getLimitStatus(spent, amount) {
  if (!amount) {
    return { label: "Нет суммы", className: "badge-neutral" };
  }

  const progress = spent / amount;

  if (progress >= 1) {
    return { label: "Лимит превышен", className: "badge-expense" };
  }

  if (progress >= 0.8) {
    return { label: "Почти достигнут", className: "badge-warning" };
  }

  return { label: "Не достигнут", className: "badge-positive" };
}

export function renderLimits() {
  const tbody = $("#limitsTbody");
  const count = $("#limitsCount");
  const template = $("#limitRowTemplate");

  const rawLimits = Array.isArray(state.limits) ? state.limits : [];
  const currentLimitIds = new Set(
    state.categories
      .map((category) => Number(getCategoryLimitId(category)))
      .filter((id) => Number.isFinite(id) && id > 0),
  );
  const currentCategoryIds = new Set(
    state.categories
      .map((category) => Number(getCategoryId(category)))
      .filter((id) => Number.isFinite(id) && id > 0),
  );
  const limits = rawLimits.filter((limit) => {
    const limitId = Number(getLimitId(limit));
    const linkedCategory = getCategoryByLimitId(limitId);
    const categoryId = Number(limit.category_id ?? limit.categoryId ?? limit.CategoryID ?? limit.CategoryId);
    const userId = Number(limit.user_id ?? limit.userId ?? limit.UserID ?? limit.UserId);

    if (Number.isFinite(userId) && userId > 0) return Number(userId) === Number(state.userId);
    if (linkedCategory) return true;
    if (Number.isFinite(categoryId) && categoryId > 0) return currentCategoryIds.has(categoryId);
    return Number.isFinite(limitId) && currentLimitIds.has(limitId);
  });

  if (count) count.textContent = `${limits.length} записей`;
  if (!tbody || !template) return;

  tbody.innerHTML = "";

  if (!limits.length) {
    tbody.innerHTML = `<tr><td colspan="7" class="empty-cell">Лимиты пока не созданы.</td></tr>`;
    return;
  }

  limits.forEach((limit) => {
    const limitId = getLimitId(limit);
    const amount = Number(getLimitAmount(limit)) || 0;
    const duration = getLimitDuration(limit);
    const linkedCategory = getCategoryByLimitId(limitId);
    const categoryId = linkedCategory
      ? getCategoryId(linkedCategory)
      : limit.category_id ?? limit.categoryId ?? limit.CategoryID ?? limit.CategoryId;
    const spent = categoryId ? getExpenseSumByCategory(categoryId, duration) : 0;
    const status = categoryId
      ? getLimitStatus(spent, amount)
      : { label: "Не привязан", className: "badge-neutral" };

    const row = template.content.firstElementChild.cloneNode(true);
    row.dataset.id = limitId ?? "";
    row.querySelector('[data-cell="id"]').textContent = limitId ?? "—";
    row.querySelector('[data-cell="category"]').textContent = categoryId ? getCategoryName(categoryId) : "Не привязан";
    row.querySelector('[data-cell="amount_limit"]').textContent = formatMoney(amount);
    row.querySelector('[data-cell="duration"]').textContent = formatDate(duration);
    row.querySelector('[data-cell="spent"]').textContent = categoryId ? formatMoney(spent) : "—";
    row.querySelector('[data-cell="status"]').innerHTML = `<span class="badge ${status.className}">${status.label}</span>`;
    tbody.append(row);
  });
}

function getActiveUser() {
  return state.currentUser;
}

function getInitials(name) {
  const parts = String(name || "")
    .trim()
    .split(/\s+/)
    .filter(Boolean);

  if (!parts.length) return "?";

  return parts
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("");
}

export function renderProfile() {
  const user = getActiveUser();
  const profileName = $("#profileName");
  const profileEmail = $("#profileEmail");
  const profileEmailMeta = $("#profileEmailMeta");
  const profilePhone = $("#profilePhone");
  const profileCreated = $("#profileCreated");
  const profileAvatar = $("#profileAvatar");
  const profileBadge = $("#profileActiveBadge");
  const headerName = $("#headerProfileName");
  const headerEmail = $("#headerProfileEmail");
  const headerAvatar = $("#headerAvatar");

  if (!user) {
    if (profileName) profileName.textContent = "Профиль не загружен";
    if (profileEmail) profileEmail.textContent = "Войдите в аккаунт, чтобы увидеть личную информацию.";
    if (profileEmailMeta) profileEmailMeta.textContent = "—";
    if (profilePhone) profilePhone.textContent = "—";
    if (profileCreated) profileCreated.textContent = "—";
    if (profileAvatar) profileAvatar.textContent = "?";
    if (headerName) headerName.textContent = "Профиль не загружен";
    if (headerEmail) headerEmail.textContent = "Нужен вход";
    if (headerAvatar) headerAvatar.textContent = "?";
    if (profileBadge) {
      profileBadge.textContent = "Нет сессии";
      profileBadge.className = "badge badge-warning";
    }
    return;
  }

  const name = getUserFullName(user) ?? "Пользователь";
  const email = getUserEmail(user) ?? "Email не указан";
  const initials = getInitials(name);

  if (profileName) profileName.textContent = name;
  if (profileEmail) profileEmail.textContent = email;
  if (profileEmailMeta) profileEmailMeta.textContent = email;
  if (profilePhone) profilePhone.textContent = getUserPhone(user) ?? "—";
  if (profileCreated) profileCreated.textContent = formatDate(getUserTimeAdd(user));
  if (profileAvatar) profileAvatar.textContent = initials;
  if (headerName) headerName.textContent = name;
  if (headerEmail) headerEmail.textContent = email;
  if (headerAvatar) headerAvatar.textContent = initials;
  if (profileBadge) {
    profileBadge.textContent = "Активный";
    profileBadge.className = "badge badge-positive";
  }
}

export function renderUsers() {
  renderProfile();
}

export function renderStats() {
  const stats = state.stats || {};

  $("#statIncome").textContent = formatMoney(stats.sum_income);
  $("#statExpense").textContent = formatMoney(stats.sum_expenditure);
  $("#statDifference").textContent = formatMoney(stats.difference);
  $("#statCount").textContent = stats.count_operation ?? 0;
  $("#statAvgExpense").textContent = formatMoney(stats.avg_expenditure);
  $("#statAvgIncome").textContent = formatMoney(stats.avg_income);
  $("#statCostCategory").textContent = stats.cost_category || "—";
  $("#statMaxExpense").textContent = formatMoney(stats.max_expenditure);

  const shareList = $("#categoryShareList");
  if (!shareList) return;

  const entries = Object.entries(stats.share_category || {});
  shareList.innerHTML = "";

  if (!entries.length) {
    shareList.innerHTML = `<p class="empty-state">За выбранный период нет данных по категориям.</p>`;
    return;
  }

  entries
    .sort(([, a], [, b]) => Number(b) - Number(a))
    .forEach(([name, percent]) => {
      const safePercent = Math.max(0, Math.min(100, Number(percent) || 0));
      const row = document.createElement("div");
      row.className = "share-row";
      row.innerHTML = `
        <strong>${escapeHtml(name)}</strong>
        <div class="share-bar" aria-label="${escapeHtml(name)}: ${safePercent.toFixed(1)}%">
          <span style="width: ${safePercent}%"></span>
        </div>
        <span>${safePercent.toFixed(1)}%</span>
      `;
      shareList.append(row);
    });
}

export function fillTransactionForm(transaction) {
  $("#transactionId").value = transaction.id ?? transaction.ID ?? transaction.Id ?? "";
  $("#transactionType").value = transaction.type ?? transaction.Type ?? "Expenditure";
  $("#transactionSum").value = transaction.sum ?? transaction.Sum ?? "";
  $("#transactionDate").value = toLocalDateTimeInputValue(transaction.date ?? transaction.Date);
  $("#transactionCategoryId").value = transaction.category_id ?? transaction.categoryId ?? transaction.CategoryID ?? transaction.CategoryId ?? "";
  $("#transactionComments").value = transaction.comments ?? transaction.Comments ?? "";
  $("#saveTransactionBtn").textContent = "Сохранить изменения";
}

export function fillCategoryForm(category) {
  $("#categoryId").value = getCategoryId(category) ?? "";
  $("#categoryName").value = getCategoryTitle(category) ?? "";
}

export function fillLimitForm(limit) {
  const limitId = getLimitId(limit);
  const linkedCategory = getCategoryByLimitId(limitId);

  $("#limitId").value = limitId ?? "";
  $("#limitCategoryId").value = linkedCategory ? getCategoryId(linkedCategory) : "";
  $("#limitAmount").value = getLimitAmount(limit) ?? "";
  $("#limitDuration").value = toLocalDateTimeInputValue(getLimitDuration(limit));
}

export function fillUserForm(user) {
  $("#userId").value = getUserId(user) ?? "";
  $("#userFullName").value = getUserFullName(user) ?? "";
  $("#userEmail").value = getUserEmail(user) ?? "";
  $("#userPhone").value = getUserPhone(user) ?? "";
  $("#saveUserBtn").textContent = "Сохранить профиль";
}

export function resetTransactionForm() {
  $("#transactionForm").reset();
  $("#transactionId").value = "";
  $("#saveTransactionBtn").textContent = "Сохранить";
  $("#transactionDate").value = nowToLocalDateTimeInputValue();
}

export function resetCategoryForm() {
  $("#categoryForm").reset();
  $("#categoryId").value = "";
}

export function resetLimitForm() {
  $("#limitForm").reset();
  $("#limitId").value = "";
}

export function resetUserForm() {
  fillUserForm(state.currentUser || {});
  $("#saveUserBtn").textContent = "Сохранить профиль";
}
