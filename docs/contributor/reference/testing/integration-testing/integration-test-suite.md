---
myst:
  html_meta:
    description: "Build integration test suites for Juju with organized setup, teardown logic, and reusable test includes for bash-based testing."
---

(integration-test-suite)=
# Integration test suite
> Source: https://github.com/juju/juju/tree/main/tests/suites

An **integration test suite** is a collection of integration tests. Each suite has a distinct set-up and tear-down
logic. Integration test suites are often composed of {ref}`test includes <test-include>`.