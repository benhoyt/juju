---
myst:
  html_meta:
    description: "Retrieve the machine agent name using juju_machine_agent_name introspect macro for identifying running Juju agents."
---

(juju-machine-agent-name)=
# `juju_controller_agent_name`

The `juju_machine_agent_name` macro returns the name of the machine agent
that is currently running.

## Usage

```python
juju_machine_agent_name
```

## Example output

```bash
machine-0
```
