#!/bin/bash
# Exit if anything fails
set -e
# Ensure we're working in the current directory
cd "$(dirname "$0")"

# Use Python to extract the command line options
# (I tried pure Bash, it was a pain)
goose_cmdline=$(python3 - <<'EOF'
import toml
import sys

config = toml.load("./config.toml")
print(
  "user='{user}' password='{password}' dbname='{name}' host='{host}' port='{port}' sslmode='{ssl}'"
    .format(**{k: v.replace("'", "\\'") for k, v in config['db'].items()})
)
EOF
)

# Run DB migrations
cd ./db/migrations
echo \$ goose postgres "${goose_cmdline}" up
goose postgres "${goose_cmdline}" up

# Now run Hummingbard itself
cd ../..
echo \$ ./bin/hummingbard
./bin/hummingbard
