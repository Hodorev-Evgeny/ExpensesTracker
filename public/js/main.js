import { categoriesApi, limitsApi, statsApi, transactionsApi, usersApi } from "./api.js";
import { DEFAULT_PAGE_SIZE } from "./config.js";
import { getCategoryId, getCategoryTitle, setActiveUserId, state } from "./state.js";
import {
  assertDateTimeLocal,
  assertPositiveInteger,
  assertTransactionType,
  dateInputToRFC3339,
  nowToLocalDateTimeInputValue,
  toApiDateTime,
  validateCategoryName,
  validateUserPayload,
} from "./validators.js";
import {
  $,
  $$,
  fillCategoryForm,
  fillLimitForm,
  fillTransactionForm,
  renderCategories,
  renderCategorySelects,
  renderLimits,
  renderStats,
  renderTransactions,
  renderUsers,
  resetCategoryForm,
  resetLimitForm,
  resetTransactionForm,
  setLoading,
  showToast,
} from "./render.js";

function getFormData(form) {
  return Object.fromEntries(new FormData(form).entries());
}

function onlyFilled(object) {
  return Object.fromEntries(
    Object.entries(object).filter(([, value]) => value !== undefined && value !== null && value !== ""),
  );
}

async function runSafely(action, successMessage) {
  try {
    await action();
    if (successMessage) showToast(successMessage);
  } catch (error) {
    console.error(error);
    showToast(error.message || "Произошла ошибка", "error");
  }
}

function confirmAction({ title = "Подтвердите действие", text = "Вы уверены?" } = {}) {
  const dialog = $("#confirmDialog");
  const titleNode = $("#confirmTitle");
  const textNode = $("#confirmText");

  if (!dialog || typeof dialog.showModal !== "function") {
    return Promise.resolve(window.confirm(text));
  }

  titleNode.textContent = title;
  textNode.textContent = text;

  return new Promise((resolve) => {
    dialog.addEventListener(
      "close",
      () => {
        resolve(dialog.returnValue === "ok");
      },
      { once: true },
    );

    dialog.showModal();
  });
}

function activateTab(tabName) {
  $$(".tab").forEach((tab) => {
    tab.classList.toggle("is-active", tab.dataset.tab === tabName);
  });

  $$(".page").forEach((page) => {
    page.classList.toggle("is-active", page.id === tabName);
  });
}

function prepareTransactionQuery(filterData = {}) {
  return onlyFilled({
    user_id: state.userId,
    category_id: filterData.category_id,
    sum: filterData.sum,
    from: dateInputToRFC3339(filterData.from),
    to: dateInputToRFC3339(filterData.to, true),
    limit: filterData.limit || DEFAULT_PAGE_SIZE,
    offset: filterData.offset || 0,
  });
}

function applyLocalTransactionFilters(transactions, filterData = {}) {
  let filtered = Array.isArray(transactions) ? [...transactions] : [];

  filtered = filtered.filter((item) => {
    const itemUserId = item.user_id ?? item.userId ?? item.UserID ?? item.UserId;

    if (!itemUserId) {
      return true;
    }

    return Number(itemUserId) === Number(state.userId);
  });

  if (filterData.type) {
    filtered = filtered.filter((item) => {
      const type = item.type ?? item.Type;
      return type === filterData.type;
    });
  }

  if (filterData.category_id) {
    filtered = filtered.filter((item) => {
      const categoryId =
          item.category_id ??
          item.categoryId ??
          item.CategoryID ??
          item.CategoryId;

      return Number(categoryId) === Number(filterData.category_id);
    });
  }

  if (filterData.sum) {
    filtered = filtered.filter((item) => {
      const sum = item.sum ?? item.Sum;
      return Number(sum) >= Number(filterData.sum);
    });
  }

  if (filterData.from) {
    const from = new Date(dateInputToRFC3339(filterData.from)).getTime();

    filtered = filtered.filter((item) => {
      const date = new Date(item.date ?? item.Date).getTime();
      return !Number.isNaN(date) && date >= from;
    });
  }

  if (filterData.to) {
    const to = new Date(dateInputToRFC3339(filterData.to, true)).getTime();

    filtered = filtered.filter((item) => {
      const date = new Date(item.date ?? item.Date).getTime();
      return !Number.isNaN(date) && date <= to;
    });
  }

  if (filterData.comment) {
    const needle = filterData.comment.trim().toLowerCase();

    filtered = filtered.filter((item) => {
      const comments = item.comments ?? item.Comments ?? "";
      return String(comments).toLowerCase().includes(needle);
    });
  }

  filtered.sort((a, b) => {
    const dateA = new Date(a.date ?? a.Date ?? 0).getTime();
    const dateB = new Date(b.date ?? b.Date ?? 0).getTime();

    return dateB - dateA;
  });

  const limit = Number(filterData.limit || DEFAULT_PAGE_SIZE);
  const offset = Number(filterData.offset || 0);

  return filtered.slice(offset, offset + limit);
}

