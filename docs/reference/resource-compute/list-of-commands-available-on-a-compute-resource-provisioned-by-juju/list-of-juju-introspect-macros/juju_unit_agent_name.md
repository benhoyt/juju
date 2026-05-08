---
myst:
  html_meta:
    description: "Get the unit agent name using juju_unit_agent_name introspect macro to identify running Juju unit and application agents."
---

(juju-unit-agent-name)=
# `juju_unit_agent_name`

The `juju_unit_agent_name` macro returns the name of the unit
agent that is currently running.

## Usage

```python
juju_unit_agent_name
```

Returns a string containing the agent (unit / application) name.

## Example output:

```bash
unit-ubuntu-0
```
