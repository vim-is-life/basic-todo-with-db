CREATE TABLE `TodoList` (
  `todo_id`     INTEGER  NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
  `name`        TEXT     NOT NULL UNIQUE,
  `kind`        INTEGER  NOT NULL DEFAULT 0,
  `state`       INTEGER  NOT NULL DEFAULT 0
);

INSERT INTO TodoList
VALUES (1, 'finish todo app', 0, 0);

INSERT INTO TodoList (name, kind, state)
VALUES ('learn golang', 0, 0);

INSERT INTO TodoList (name, kind, state)
VALUES ('learn rust', 0, 0);

INSERT INTO TodoList (name, kind, state)
VALUES ('learn haskell?', 0, 0);

INSERT INTO TodoList (name, kind, state)
VALUES ('look thru doom source code', 0, 0);
