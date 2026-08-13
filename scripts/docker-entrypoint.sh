#!/bin/bash
set -e

# ─── Fix ownership of bind-mounted directories ───
# When users bind-mount host directories (e.g. ./skills/preloaded),
# the mount inherits the host UID/GID which may differ from the
# container's appuser. This entrypoint runs as root, fixes ownership,
# then drops privileges to appuser via gosu — the same pattern used
# by official postgres/redis images.

# Directories that may be bind-mounted and need appuser access
MOUNT_DIRS=(
    /app/skills/preloaded
)

for dir in "${MOUNT_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        chown -R appuser:appuser "$dir" 2>/dev/null || true
    fi
done

# /data/files can contain a large production corpus. Recursively chowning it on
# every container start can block readiness for minutes or hours, so only ensure
# the mount root itself is writable by default. Operators can opt into a one-off
# recursive repair with WEKNORA_FIX_DATA_FILES_OWNERSHIP_RECURSIVE=true.
if [ -d /data/files ]; then
    if [ "${WEKNORA_FIX_DATA_FILES_OWNERSHIP_RECURSIVE:-false}" = "true" ]; then
        chown -R appuser:appuser /data/files 2>/dev/null || true
    else
        chown appuser:appuser /data/files 2>/dev/null || true
    fi
fi

# ─── Merge built-in skills into preloaded ───
# Built-in skills are backed up at /app/skills/_builtin during image build.
# After a bind-mount replaces /app/skills/preloaded, copy back any
# missing built-in skills (without overwriting user-provided ones).
BUILTIN_DIR="/app/skills/_builtin"
PRELOADED_DIR="/app/skills/preloaded"

if [ -d "$BUILTIN_DIR" ]; then
    mkdir -p "$PRELOADED_DIR"
    for skill_dir in "$BUILTIN_DIR"/*/; do
        [ -d "$skill_dir" ] || continue
        skill_name="$(basename "$skill_dir")"
        if [ ! -d "$PRELOADED_DIR/$skill_name" ]; then
            cp -r "$skill_dir" "$PRELOADED_DIR/$skill_name"
        fi
    done
    chown -R appuser:appuser "$PRELOADED_DIR"
fi

# ─── Drop privileges and exec the main process ───
exec gosu appuser "$@"
