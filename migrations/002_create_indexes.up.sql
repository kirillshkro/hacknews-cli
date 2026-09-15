CREATE INDEX idx_items_name ON items(name);
CREATE INDEX idx_items_type_score ON items(type, score DESC) WHERE deleted = FALSE AND dead = FALSE;
CREATE INDEX idx_items_parent ON items(parent_id);
CREATE INDEX idx_items_created ON items(created_at DESC);
