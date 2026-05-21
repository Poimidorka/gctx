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
(didn't find active profile)
p1 p2 p3
```

When a profile is active:

```text
current used profile: p1
p1 p2 p3
```

Apply a profile to the current Git repository:

```bash
gctx p1
```

Save the current repository Git config as a profile:

```bash
gctx p1 --save
gctx p1 -s
```

Remove a profile:

```bash
gctx p1 --remove
gctx p1 -r
```

Use a custom profile directory:

```bash
gctx --config /path/to/profiles
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
