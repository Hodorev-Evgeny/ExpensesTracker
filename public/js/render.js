import { getCategoryId, getCategoryName, getCategoryTitle, getCategoryUserId, state } from "./state.js";
import { TRANSACTION_TYPES } from "./config.js";
import { nowToLocalDateTimeInputValue, toLocalDateTimeInputValue } from "./validators.js";

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

export function renderCategorySelects() {
  const selects = [
    $("#transactionCategoryId"),
    $("#filterCategoryId"),
    $("#statsCategoryId"),
    $("#limitCategoryId"),
  ].filter(Boolean);

  selects.forEach((select) => {
    const currentValue = select.value;
    const firstOptionText = select.id === "transactionCategoryId"
      ? "Выберите категорию"
      : select.id === "limitCategoryId"
        ? "Привязка к категории пока не описана в API"
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
    const row = template.content.firstElementChild.cloneNode(true);
    row.dataset.id = categoryId ?? "";
    row.querySelector('[data-cell="id"]').textContent = categoryId ?? "—";
    row.querySelector('[data-cell="category_name"]').textContent = getCategoryTitle(category) ?? "—";
    row.querySelector('[data-cell="user_id"]').textContent = getCategoryUserId(category) ?? "—";
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

export function renderLimits() {
  const tbody = $("#limitsTbody");
  const count = $("#limitsCount");
  const template = $("#limitRowTemplate");

  if (count) count.textContent = `${state.limits.length} записей`;
  if (!tbody || !template) return;

  tbody.innerHTML = "";

  if (!state.limits.length) {
    tbody.innerHTML = `<tr><td colspan="5" class="empty-cell">Лимиты пока не созданы.</td></tr>`;
    return;
  }

  state.limits.forEach((limit) => {
    const row = template.content.firstElementChild.cloneNode(true);
    row.dataset.id = limit.id;
    row.querySelector('[data-cell="id"]').textContent = limit.id ?? "—";
    row.querySelector('[data-cell="amount_limit"]').textContent = formatMoney(limit.amount_limit);
    row.querySelector('[data-cell="duration"]').textContent = formatDate(limit.duration);
    row.querySelector('[data-cell="status"]').innerHTML = `<span class="badge badge-neutral">Нет данных</span>`;
    tbody.append(row);
  });
}

export function renderUsers() {
  const tbody = $("#usersTbody");
  const count = $("#usersCount");
  const template = $("#userRowTemplate");

  if (count) count.textContent = `${state.users.length} записей`;
  if (!tbody || !template) return;

  tbody.innerHTML = "";

  if (!state.users.length) {
    tbody.innerHTML = `<tr><td colspan="6" class="empty-cell">Пользователи не найдены.</td></tr>`;
    return;
  }

  state.users.forEach((user) => {
    const row = template.content.firstElementChild.cloneNode(true);
    row.dataset.id = user.id;
    row.querySelector('[data-cell="id"]').textContent = user.id ?? "—";
    row.querySelector('[data-cell="full_name"]').textContent = user.full_name ?? "—";
    row.querySelector('[data-cell="email"]').textContent = user.email ?? "—";
    row.querySelector('[data-cell="phone"]').textContent = user.phone ?? "—";
    row.querySelector('[data-cell="time_add"]').textContent = formatDate(user.time_add);
    tbody.append(row);
  });
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
  $("#transactionId").value = transaction.id ?? "";
  $("#transactionType").value = transaction.type ?? "Expenditure";
  $("#transactionSum").value = transaction.sum ?? "";
  $("#transactionDate").value = toLocalDateTimeInputValue(transaction.date);
  $("#transactionCategoryId").value = transaction.category_id ?? "";
  $("#transactionComments").value = transaction.comments ?? "";
  $("#saveTransactionBtn").textContent = "Сохранить изменения";
}

export function fillCategoryForm(category) {
  $("#categoryId").value = getCategoryId(category) ?? "";
  $("#categoryName").value = getCategoryTitle(category) ?? "";
}

export function fillLimitForm(limit) {
  $("#limitId").value = limit.id ?? "";
  $("#limitAmount").value = limit.amount_limit ?? "";
  $("#limitDuration").value = toLocalDateTimeInputValue(limit.duration);
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
