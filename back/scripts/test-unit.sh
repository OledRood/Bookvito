#!/usr/bin/env sh
set -eu

TOTAL_MIN_COVERAGE="${TOTAL_MIN_COVERAGE:-65}"
USECASE_MIN_COVERAGE="${USECASE_MIN_COVERAGE:-80}"

if go tool covdata >/dev/null 2>&1; then
  go test ./... -coverpkg=./internal/... -coverprofile=coverage.out -count=1

  awk \
    -v total_min="${TOTAL_MIN_COVERAGE}" \
    -v usecase_min="${USECASE_MIN_COVERAGE}" \
    '
      NR == 1 { next }
      {
        split($1, location, ":")
        file = location[1]
        statements = $2 + 0
        hits = $3 + 0

        total_statements += statements
        if (hits > 0) {
          covered_statements += statements
        }

        if (file ~ /internal\/usecase\//) {
          usecase_total += statements
          if (hits > 0) {
            usecase_covered += statements
          }
        }
      }
      END {
        if (total_statements == 0) {
          print "coverage profile is empty"
          exit 1
        }

        total_pct = (covered_statements * 100) / total_statements
        usecase_pct = usecase_total == 0 ? 0 : (usecase_covered * 100) / usecase_total

        printf("coverage(total)=%.2f%%, coverage(usecase)=%.2f%%\n", total_pct, usecase_pct)

        if (total_pct + 0.0001 < total_min) {
          printf("total coverage gate failed: %.2f%% < %.2f%%\n", total_pct, total_min)
          exit 2
        }
        if (usecase_pct + 0.0001 < usecase_min) {
          printf("usecase coverage gate failed: %.2f%% < %.2f%%\n", usecase_pct, usecase_min)
          exit 3
        }
      }
    ' coverage.out
else
  echo "go tool covdata is unavailable; running unit tests without coverage gate"
  go test ./... -count=1
fi
