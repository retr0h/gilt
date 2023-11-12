![codecov](https://img.shields.io/codecov/c/github/retr0h/go-gilt?token=clAMnFQCEQ&style=flat-square)
[![go report card](https://goreportcard.com/badge/github.com/retr0h/go-gilt?style=flat-square)](https://goreportcard.com/report/github.com/retr0h/go-gilt)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)

# Gilt

<img src="asset/gilt.png" align="left" width=20% height=20%>

Gilt is a tool which aims to make repo management, manageable.  Gilt
clones repositories at a particular version, then overlays the repository to
the provided destination.  An alternate approach to "vendoring".

What makes Gilt interesting, is the ability to overlay particular files and/or
directories from the specified repository to given destinations.  This is quite
helpful for those using Ansible, since libraries, plugins, and playbooks are
often shared, but Ansible's [Galaxy][] has no mechanism to handle this.

<br clear="left"/>

## Port

This project is a port of [Gilt][], it is
not 100% compatible with the python version, and aims to correct some poor decisions
made in the python version of Gilt.

This version of Gilt does not provide built in locking, unlike our python friend.

## Installation

    $  go install github.com/retr0h/go-gilt@latest

## Usage

### Overlay Repository

Create the giltfile (`gilt.yml`).

Clone the specified `url`@`version` to the configurable path `--giltdir`.
Extract the repo the `dstDir` when `dstDir` is provided.  Otherwise, copy files
and/or directories to the desired destinations.

```yaml
---
- git: https://github.com/retr0h/ansible-etcd.git
  version: 77a95b7
  dstDir: roles/retr0h.ansible-etcd

- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: "*_manage"
      dstDir: library
    - src: nova_quota
      dstDir: library
    - src: neutron_router
      dstFile: library/neutron_router.py
    - src: tests
      dstDir: tests
```

Overlay a remote repository into the destination provided.

    $ gilt overlay

Use an alternate config file (default `gilt.yml`).

    $ gilt overlay --filename /path/to/gilt.yml

Optionally, override gilt's cache location (defaults to `~/.gilt/clone`):

    $ gilt --giltdir ~/alternate/directory overlay

### Debug

Display the git commands being executed.

    $ gilt --debug overlay

## Building

    $ task build

## Testing

### Dependencies

Check installed dependencies:

    $ task deps:check

To execute tests:

    $ task test

Auto format code:

    $ task fmt

List helpful targets:

    $ task

## License

The [MIT][] License.

The logo is licensed under the [Creative Commons NoDerivatives 4.0 License][],
and designed by [@nanotron][].

[Galaxy]: https://docs.ansible.com/ansible/latest/reference_appendices/galaxy.html
[Gilt]: http://gilt.readthedocs.io/en/latest/
[MIT]: LICENSE
[Creative Commons NoDerivatives 4.0 License]: https://creativecommons.org/licenses/by-nd/4.0/
[@nanotron]: https://github.com/nanotron
