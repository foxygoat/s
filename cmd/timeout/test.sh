#!/bin/bash

test_count=0
fail_count=0

main() {
  binary="${1:-./out/timeout}"
  if [[ ! -x "${binary}" ]]; then
    printf 'FAIL: binary '%s' not found\n' "${binary}"
    exit 1
  fi

  # Test: Process exits on signal after timeout (1 second)
  run_test 1 1 "${binary}" 1s sleep 5

  # Test: Process ignores signal, killed after grace after timeout (1 + 1 seconds)
  # Note we use SIGCHLD because that is ignored by default
  run_test 2 1 "${binary}" -s "$(kill -l CHLD)" -g 1s 1s sleep 5

  # Test: Process exits successfully before timeout
  run_test 0 0 "${binary}" 1s true

  # Test: Process exits unsuccessfully before timeout
  run_test 0 2 "${binary}" 1s sh -c 'exit 2'

  # Test: Command with slash that exists
  run_test 0 0 "${binary}" 1s /bin/true

  # Test: Command with slash that does not exist
  run_test 0 1 "${binary}" 1s /bin/non-existent

  # Test: Command without slash that does not exist
  run_test 0 1 "${binary}" 1s non-existent

  # TODO(camh): Test stdin/out/err are preserved

  if (( fail_count > 0 )); then
    printf 'FAIL\t%s : %s/%s tests failed\n' "$0" "${fail_count}" "${test_count}"
    exit 1
  fi

  printf 'ok\t%s\n' "$0"
}

run_test() {
  (( test_count++ ))
  run "$@" || (( fail_count++ ))
}

run() {
  # Set TIMEFORMAT to show only real time (elapsed) and no decimal places.
  # From testing this truncates. Since all our tests run in very close to
  # the second boundary, this gives us a little leeway when running under
  # load, but it does mean we cannot accurately test the timeout period.
  local TIMEFORMAT='%0R'

  expected_duration="$1"; shift
  expected_exit="$1"; shift

  # Discard output of command. Just capture the elapsed time, which needs
  # the braces to the redirection binds to `time`, not the command.
  actual_duration=$({ time "$@" >/dev/null 2>&1; } 2>&1)
  actual_exit=$?

  if (( expected_exit != actual_exit )); then
    printf -- '--- FAIL: %s\n    exit code mismatch: expected %d, actual %d\n' \
      "$*" "${expected_exit}" "${actual_exit}"
    return 1
  fi

  if (( expected_duration != actual_duration )); then
    printf -- '--- FAIL: %s\n    duration mismatch: expected %d, actual %d\n' \
      "$*" "${expected_duration}" "${actual_duration}"
    return 1
  fi

  return 0
}

main "$@"
