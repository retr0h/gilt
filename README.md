[![codecov](https://img.shields.io/codecov/c/github/retr0h/go-gilt?token=clAMnFQCEQ&style=flat-square)](https://codecov.io/gh/retr0h/go-gilt)
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

This project is a golang port of [Gilt][], and aims to correct poor decisions
made in the python version, primarially around config syntax, and not supporting
braches in favor of reproducibility.

## Installation

### Homebrew Tap

    $ brew install retr0h/tap/go-gilt

### Go Install

    $  go install github.com/retr0h/go-gilt@latest

## Configuration

Gilt uses [Viper][] to load configuation through multpile methods.

### Config File

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
```

### Env Vars

The config file can be overriden/defined through env vars.

    $ GILT_GILTFILE=Giltfile.yaml \
      GILT_GILTDIR=~/.gilt/clone \
      GILT_DEBUG=false \
      go-gilt overlay

### Command Flags

The config file and/or env vars can be overriden/defined through cli flags.

    $ go-gilt \
      --gilt-file=Giltfile.yaml \
      --gilt-dir=~/.gilt/clone \
      --debug \
      overlay

## Usage

### CLI

#### Overlay Repository

Overlay a remote repository into the destination provided.

    $ gilt overlay

#### Debug

Display the git commands being executed.

    $ gilt --debug overlay

### Package

#### Overlay Repository

See example client in `examples/go-client/`.

```golang
func main() {
	debug := true
	logger := getLogger(debug)

	c := config.Repositories{
		Debug:   debug,
		GiltDir: "~/.gilt",
		Repositories: []config.Repository{
			{
				Git:     "https://github.com/retr0h/ansible-etcd.git",
				SHA:     "77a95b7",
				DstDir:  "../tmp/retr0h.ansible-etcd",
			},
		},
	}

	var r repositoriesManager = repositories.New(c, logger)
	r.Overlay()
}
```

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
If you have some other use in mind, contact us.

[Galaxy]: https://docs.ansible.com/ansible/latest/reference_appendices/galaxy.html
[Gilt]: http://gilt.readthedocs.io/en/latest/
[Viper]: https://github.com/spf13/viper
[MIT]: LICENSE
[Creative Commons NoDerivatives 4.0 License]: https://creativecommons.org/licenses/by-nd/4.0/
[@nanotron]: https://github.com/nanotron
