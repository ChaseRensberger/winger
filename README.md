# winger

A modern take on the [finger protocol](https://en.wikipedia.org/wiki/Finger_(protocol)). Maintain a `~/.plan` file, sign it with Ed25519 keys, sync it to a relay. Anyone can fetch 
your plan and learn about you or see what you're working on.

If you're curious about the motivation for this project or want ', check out the [blog post](wingman.actor/blog/winger)

## How it works

Users keep a `~/.plan` file locally. Every sync cryptographically signs the content and pushes it to a relay server. No accounts, no passwords — just a handle tied to your public key.

Anyone can read anyone's plan, just like the original finger protocol.

## CLI

### Install

```bash
go install winger/cmd/winger@latest
```

### Usage

```bash
winger init            # generate keys, pick a handle, register with relay
winger sync            # sign and push ~/.plan to relay
winger daemon          # watch ~/.plan and auto-sync on changes
winger <handle>        # look up someone's plan
```

### What `winger init` does

1. Prompts for a handle and relay URL
2. Generates an Ed25519 keypair at `~/.winger/identity.key` and `~/.winger/identity.pub`
3. Writes config to `~/.winger/config.toml`
4. Creates a starter `~/.plan` file
5. Registers your handle with the relay

```
~/.winger/
    identity.key    # Ed25519 private key
    identity.pub    # public key
    config.toml     # relay URL, handle
~/.plan             # your plan file
```

## Relay

The relay is the server that stores and serves plan files. It exposes both a JSON API (for the CLI) and an HTMX web UI.

### Run locally

```bash
# build css
bun install
bunx @tailwindcss/cli -i input.css -o static/output.css

# start relay
go run ./cmd/relay
```

The relay listens on `:2323` by default. Set `PORT` and `DB_PATH` environment variables to customize.

### API

| Method | Route | Description |
|--------|-------|-------------|
| `POST` | `/register` | Register a handle with a public key |
| `POST` | `/sync` | Upload a signed plan |
| `GET` | `/{handle}` | View a plan (text/plain from CLI, HTML from browser) |
| `GET` | `/random` | View a random plan |
| `GET` | `/` | Leaderboard — top plans by retrieval count |

### Deploy to Render

Push to GitHub and connect to Render. The included `render.yaml` configures a web service with a persistent disk for the SQLite database.

## Example `.plan`

```
[chase]
Plan:

This is my daily work ...

When I accomplish something, I write a * line that day.

Whenever a bug / missing feature is mentioned during the day and
I don't fix it, I make a note of it. Some things get noted many times
before they get fixed.

Occasionally I go back through the old notes and mark with a +
the things I have since fixed.

= mar 21 ===================================
* winger init flow
* ed25519 key generation
* relay with sqlite storage
* htmx leaderboard
+ plan retrieval counter
```
