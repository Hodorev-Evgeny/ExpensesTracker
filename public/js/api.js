import { API_BASE } from "./config.js?v=11";

function buildUrl(path, query = {}) {
  const url = new URL(`${API_BASE}${path}`, window.location.origin);

  Object.entries(query).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") return;
    url.searchParams.set(key, value);
  });

  return url.toString();
}

function normalizeArray(payload) {
  if (Array.isArray(payload)) return payload;
  if (!payload) return [];

  const candidates = [
    payload.items,
    payload.data,
    payload.results,
    payload.rows,
    payload.list,
    payload.transactions,
    payload.categories,
    payload.limits,
    payload.users,
  ];

  const arrayCandidate = candidates.find(Array.isArray);
  return arrayCandidate ?? [payload];
}

function getErrorMessage(payload, fallback) {
  if (!payload) return fallback;
  return payload.message || payload.massage || payload.error || payload.detail || fallback;
}

function looksLikeHtml(text) {
  return /^\s*<!doctype html/i.test(text) || /^\s*<html/i.test(text);
}

function parseResponseBody(text, response) {
  if (!text) return null;

  try {
    return JSON.parse(text);
  } catch (error) {
    if (response.redirected || looksLikeHtml(text)) {
      throw new Error("Сессия не найдена или истекла. Войдите в аккаунт заново.");
    }

    console.error("Backend вернул невалидный JSON");
    console.error("URL:", response.url);
    console.error("STATUS:", response.status);
    console.error("BODY:", text);

    throw new Error("Backend вернул невалидный JSON. Подробности в Console.");
  }
}

export function getCookie(name) {
  const encodedName = `${encodeURIComponent(name)}=`;
  const item = document.cookie
    .split(";")
    .map((part) => part.trim())
    .find((part) => part.startsWith(encodedName));

  return item ? decodeURIComponent(item.slice(encodedName.length)) : "";
}

async function request(path, options = {}) {
  const { query, body, headers, ...fetchOptions } = options;

  const response = await fetch(buildUrl(path, query), {
    credentials: "include",
    ...fetchOptions,
    headers: {
      ...(body ? { "Content-Type": "application/json" } : {}),
      ...headers,
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (response.redirected && !response.url.includes(API_BASE)) {
    window.location.href = "/login";
    throw new Error("Нужно войти в аккаунт");
  }

  if (response.status === 204) return null;

  const text = await response.text();
  const payload = parseResponseBody(text, response);

  if (!response.ok) {
    if (response.status === 401 || response.status === 403) {
      window.location.href = "/login";
    }
    throw new Error(getErrorMessage(payload, `Ошибка запроса: ${response.status}`));
  }

  return payload;
}

export const usersApi = {
  async list(query = {}) {
    return normalizeArray(await request("/users", { query }));
  },

  async create(data) {
    return request("/users", {
      method: "POST",
      body: data,
    });
  },

  async login(data) {
    return request("/users/login", {
      method: "POST",
      body: data,
    });
  },

  async get(id) {
    return request(`/users/${id}`);
  },

  async update(id, data) {
    return request(`/users/${id}`, {
      method: "PATCH",
      body: data,
    });
  },

  async delete(id) {
    return request(`/users/${id}`, {
      method: "DELETE",
    });
  },
};

export const sessionApi = {
  async get(sessionID) {
    return request(`/session/${encodeURIComponent(sessionID)}`);
  },

  async current() {
    const sessionID = getCookie("sessionID");

    if (!sessionID) {
      throw new Error("Cookie sessionID не найдена. Войдите в аккаунт.");
    }

    return this.get(sessionID);
  },
};

export const categoriesApi = {
  async list() {
    return normalizeArray(await request("/category"));
  },

  async create(data) {
    return request("/category", {
      method: "POST",
      body: data,
    });
  },

  async get(id) {
    return request(`/category/${id}`);
  },

  async update(id, data) {
    return request(`/category/${id}`, {
      method: "PATCH",
      body: data,
    });
  },

  async delete(id) {
    return request(`/category/${id}`, {
      method: "DELETE",
    });
  },
};

export const transactionsApi = {
  async list(query = {}) {
    return normalizeArray(await request("/transactions", { query }));
  },

  async create(data) {
    return request("/transactions", {
      method: "POST",
      body: data,
    });
  },

  async get(id) {
    return request(`/transactions/${id}`);
  },

  async update(id, data) {
    return request(`/transactions/${id}`, {
      method: "PATCH",
      body: data,
    });
  },

  async delete(id) {
    return request(`/transactions/${id}`, {
      method: "DELETE",
    });
  },
};

export const limitsApi = {
  async list(query = {}) {
    return normalizeArray(await request("/limit", { query }));
  },

  async create(data) {
    return request("/limit", {
      method: "POST",
      body: data,
    });
  },

  async get(id) {
    return request(`/limit/${id}`);
  },

  async update(id, data) {
    return request(`/limit/${id}`, {
      method: "PATCH",
      body: data,
    });
  },

  async delete(id) {
    return request(`/limit/${id}`, {
      method: "DELETE",
    });
  },
};

export const statsApi = {
  async get(userId, query = {}) {
    if (!userId) {
      return request("/static", { query });
    }

    return request(`/static/${encodeURIComponent(userId)}`, { query });
  },
};
