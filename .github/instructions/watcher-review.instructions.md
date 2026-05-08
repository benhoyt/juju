---
applyTo: "**/*_test.go"
---

When using the watcher test harness (`core/watcher/watchertest`):
 * each test must use a single transaction for all mutations that cause the watcher under test to fire. 
   * This is to ensure that those change stream events are in the same term.
* if mutations emit change events for a watcher under test, AssertChangeStreamIdle shall be called before creating the watcher.
    * This is to ensure initial events are tested correctly in your watcher.

