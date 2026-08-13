-- Migration 000086: opt-in automatic association of existing knowledge-base tags.
ALTER TABLE knowledge_bases DROP COLUMN IF EXISTS auto_tag_config;
