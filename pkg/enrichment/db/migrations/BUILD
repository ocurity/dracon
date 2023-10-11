filegroup(
    name = "migrations",
    srcs = glob(["*.sql"]),
    visibility = ["PUBLIC"],
)

gentest(
    name = "migrations_test",
    data = [":migrations"],
    exit_on_error = True,
    no_test_output = True,
    test_cmd = """
set -Eeuo pipefail

# ensure all migrations end in .down.sql or .up.sql
for src in $DATA; do
    if ! echo "$src" | grep "\\.down\\.sql$\\|\\.up\\.sql$" > /dev/null; then
        echo "'$src' should end in '.down.sql' or '.up.sql'"
        exit 1
    fi
done
    """,
)
