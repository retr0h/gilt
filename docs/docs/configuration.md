---
sidebar_position: 3
---

# Configuration

Gilt uses [Viper][] to load configuation through multiple methods.

## Config File

Create the giltfile (`Giltfile.yaml`).

Clone the specified `url`@`version` to the configurable path `--gilt-dir`.
Extract the repo the `dstDir` when `dstDir` is provided. Otherwise, copy files
and/or directories to the desired destinations.

```yaml
---
giltDir: ~/.gilt/clone
debug: false
parallel: true
repositories:
  - git: https://github.com/retr0h/ansible-etcd.git
    version: 77a95b7
    dstDir: roles/retr0h.ansible-etcd
  - git: https://github.com/retr0h/ansible-etcd.git
    version: 1.1
    dstDir: roles/retr0h.ansible-etcd-tag
  - git: https://github.com/lorin/openstack-ansible-modules.git
    version: 2677cc3
    sources:
      - src: '*_manage'
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

### Configuration Options

#### `debug`

- Type: boolean
- Default: `false`
- Required: no

Enable / disable debug output

#### `parallel`

- Type: boolean
- Default: `true`
- Required: no

Enable / disable fetching clones concurrently. The default is to fetch clones in
parallel, with one fetch per CPU, and a maximum of 8 concurrent processes.
Setting `parallel: false` will cause Gilt to fetch each clone one-at-a-time.

#### `giltDir`

- Type: string
- Default: `~/.gilt/clone`
- Required: no

Specifies the directory to use for storing cached clones for use by Gilt. The
directory will be created if it does not exist.

#### `repositories`

- Type: list
- Default: `[]`
- Required: no

The list of repositories for Gilt to vendor in. They will be processed in the
order they are defined.

##### `repositories[].git`

- Type: string
- Default: None
- Required: yes

The Git URL of the repository to clone. Any URL format supported by Git may be
used.

##### `repositories[].version`

- Type: string
- Default: None
- Required: yes

The Git commit-ish to use as the source. Any valid branch name, tag name, or
commit hash may be used.

##### `repositories[].dstDir`

- Type: string
- Default: None
- Required: no

The local directory to copy files into. All files in the repository will be
copied. Relative paths will be installed into the directory where `gilt` was
invoked. If `dstDir` already exists, it will be destroyed and overwritten; as
such, `.` and `..` are not allowed.

To copy only a subset of files, use the `repositories.sources` option instead.

This option cannot be used with `repositories.sources`.

##### `repositories[].sources`

- Type: list
- Default: `[]`
- Required: no

A list of subtrees and their targets for Gilt to copy. Relative paths will
read/write into the directory where `gilt` was invoked.

This option cannot be used with `repositories.dstDir`.

###### `repositories[].sources[].src`

- Type: string
- Default: None
- Required: yes

The pathname of the source file/directory to copy.

###### `repositories[].sources[].dstDir`

- Type: string
- Default: None
- Required: no

The pathname of the destination directory. If `src` is a file, it will be placed
inside the named directory. If `src` is a directory, its contents will be copied
into the named directory. All parent directories will be created if they do not
exist. If `dstDir` already exists, it will be destroyed and overwritten; as
such, `.` and `..` are not allowed.

This option cannot be used with `repositories[].sources[].dstFile`.

###### `repositories[].sources[].dstFile`

- Type: string
- Default: None
- Required: no

The pathname of the destination file. If `src` is a directory, an error is
thrown. All parent directories will be created if they do not exist, with an
equivalent set of permissions, i.e., a `src` file with mode `0640` will create
all nonexistant intermediate directories with mode `0750`.

This option cannot be used with `repositories[].sources[].dstDir`.

##### `repositories[].commands`

- Type: list
- Default: `[]`
- Required: no

A list of commands to run after overlaying files. These commands are run in the
same directory used to invoke `gilt`. They will be executed in the order they
are defined, and a non-zero exit status will cause Gilt to abort.

###### `repositories[].commands[].cmd`

- Type: string
- Default: None
- Required: yes

The name of the command to run. The current value of `$PATH` will be used to
find it. This does **NOT** invoke a shell, so variable interpolation, output
redirection, etc., is not supported.

###### `repositories[].commands[].args`

- Type: list of strings
- Default: `[]`
- Required: no

Any and all arguments to the given command. This does **NOT** invoke a shell, so
variable interpolation, output redirection, etc. is not supported. Similarly,
arguments are not split on spaces, so each argument must be a separate list
entry.

## Env Vars

The config file can be overriden/defined through env vars.

```bash
GILT_GILTFILE=Giltfile.yaml \
GILT_GILTDIR=~/.gilt/clone \
GILT_DEBUG=false \
GILT_PARALLEL=false \
gilt overlay
```

### `GILT_DEBUG`

- Default: `false`

Enable/disable debug output.

### `GILT_PARALLEL`

- Default: `true`

Enable / disable fetching clones concurrently. The default is to fetch clones in
parallel, with one fetch per CPU, and a maximum of 8 concurrent processes.
Setting `GIT_PARALLEL=false` will cause Gilt to fetch each clone one-at-a-time.

### `GILT_GILTFILE`

- Default: `./Giltfile.yaml`

Configuration file to use.

### `GILT_GILTDIR`

- Default: `~/.gilt/clone`

Specifies the directory to use for storing cached clones for use by Gilt. The
directory will be created if it does not exist.

### `GILT_SKIPCOMMANDS`

- Default: `false`

If set, Gilt will skip running any post-commands when overlaying files. This can
be useful when debugging.

## Command Flags

The config file and/or env vars can be overriden/defined through cli flags.

```bash
gilt \
  --gilt-file=Giltfile.yaml \
  --gilt-dir=~/.gilt/clone \
  --debug \
  --parallel=false \
  overlay
```

### `-d`, `--debug`

Enable debug output.

### `-c`, `--gilt-dir`

Path to Gilt's clone dir. (default `~/.gilt/clone`)

### `-f`, `--gilt-file`

Path to config file. (default `./Giltfile.yaml`)

### `--no-commands`

If set, Gilt will skip running any post-commands when overlaying files. This can
be useful when debugging.

### `-p`, `--parallel`

Enable / disable fetching clones concurrently. The default is to fetch clones in
parallel, with one fetch per CPU, and a maximum of 8 concurrent processes.
Setting `--parallel=false` will cause Gilt to fetch each clone one-at-a-time.

<!-- prettier-ignore-start -->
[Viper]: https://github.com/spf13/viper
<!-- prettier-ignore-end -->