async function loadCategories() {
  state.categories = await categoriesApi.list();
  renderCategories();
  renderCategorySelects();
}

async function loadTransactions(filterData = state.filters.transactions) {
  state.filters.transactions = {
    limit: DEFAULT_PAGE_SIZE,
    offset: 0,
    ...filterData,
  };

  const transactions = await transactionsApi.list();

  state.transactions = Array.isArray(transactions)
      ? transactions
      : transactions
          ? [transactions]
          : [];

  const visibleTransactions = applyLocalTransactionFilters(
      state.transactions,
      state.filters.transactions,
  );

  renderTransactions(visibleTransactions);
}

async function loadStats(filterData = {}) {
  const query = onlyFilled({
    from: dateInputToRFC3339(filterData.from),
    to: dateInputToRFC3339(filterData.to, true),
    category_id: filterData.category_id,
  });

  state.stats = await statsApi.get(query);
  renderStats();
}

async function loadLimits() {
  state.limits = await limitsApi.list({ limit: 100, offset: 0 });
  renderLimits();
}

async function loadUsers() {
  state.users = await usersApi.list({ limit: 100, offset: 0 });
  renderUsers();
}

async function reloadMainData() {
  const app = $(".app");
  setLoading(app, true);

  try {
    await loadCategories();

    const results = await Promise.allSettled([
      loadTransactions(),
      loadLimits(),
      loadUsers(),
      loadStats(getFormData($("#statsFilterForm"))),
    ]);

    results.forEach((result) => {
      if (result.status === "rejected") {
        console.error(result.reason);
      }
    });
  } finally {
    setLoading(app, false);
  }
}

function setupTabs() {
  $$(".tab").forEach((tab) => {
    tab.addEventListener("click", () => activateTab(tab.dataset.tab));
  });
}

function setupUserPanel() {
  const input = $("#activeUserId");
  input.value = state.userId;

  $("#saveUserIdBtn").addEventListener("click", () => {
    runSafely(async () => {
      setActiveUserId(input.value);
      await Promise.allSettled([loadCategories(), loadTransactions(), loadStats(getFormData($("#statsFilterForm")))]);
    }, "Пользователь выбран");
  });
}

function setupStats() {
  $("#statsFilterForm").addEventListener("submit", (event) => {
    event.preventDefault();
    runSafely(() => loadStats(getFormData(event.currentTarget)), "Статистика обновлена");
  });
}

