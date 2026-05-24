# gctx

`gctx` is a small Git context switcher. It saves named local Git configs and applies them to the current repository.

Profiles are stored in:

```text
~/.config/gctx
```

Each profile is saved as a `.config` file.

## Usage

List profiles:

```bash
gctx
```

When no profile is active:
   
```text
No active profile in current git config.
p1 p2 p3
```

When a profile is active:

```text
Current git profile: "p1".
p1 p2 p3
```

Apply a profile to the current Git repository:

```bash
gctx p1
# ✔ Switched to profile "p1".
```

Save the current repository Git config as a profile:

```bash
gctx p1 --save
gctx p1 -s
# ✔ Saved context to profile: "p1".
```

Remove a profile:

```bash
gctx p1 --remove
gctx p1 -r
# ✔ Removed profile "p1".
```

Use the global Git config instead of the current repository config:

```bash
gctx --global
gctx -g
gctx p1 --global
gctx p1 -g
gctx p1 --save --global
```

Run a command scoped to a profile without changing the current repository or global config:

```bash
gctx p1 -c /bin/bash -i
gctx p1 -c "git config --global user.name"
```

Use a custom profile directory:

```bash
gctx --config /path/to/profiles
```

If a context does not exist, `gctx` prints the saved contexts that are available:

```text
✘ Profile "missing" not found. Available profiles: p1 p2.
```

## Development

Build the binary:

```bash
make build
```

Run tests:

```bash
make test
```

Run the CLI through Go:

```bash
make run
make run ARGS="p1 --save"
make run ARGS="p1"
make run ARGS="p1 --remove"
make run ARGS="p1 --global"
make run ARGS='p1 -c "git config --global user.name"'
```

The built binary is written to:

```text
build/gctx
```

Run it directly:

```bash
./build/gctx
./build/gctx p1
```

Clean build output:

```bash
make clean
```

## Installation

Install with Go:

```bash
go install github.com/Poimidorka/gctx@latest
```

After a tagged release, install with Homebrew from the tap:

```bash
brew tap Poimidorka/tap
brew install --cask gctx
```

## Releases 

Releases are built by GoReleaser when a `v*` tag is pushed:

```bash   
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds Linux, macOS, and Windows archives, publishes checksums, and updates the Homebrew tap configured in `.goreleaser.yaml`.
