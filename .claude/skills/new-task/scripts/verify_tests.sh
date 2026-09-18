#!/usr/bin/env bash
#
# verify_tests.sh - run the tests of a task folder against one solution file.
#
# Usage:
#   verify_tests.sh <task_dir> <solution.go> [go test flags...]
#
# The script copies the *_test.go files of the task folder and the solution
# into a throwaway module, then runs go test there. The task folder and the
# repo stay as they are, so the solution never lands in git.
#
# The solution file must hold the package clause, the function, and any helper
# type the tests need, such as ListNode.
#
# The exit code is the exit code of go test: 0 on a pass. A correct solution
# must give 0. Every wrong solution must give something else.

set -u

if [ $# -lt 2 ]; then
	sed -n '3,12p' "$0" >&2
	exit 2
fi

task_dir=$1
solution=$2
shift 2

if [ ! -d "$task_dir" ]; then
	echo "verify_tests.sh: no such task folder: $task_dir" >&2
	exit 2
fi
if [ ! -f "$solution" ]; then
	echo "verify_tests.sh: no such solution file: $solution" >&2
	exit 2
fi

tests=$(ls "$task_dir"/*_test.go 2>/dev/null) || true
if [ -z "$tests" ]; then
	echo "verify_tests.sh: no test file in $task_dir" >&2
	exit 2
fi

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

cp "$task_dir"/*_test.go "$work"/ || exit 2
cp "$solution" "$work"/solution.go || exit 2

# Take the go directive of the repo, so the module builds with the same rules.
root=$(git -C "$task_dir" rev-parse --show-toplevel 2>/dev/null)
go_line=$(grep -m1 '^go ' "${root:-.}/go.mod" 2>/dev/null || echo "go 1.21")
printf 'module verify\n\n%s\n' "$go_line" > "$work/go.mod"

cd "$work" || exit 2
go test "$@" ./...
