export function assertPositiveInteger(value, fieldName) {
  const numberValue = Number(value);

  if (!Number.isInteger(numberValue) || numberValue <= 0) {
    throw new Error(`${fieldName}: значение должно быть положительным целым числом`);
  }

  return numberValue;
}

export function assertRequiredString(value, fieldName) {
  const stringValue = String(value ?? "").trim();

  if (!stringValue) {
    throw new Error(`${fieldName}: поле обязательно для заполнения`);
  }

  return stringValue;
}

export function assertTransactionType(value) {
  const type = assertRequiredString(value, "Тип операции");

  if (!["Expenditure", "Income"].includes(type)) {
    throw new Error("Тип операции должен быть Expenditure или Income");
  }

  return type;
}

export function assertDateTimeLocal(value, fieldName) {
  const rawValue = assertRequiredString(value, fieldName);
  const date = new Date(rawValue);

  if (Number.isNaN(date.getTime())) {
    throw new Error(`${fieldName}: укажите корректную дату`);
  }

  return date;
}

export function toApiDateTime(value) {
  return assertDateTimeLocal(value, "Дата").toISOString();
}

export function dateInputToRFC3339(value, endOfDay = false) {
  if (!value) return undefined;

  const [year, month, day] = String(value).split("-").map(Number);

  if (!year || !month || !day) {
    throw new Error("Фильтр даты заполнен некорректно");
  }

  const date = new Date(
    year,
    month - 1,
    day,
    endOfDay ? 23 : 0,
    endOfDay ? 59 : 0,
    endOfDay ? 59 : 0,
    endOfDay ? 999 : 0,
  );

  if (Number.isNaN(date.getTime())) {
    throw new Error("Фильтр даты заполнен некорректно");
  }

  return date.toISOString();
}

export function dateInputToUnixSeconds(value, endOfDay = false) {
  const rfc3339 = dateInputToRFC3339(value, endOfDay);
  if (!rfc3339) return undefined;
  return Math.floor(new Date(rfc3339).getTime() / 1000);
}

export function toLocalDateTimeInputValue(value) {
  if (!value) return "";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";

  const timezoneOffset = date.getTimezoneOffset() * 60000;
  return new Date(date.getTime() - timezoneOffset).toISOString().slice(0, 16);
}

export function nowToLocalDateTimeInputValue() {
  const now = new Date();
  const timezoneOffset = now.getTimezoneOffset() * 60000;
  return new Date(now.getTime() - timezoneOffset).toISOString().slice(0, 16);
}

export function validateCategoryName(value) {
  const categoryName = assertRequiredString(value, "Название категории");

  if (categoryName.length < 3 || categoryName.length > 20) {
    throw new Error("Название категории должно быть от 3 до 20 символов");
  }

  return categoryName;
}

export function validateUserPayload(payload) {
  return {
    full_name: assertRequiredString(payload.full_name, "Имя"),
    email: assertRequiredString(payload.email, "Email"),
    phone: assertRequiredString(payload.phone, "Телефон"),
    password: assertRequiredString(payload.password, "Пароль"),
  };
}
