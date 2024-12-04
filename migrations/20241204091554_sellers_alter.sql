ALTER TABLE sellers_json_history
  ADD column deleted_domains  TEXT NOT NULL,
  ADD column deleted_publishers  TEXT NOT NULL;

