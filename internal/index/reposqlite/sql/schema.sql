CREATE TABLE IF NOT EXISTS symbols (
    id        INTEGER PRIMARY KEY,
    name      TEXT NOT NULL DEFAULT '',
    kind      TEXT NOT NULL DEFAULT '',
    package   TEXT NOT NULL DEFAULT '',
    file      TEXT NOT NULL DEFAULT '',
    signature TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS relations (
    id        INTEGER PRIMARY KEY,
    from_id   INTEGER NOT NULL REFERENCES symbols(id),
    to_id     INTEGER NOT NULL REFERENCES symbols(id),
    edge_kind TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_relations_from ON relations(from_id, edge_kind);
CREATE INDEX IF NOT EXISTS idx_relations_to ON relations(to_id, edge_kind);

CREATE TABLE IF NOT EXISTS methods (
    id        INTEGER PRIMARY KEY,
    parent_id INTEGER NOT NULL REFERENCES symbols(id),
    name      TEXT NOT NULL DEFAULT '',
    signature TEXT NOT NULL DEFAULT '',
    kind      TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS files (
    id      INTEGER PRIMARY KEY,
    path    TEXT NOT NULL UNIQUE,
    lang    TEXT NOT NULL DEFAULT '',
    hash    TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT ''
);
