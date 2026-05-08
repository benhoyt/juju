---
myst:
  html_meta:
    description: "Start a stopped unit agent within a single Juju process using juju_start_unit introspect function on deployed machines."
---

(juju_start_unit)=
# `juju_start_unit`

In 2.9 the machine and unit agents were combined into a single process running on Juju deployed machines. This tools allows you to see the start a stopped unit agent running inside of that single process.  It takes a unit name as input. Example output:

```text
$ juju_start_unit neutron-openvswitch/0
neutron-openvswitch/0: started
```
