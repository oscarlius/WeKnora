-- Migration 000080: ensure knowledge_tag_relations exists for upgraded databases.
-- Some long-lived deployments were marked past 000063 without the relation table.
CREATE TABLE IF NOT EXISTS knowledge_tag_relations (
    knowledge_id VARCHAR(36) NOT NULL,
    tag_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (knowledge_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_ktr_knowledge ON knowledge_tag_relations (knowledge_id);
CREATE INDEX IF NOT EXISTS idx_ktr_tag ON knowledge_tag_relations (tag_id);
