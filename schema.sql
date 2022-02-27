CREATE TABLE user_accounts (
  id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  email TEXT UNIQUE
);

CREATE TABLE budgets (
  id UUID PRIMARY KEY,
  user_account_id UUID NOT NULL REFERENCES user_accounts(id),
  name TEXT NOT NULL,
  UNIQUE (user_account_id, name)
);

CREATE TABLE bank_accounts (
  id UUID PRIMARY KEY,
  budget_id UUID NOT NULL REFERENCES budgets(id),
  name TEXT NOT NULL,
  UNIQUE (budget_id, name)
);