function setupTransactions() {
  $("#transactionForm").addEventListener("submit", (event) => {
    event.preventDefault();

    runSafely(async () => {
      const form = event.currentTarget;
      const data = getFormData(form);
      const id = data.id;
      const type = assertTransactionType(data.type);
      const sum = assertPositiveInteger(data.sum, "Сумма");
      const categoryId = assertPositiveInteger(data.category_id, "Категория");
      assertDateTimeLocal(data.date, "Дата операции");

      if (id) {
        await transactionsApi.update(id, onlyFilled({
          categoryId,
          comments: data.comments?.trim() || "",
          sum,
          typeTransaction: type,
        }));
      } else {
        await transactionsApi.create({
          category_id: categoryId,
          comments: data.comments?.trim() || "",
          date: toApiDateTime(data.date),
          sum,
          type,
          user_id: state.userId,
        });
      }

      resetTransactionForm();
      await Promise.allSettled([
        loadTransactions(),
        loadStats(getFormData($("#statsFilterForm"))),
      ]);
    }, "Операция сохранена");
  });

  $("#transactionFilterForm").addEventListener("submit", (event) => {
    event.preventDefault();
    runSafely(() => loadTransactions(getFormData(event.currentTarget)), "Фильтры применены");
  });

  $("#resetTransactionFiltersBtn").addEventListener("click", () => {
    window.setTimeout(() => {
      runSafely(() => loadTransactions({ limit: DEFAULT_PAGE_SIZE, offset: 0 }));
    });
  });

  $("#reloadTransactionsBtn").addEventListener("click", () => {
    runSafely(() => loadTransactions(), "Операции обновлены");
  });

  $("#resetTransactionFormBtn").addEventListener("click", () => {
    window.setTimeout(resetTransactionForm);
  });

  $("#transactionsTbody").addEventListener("click", async (event) => {
    const button = event.target.closest("button[data-action]");
    if (!button) return;

    const row = button.closest("tr");
    const id = row?.dataset.id;
    const action = button.dataset.action;
    const transaction = state.transactions.find((item) => Number(item.id) === Number(id));

    if (!id || !transaction) return;

    if (action === "edit") {
      fillTransactionForm(transaction);
      window.scrollTo({ top: $("#transactionForm").offsetTop - 40, behavior: "smooth" });
      return;
    }

    if (action === "delete") {
      const confirmed = await confirmAction({
        title: "Удалить операцию?",
        text: `Операция #${id} будет удалена и не должна участвовать в статистике.`,
      });

      if (!confirmed) return;

      runSafely(async () => {
        await transactionsApi.delete(id);
        await Promise.allSettled([
          loadTransactions(),
          loadStats(getFormData($("#statsFilterForm"))),
        ]);
      }, "Операция удалена");
    }
  });
}

function setupCategories() {
  $("#categoryForm").addEventListener("submit", (event) => {
    event.preventDefault();

    runSafely(async () => {
      const data = getFormData(event.currentTarget);
      const id = data.id;
      const categoryName = validateCategoryName(data.category_name);

      if (id) {
        await categoriesApi.update(id, { title: categoryName });
      } else {
        await categoriesApi.create({
          category_name: categoryName,
          user_id: state.userId,
        });
      }

      resetCategoryForm();
      await loadCategories();
    }, "Категория сохранена");
  });

  $("#reloadCategoriesBtn").addEventListener("click", () => {
    runSafely(() => loadCategories(), "Категории обновлены");
  });

  $("#resetCategoryFormBtn").addEventListener("click", () => {
    window.setTimeout(resetCategoryForm);
  });

  $("#categoriesTbody").addEventListener("click", async (event) => {
    const button = event.target.closest("button[data-action]");
    if (!button) return;

    const id = button.closest("tr")?.dataset.id;
    const action = button.dataset.action;
    const category = state.categories.find((item) => Number(getCategoryId(item)) === Number(id));

    if (!id || !category) return;

    if (action === "edit") {
      fillCategoryForm(category);
      window.scrollTo({ top: $("#categoryForm").offsetTop - 40, behavior: "smooth" });
      return;
    }

    if (action === "delete") {
      const confirmed = await confirmAction({
        title: "Удалить категорию?",
        text: `Категория «${getCategoryTitle(category) ?? id}» будет удалена. Убедитесь, что backend корректно обрабатывает связанные операции.`,
      });

      if (!confirmed) return;

      runSafely(async () => {
        await categoriesApi.delete(id);
        await loadCategories();
      }, "Категория удалена");
    }
  });
}

