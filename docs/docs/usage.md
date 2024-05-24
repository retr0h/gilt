---
sidebar_position: 4
---

# Usage

## CLI

### Init Configuration

Initializes config file in the shell's current working directory:

```bash
gilt init
```

### Overlay Repository

Overlay a remote repository into the destination provided.

```bash
gilt overlay
```

### Debug

Display the git commands being executed.

```bash
gilt --debug overlay
```

### Skipping post-commands

Overlay files only, but run no other commands.

```bash
gilt overlay --no-commands
```

## Package

### Overlay Repository

See example client in `examples/go-client/`.

```go
func main() {
	debug := true
	logger := getLogger(debug)

	c := config.Repositories{
		Debug:   debug,
		GiltDir: "~/.gilt",
		Repositories: []config.Repository{
			{
				Git:     "https://github.com/retr0h/ansible-etcd.git",
				Version: "77a95b7",
				DstDir:  "../tmp/retr0h.ansible-etcd",
			},
		},
	}

	var r repositoriesManager = repositories.New(c, logger)
	r.Overlay()
}
```
