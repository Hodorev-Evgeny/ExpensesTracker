import { categoriesApi, getCookie, limitsApi, sessionApi, statsApi, transactionsApi, usersApi } from "./api.js?v=11";
import { DEFAULT_PAGE_SIZE, STORAGE_KEYS } from "./config.js?v=11";
import {
  getCategoryByLimitId,
  getCategoryId,
  getCategoryLimitId,
  getCategoryTitle,
  getCategoryUserId,
  getLimitId,
  getUserId,
  setCategoryLimitId,
  setCurrentUser,
  setTheme,
  state,
  toggleTheme,
} from "./state.js?v=11";
import {
  assertDateTimeLocal,
  assertPositiveInteger,
  assertTransactionType,
  dateInputToUnixSeconds,
  nowToLocalDateTimeInputValue,
  toApiDateTime,
  validateCategoryName,
  validateUserPatchPayload,
} from "./validators.js?v=11";
import {
  $,
  $$,
  fillCategoryForm,
  fillLimitForm,
  fillTransactionForm,
  fillUserForm,
  renderCategories,
  renderCategorySelects,
  renderLimits,
  renderStats,
  renderTransactions,
  renderUsers,
  resetCategoryForm,
  resetLimitForm,
  resetTransactionForm,
  resetUserForm,
  setLoading,
  showToast,
  updateThemeButton,
} from "./render.js?v=11";

const CATEGORY_LIMIT_PATCH_FIELD = "limit_id";

function getFormData(form) {
  return Object.fromEntries(new FormData(form).entries());
}

function onlyFilled(object) {
  return Object.fromEntries(
    Object.entries(object).filter(([, value]) => value !== undefined && value !== null && value !== ""),
  );
}

function getEntityId(entity) {
  return entity?.id ?? entity?.ID ?? entity?.Id;
}

function getTransactionCategoryId(transaction) {
  return transaction?.category_id ?? transaction?.categoryId ?? transaction?.CategoryID ?? transaction?.CategoryId;
}

function getTransactionType(transaction) {
  return transaction?.type ?? transaction?.Type;
}

function getTransactionDate(transaction) {
  return transaction?.date ?? transaction?.Date;
}

function getTransactionSum(transaction) {
  return transaction?.sum ?? transaction?.Sum;
}

function getTransactionUserId(transaction) {
  return transaction?.user_id ?? transaction?.userId ?? transaction?.UserID ?? transaction?.UserId;
}

function getTransactionComments(transaction) {
  return transaction?.comments ?? transaction?.Comments ?? "";
}

function getSessionUserId(cookieData) {
  return cookieData?.userID
    ?? cookieData?.userId
    ?? cookieData?.user_id
    ?? cookieData?.UserID
    ?? cookieData?.id
    ?? cookieData?.ID;
}

function getStoredCurrentUserId() {
  const raw = localStorage.getItem(STORAGE_KEYS.currentUserId);
  const numericId = Number(raw);
  return Number.isFinite(numericId) && numericId > 0 ? numericId : null;
}

function setStoredCurrentUserId(userId) {
  const numericId = Number(userId);
  if (Number.isFinite(numericId) && numericId > 0) {
    localStorage.setItem(STORAGE_KEYS.currentUserId, String(numericId));
  }
}

async function resolveCurrentUserId() {
  const sessionID = getCookie("sessionID");

  if (sessionID) {
    try {
      const session = await sessionApi.get(sessionID);
      const userId = getSessionUserId(session);
      if (userId) {
        setStoredCurrentUserId(userId);
        return userId;
      }
    } catch (error) {
      console.warn("Не удалось получить пользователя по sessionID, пробую локальный fallback", error);
    }
  }

  const storedUserId = getStoredCurrentUserId();
  if (storedUserId) {
    return storedUserId;
  }

  throw new Error("Не удалось определить пользователя. Войдите в аккаунт заново.");
}

