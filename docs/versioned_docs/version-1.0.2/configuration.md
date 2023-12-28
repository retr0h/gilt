---
sidebar_position: 3
---

# Configuration

Gilt uses [Viper][] to load configuation through multpile methods.

## Config File

Create the giltfile (`Giltfile.yaml`).

Clone the specified `url`@`version` to the configurable path `--gilt-dir`.
Extract the repo the `dstDir` when `dstDir` is provided.  Otherwise, copy files
and/or directories to the desired destinations.

```yaml
---
giltDir: ~/.gilt/clone
debug: false
repositories:
  - git: https://github.com/retr0h/ansible-etcd.git
    sha: 77a95b7
    dstDir: roles/retr0h.ansible-etcd
  - git: https://github.com/retr0h/ansible-etcd.git
    tag: 1.1
    dstDir: roles/retr0h.ansible-etcd-tag
  - git: https://github.com/lorin/openstack-ansible-modules.git
    sha: 2677cc3
    sources:
      - src: "*_manage"
        dstDir: library
      - src: nova_quota
        dstDir: library
      - src: neutron_router
        dstFile: library/neutron_router.py
      - src: tests
        dstDir: tests
    commands:
      - cmd: ansible-playbook
        args:
          - -i,
          - playbook.yml
      - cmd: bash
        args:
          - -c
          - who | grep tty
```

## Env Vars

The config file can be overriden/defined through env vars.

```
GILT_GILTFILE=Giltfile.yaml \
GILT_GILTDIR=~/.gilt/clone \
GILT_DEBUG=false \
go-gilt overlay
```

## Command Flags

The config file and/or env vars can be overriden/defined through cli flags.

```
go-gilt \
  --gilt-file=Giltfile.yaml \
  --gilt-dir=~/.gilt/clone \
  --debug \
  overlay
```

[Viper]: https://github.com/spf13/viper
