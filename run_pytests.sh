#!/usr/bin/env bash
# If Python files are present, run pytest and exit with its status; otherwise, skip
set -e
if command -v rg >/dev/null 2>&1; then
  has_py=$(rg --files --iglob '*.py' | head -n1 || true)
else
  has_py=$(find . -type f -name '*.py' | head -n1 || true)
fi
if [ -n "$has_py" ]; then
  echo "Python files detected, running pytest..."
  if ! command -v pytest >/dev/null 2>&1; then
    echo "Error: pytest not found. Please install pytest to run tests." >&2
    exit 1
  fi
  pytest "$@"
  exit $?
else
  echo "No Python files found, skipping pytest."
  exit 0
fi