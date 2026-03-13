# jitzu

Commitizen for [Jujutsu](https://jj-vcs.github.io/jj/) — interactive conventional commits and bookmark management for `jj`.

## Install

```sh
go install github.com/ATTron/jitzu@latest
```

Or build from source:

```sh
git clone https://github.com/ATTron/jitzu.git
cd jitzu
make install
```

This installs both `jitzu` and a `jz` symlink.

## Setup

```sh
jitzu init
```

Creates a `.jitzu.toml` config in your project. To also install the `jj z` alias:

```sh
jitzu init --install-alias
```

This lets you run `jj z` from any jj repo to start the interactive prompt.

## Usage

### Describe (default)

```sh
jitzu              # describe the working copy
jitzu -r @-        # describe a previous revision
jitzu describe     # same as bare jitzu
```

Walks you through type, scope, subject, body, breaking changes, and issue refs — then runs `jj describe -m "..."`.

### Commit

```sh
jitzu commit
```

Same interactive flow, runs `jj commit -m "..."`.

### Bookmark

```sh
jitzu bookmark
```

Intelligently detects your bookmark situation:

- **Current revision has a bookmark** — keep it or create a new one
- **Parent/grandparent has a non-trunk bookmark** — advance it, set it here, or create a new one
- **Branching off main/trunk** — guided creation of a new bookmark

New bookmarks are created through a form that produces structured names like `feat/auth-login` or `fix/PROJ-123-null-pointer`.

### Check

```sh
jitzu check        # validate current revision
jitzu check @-     # validate a specific revision
```

### Changelog

```sh
jitzu changelog
jitzu changelog -r "trunk()..@"
```

Generates grouped markdown from jj history.

## Configuration

jitzu works with zero config. Optionally create `.jitzu.toml` in your project root:

```toml
# Restrict scopes to a predefined list
scopes = ["api", "ui", "core"]

# Require scope on every commit
scope_required = false

# Require body on every commit
body_required = false

# Maximum subject line length (default: 72)
subject_max_len = 72

# Custom commit types (overrides defaults)
[[types]]
name = "feat"
description = "A new feature"

[[types]]
name = "fix"
description = "A bug fix"
```

Config is discovered by searching upward from CWD, then `~/.config/jitzu/config.toml`, then built-in defaults.

## License

MIT