function findCategoryById(categoryId) {
  return state.categories.find((category) => Number(getCategoryId(category)) === Number(categoryId));
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

function requireSessionUserId() {
  if (!state.userId) {
    throw new Error("Пользователь не определён. Войдите в аккаунт заново.");
  }
  return state.userId;
}

async function loadCurrentUser() {
  const userId = await resolveCurrentUserId();
  const user = await usersApi.get(userId);
  setCurrentUser(userId, user);
  setStoredCurrentUserId(userId);
  renderUsers();
  fillUserForm(user);
}

function filterCurrentUserCategories(categories) {
  const userId = state.userId;
  return categories.filter((category) => {
    const categoryUserId = getCategoryUserId(category);
    return !categoryUserId || Number(categoryUserId) === Number(userId);
  });
}

function getLimitCategoryId(limit) {
  return limit?.category_id ?? limit?.categoryId ?? limit?.CategoryID ?? limit?.CategoryId;
}

function getLimitUserId(limit) {
  return limit?.user_id ?? limit?.userId ?? limit?.UserID ?? limit?.UserId;
}

function filterCurrentUserLimits(limits) {
  const currentCategoryIds = new Set(
    state.categories
      .map((category) => Number(getCategoryId(category)))
      .filter((id) => Number.isFinite(id) && id > 0),
  );

  const currentLimitIds = new Set(
    state.categories
      .map((category) => Number(getCategoryLimitId(category)))
      .filter((id) => Number.isFinite(id) && id > 0),
  );

  return (Array.isArray(limits) ? limits : []).filter((limit) => {
    const limitId = Number(getLimitId(limit));
    const limitUserId = Number(getLimitUserId(limit));
    const limitCategoryId = Number(getLimitCategoryId(limit));

    if (Number.isFinite(limitUserId) && limitUserId > 0) {
      return Number(limitUserId) === Number(state.userId);
    }

    if (Number.isFinite(limitCategoryId) && limitCategoryId > 0) {
      return currentCategoryIds.has(limitCategoryId);
    }

    return Number.isFinite(limitId) && currentLimitIds.has(limitId);
  });
}

function prepareTransactionQuery(filterData = {}) {
  return onlyFilled({
    user_id: requireSessionUserId(),
    category_id: filterData.category_id,
    sum: filterData.sum,
    from: dateInputToUnixSeconds(filterData.from),
    to: dateInputToUnixSeconds(filterData.to, true),
    limit: filterData.limit || DEFAULT_PAGE_SIZE,
    offset: filterData.offset || 0,
  });
}

function prepareStatsQuery(filterData = {}) {
  return onlyFilled({
    from: dateInputToUnixSeconds(filterData.from),
    to: dateInputToUnixSeconds(filterData.to, true),
    category_id: filterData.category_id,
  });
}

function applyLocalTransactionFilters(transactions, filterData = {}) {
  let filtered = Array.isArray(transactions) ? [...transactions] : [];

  filtered = filtered.filter((item) => {
    const itemUserId = getTransactionUserId(item);
    return !itemUserId || Number(itemUserId) === Number(state.userId);
  });

  if (filterData.type) {
    filtered = filtered.filter((item) => getTransactionType(item) === filterData.type);
  }

  if (filterData.category_id) {
    filtered = filtered.filter((item) => Number(getTransactionCategoryId(item)) === Number(filterData.category_id));
  }

  if (filterData.sum) {
    filtered = filtered.filter((item) => Number(getTransactionSum(item)) >= Number(filterData.sum));
  }

  if (filterData.from) {
    const from = dateInputToUnixSeconds(filterData.from) * 1000;
    filtered = filtered.filter((item) => {
      const date = new Date(getTransactionDate(item)).getTime();
      return !Number.isNaN(date) && date >= from;
    });
  }

  if (filterData.to) {
    const to = dateInputToUnixSeconds(filterData.to, true) * 1000;
    filtered = filtered.filter((item) => {
      const date = new Date(getTransactionDate(item)).getTime();
      return !Number.isNaN(date) && date <= to;
    });
  }

  if (filterData.comment) {
    const needle = filterData.comment.trim().toLowerCase();
    filtered = filtered.filter((item) => String(getTransactionComments(item)).toLowerCase().includes(needle));
  }

  filtered.sort((a, b) => {
    const dateA = new Date(getTransactionDate(a) ?? 0).getTime();
    const dateB = new Date(getTransactionDate(b) ?? 0).getTime();
    return dateB - dateA;
  });

  const limit = Number(filterData.limit || DEFAULT_PAGE_SIZE);
  const offset = Number(filterData.offset || 0);

  return filtered.slice(offset, offset + limit);
}

async function loadCategories() {
  const categories = await categoriesApi.list();
  state.categories = filterCurrentUserCategories(categories);
  renderCategories();
  renderCategorySelects();
}

async function loadTransactions(filterData = state.filters.transactions) {
  state.filters.transactions = {
    limit: DEFAULT_PAGE_SIZE,
    offset: 0,
    ...filterData,
  };

  const query = prepareTransactionQuery(state.filters.transactions);
  const transactions = await transactionsApi.list(query);

  state.transactions = Array.isArray(transactions)
    ? transactions
    : transactions
      ? [transactions]
      : [];

  const visibleTransactions = applyLocalTransactionFilters(state.transactions, state.filters.transactions);

  renderTransactions(visibleTransactions);
  renderLimits();
}

async function loadStats(filterData = {}) {
  state.filters.stats = { ...filterData };
  state.stats = await statsApi.get(requireSessionUserId(), prepareStatsQuery(filterData));
  renderStats();
}

async function loadLimits() {
  const limits = await limitsApi.list({ limit: 100, offset: 0 });
  state.limits = filterCurrentUserLimits(limits);
  renderLimits();
}

async function attachLimitToCategory(categoryId, limitId) {
  const category = findCategoryById(categoryId);

  if (!category) {
    throw new Error("Категория для лимита не найдена");
  }

  setCategoryLimitId(categoryId, limitId);

  try {
    await categoriesApi.update(categoryId, {
      title: getCategoryTitle(category),
      [CATEGORY_LIMIT_PATCH_FIELD]: limitId,
    });
  } catch (error) {
    console.warn("Не удалось сохранить limit_id в категории на backend", error);
    showToast("Лимит сохранён. Привязка к категории сохранена локально, но backend её не принял.", "warning");
  }
}

async function detachLimitFromCategory(limitId) {
  const category = getCategoryByLimitId(limitId);
  if (!category) return;

  const categoryId = getCategoryId(category);
  setCategoryLimitId(categoryId, null);

  await categoriesApi.update(categoryId, {
    title: getCategoryTitle(category),
    [CATEGORY_LIMIT_PATCH_FIELD]: null,
  });
}

async function reloadMainData() {
  const app = $(".app");
  setLoading(app, true);

  try {
    await loadCurrentUser();
    await loadCategories();

    const results = await Promise.allSettled([
      loadTransactions(),
      loadLimits(),
      loadStats(getFormData($("#statsFilterForm"))),
    ]);

    results.forEach((result) => {
      if (result.status === "rejected") {
        console.error(result.reason);
        showToast(result.reason?.message || "Часть данных не загрузилась", "warning");
      }
    });

    renderLimits();
  } catch (error) {
    console.error(error);
    renderUsers();
    showToast(error.message || "Не удалось загрузить профиль", "error");
    window.setTimeout(() => {
      window.location.href = "/login";
    }, 900);
  } finally {
    setLoading(app, false);
  }
}

function setupTabs() {
  $$(".tab").forEach((tab) => {
    tab.addEventListener("click", () => activateTab(tab.dataset.tab));
  });
}

function setupTheme() {
  setTheme(state.theme);
  updateThemeButton();

  $("#themeToggleBtn")?.addEventListener("click", () => {
    toggleTheme();
    updateThemeButton();
  });
}

function setupStats() {
  $("#statsFilterForm").addEventListener("submit", (event) => {
    event.preventDefault();
    runSafely(() => loadStats(getFormData(event.currentTarget)), "Статистика обновлена");
  });

  $("#resetStatsFiltersBtn")?.addEventListener("click", () => {
    const form = $("#statsFilterForm");
    form.reset();
    $("#statsFrom").value = "";
    $("#statsTo").value = "";
    $("#statsCategoryId").value = "";
    runSafely(() => loadStats({}), "Фильтры статистики сброшены");
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
          user_id: requireSessionUserId(),
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
    const transaction = state.transactions.find((item) => Number(getEntityId(item)) === Number(id));

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
        const category = findCategoryById(id);
        await categoriesApi.update(id, onlyFilled({
          title: categoryName,
          [CATEGORY_LIMIT_PATCH_FIELD]: getCategoryLimitId(category),
        }));
      } else {
        await categoriesApi.create({
          category_name: categoryName,
          user_id: requireSessionUserId(),
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
        setCategoryLimitId(id, null);
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
      const categoryId = assertPositiveInteger(data.category_id, "Категория лимита");
      const amountLimit = assertPositiveInteger(data.amount_limit, "Сумма лимита");
      assertDateTimeLocal(data.duration, "Период лимита");

      const payload = {
        amount_limit: amountLimit,
        duration: toApiDateTime(data.duration),
      };

      let savedLimit;

      if (id) {
        savedLimit = await limitsApi.update(id, payload);
      } else {
        savedLimit = await limitsApi.create(payload);
      }

      const savedLimitId = getLimitId(savedLimit) ?? Number(id);

      if (!savedLimitId) {
        throw new Error("Backend не вернул id лимита");
      }

      await attachLimitToCategory(categoryId, savedLimitId);

      resetLimitForm();

      await Promise.allSettled([
        loadCategories(),
        loadLimits(),
        loadTransactions(),
        loadStats(getFormData($("#statsFilterForm"))),
      ]);
    }, "Лимит сохранён и привязан к категории");
  });

  $("#reloadLimitsBtn").addEventListener("click", () => {
    runSafely(async () => {
      await Promise.allSettled([loadCategories(), loadLimits(), loadTransactions()]);
    }, "Лимиты обновлены");
  });

  $("#resetLimitFormBtn").addEventListener("click", () => {
    window.setTimeout(resetLimitForm);
  });

  $("#limitsTbody").addEventListener("click", async (event) => {
    const button = event.target.closest("button[data-action]");
    if (!button) return;

    const id = button.closest("tr")?.dataset.id;
    const action = button.dataset.action;
    const limit = state.limits.find((item) => Number(getLimitId(item)) === Number(id));

    if (!id || !limit) return;

    if (action === "edit") {
      fillLimitForm(limit);
      window.scrollTo({ top: $("#limitForm").offsetTop - 40, behavior: "smooth" });
      return;
    }

    if (action === "delete") {
      const linkedCategory = getCategoryByLimitId(id);

      const confirmed = await confirmAction({
        title: "Удалить лимит?",
        text: linkedCategory
          ? `Лимит #${id}, привязанный к категории «${getCategoryTitle(linkedCategory)}», будет удалён.`
          : `Лимит #${id} будет удалён.`,
      });

      if (!confirmed) return;

      runSafely(async () => {
        await limitsApi.delete(id);
        await detachLimitFromCategory(id).catch((error) => console.warn("Не удалось отвязать лимит от категории", error));
        await Promise.allSettled([
          loadCategories(),
          loadLimits(),
          loadTransactions(),
          loadStats(getFormData($("#statsFilterForm"))),
        ]);
      }, "Лимит удалён");
    }
  });
}

function setupUsers() {
  $("#userForm").addEventListener("submit", (event) => {
    event.preventDefault();

    runSafely(async () => {
      const payload = validateUserPatchPayload(getFormData(event.currentTarget));
      const updatedUser = await usersApi.update(requireSessionUserId(), payload);
      state.currentUser = updatedUser;
      renderUsers();
      fillUserForm(updatedUser);
    }, "Профиль изменён");
  });

  $("#resetUserFormBtn").addEventListener("click", () => {
    window.setTimeout(resetUserForm);
  });

  $("#editActiveUserBtn")?.addEventListener("click", () => {
    if (!state.currentUser) {
      showToast("Профиль ещё не загружен", "warning");
      return;
    }

    fillUserForm(state.currentUser);
    window.scrollTo({ top: $("#userForm").offsetTop - 40, behavior: "smooth" });
  });

  $("#reloadUsersBtn").addEventListener("click", () => {
    runSafely(loadCurrentUser, "Профиль обновлён");
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
  setupTheme();
  setupTabs();
  setupStats();
  setupTransactions();
  setupCategories();
  setupLimits();
  setupUsers();
  setDefaultDates();

  runSafely(reloadMainData, "Данные загружены");
}

document.addEventListener("DOMContentLoaded", setupApp);