function setupLimits() {
  $("#limitForm").addEventListener("submit", (event) => {
    event.preventDefault();

    runSafely(async () => {
      const data = getFormData(event.currentTarget);
      const id = data.id;
      const amountLimit = assertPositiveInteger(data.amount_limit, "Сумма лимита");
      assertDateTimeLocal(data.duration, "Период лимита");

      const payload = {
        amount_limit: amountLimit,
        duration: toApiDateTime(data.duration),
      };

      if (id) {
        await limitsApi.update(id, payload);
      } else {
        await limitsApi.create(payload);
      }

      resetLimitForm();
      await loadLimits();
    }, "Лимит сохранён");
  });

  $("#reloadLimitsBtn").addEventListener("click", () => {
    runSafely(() => loadLimits(), "Лимиты обновлены");
  });

  $("#resetLimitFormBtn").addEventListener("click", () => {
    window.setTimeout(resetLimitForm);
  });

  $("#limitsTbody").addEventListener("click", async (event) => {
    const button = event.target.closest("button[data-action]");
    if (!button) return;

    const id = button.closest("tr")?.dataset.id;
    const action = button.dataset.action;
    const limit = state.limits.find((item) => Number(item.id) === Number(id));

    if (!id || !limit) return;

    if (action === "edit") {
      fillLimitForm(limit);
      window.scrollTo({ top: $("#limitForm").offsetTop - 40, behavior: "smooth" });
      return;
    }

    if (action === "delete") {
      const confirmed = await confirmAction({
        title: "Удалить лимит?",
        text: `Лимит #${id} будет удалён.`,
      });

      if (!confirmed) return;

      runSafely(async () => {
        await limitsApi.delete(id);
        await loadLimits();
      }, "Лимит удалён");
    }
  });
}

function setupUsers() {
  $("#userForm").addEventListener("submit", (event) => {
    event.preventDefault();

    runSafely(async () => {
      const payload = validateUserPayload(getFormData(event.currentTarget));
      await usersApi.create(payload);
      event.currentTarget.reset();
      await loadUsers();
    }, "Пользователь создан");
  });

  $("#reloadUsersBtn").addEventListener("click", () => {
    runSafely(() => loadUsers(), "Пользователи обновлены");
  });

  $("#usersTbody").addEventListener("click", async (event) => {
    const button = event.target.closest("button[data-action]");
    if (!button) return;

    const id = button.closest("tr")?.dataset.id;
    const action = button.dataset.action;

    if (!id) return;

    if (action === "select") {
      runSafely(async () => {
        setActiveUserId(id);
        $("#activeUserId").value = id;
        await Promise.allSettled([loadCategories(), loadTransactions(), loadStats(getFormData($("#statsFilterForm")))]);
      }, `Выбран пользователь #${id}`);
      return;
    }

    if (action === "delete") {
      const confirmed = await confirmAction({
        title: "Удалить пользователя?",
        text: `Пользователь #${id} будет удалён.`,
      });

      if (!confirmed) return;

      runSafely(async () => {
        await usersApi.delete(id);
        await loadUsers();
      }, "Пользователь удалён");
    }
  });
}

function setDefaultDates() {
  $("#statsFrom").value = "";
  $("#statsTo").value = "";
  $("#filterFrom").value = "";
  $("#filterTo").value = "";
  $("#filterLimit").value = DEFAULT_PAGE_SIZE;
  $("#filterOffset").value = 0;

  const transactionDate = $("#transactionDate");
  if (transactionDate) {
    transactionDate.value = nowToLocalDateTimeInputValue();
  }
}

function setupApp() {
  setupTabs();
  setupUserPanel();
  setupStats();
  setupTransactions();
  setupCategories();
  setupLimits();
  setupUsers();
  setDefaultDates();

  runSafely(reloadMainData, "Данные загружены");
}

document.addEventListener("DOMContentLoaded", setupApp);
