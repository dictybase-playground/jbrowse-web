# dictybase/jbrowse-web

Static server that serves a JBrowse 2 app for local development.

## Table of Contents

- [Requirements](#requirements)
- [Project Structure](#project-structure)
- [Start](#start)
- [Adding Data Files](#adding-data-files)
- [Configuration](#configuration)
- [Static Server](#static-server)
- [Recreating the JBrowse Build](#recreating-the-jbrowse-build)

## Requirements

- [bun](https://bun.sh/)
- [just 1.4+](https://github.com/casey/just)
- [samtools](https://www.htslib.org/)
- [tabix](https://www.htslib.org/doc/tabix.html)
- [bgzip](https://www.htslib.org/doc/bgzip.html)

## Project Structure

```
jbrowse-launch/
├── jbrowse2/          # JBrowse 2 web app build
│   └── config.json    # Active JBrowse configuration (do not edit directly)
├── server/            # Static file server (TypeScript/Bun)
├── test_data/         # Sample data files (committed)
├── config.json        # JBrowse config referencing test_data/
└── justfile           # Task runner
```

## Start

The repo includes sample data (`test_data/canonical_core.fa`, `test_data/DDB0166986.gff.sorted.gff3.gz`, `test_data/DDB0166986.gff.sorted.gff3.gz.tbi`) and a `config.json` that references them, so no data setup is needed.

**1. Clone and enter the repo**

```sh
git clone <repo-url>
cd jbrowse-launch
```

**2. Install dependencies**

```sh
bun install
```

**3. Start the dev server**

```sh
just dev
```

Loads `config.json` into `jbrowse2/config.json` and starts the server. Open `http://localhost:3000`.

## Adding Data Files

Place new files in `test_data/`, then use the commands below to index them and register them in `config.json`. Run `just dev` afterward to reload the config.

#### FASTA

```sh
just add-assembly <filename.fa>
```

Runs `samtools faidx` to produce a `.fai` index, then adds the assembly to `config.json`. JBrowse uses the FASTA as the reference sequence (genome) that tracks are displayed against.

#### GFF3

```sh
just add-track <filename.gff3>
```

Sorts the GFF3 by chromosome and position, compresses it with `bgzip` (producing `<filename>.sorted.gff3.gz`), indexes it with `tabix`, and adds the track to `config.json`.

## Configuration

| File | Purpose |
|---|---|
| `config.json` | JBrowse config; committed, references `test_data/` |
| `jbrowse2/config.json` | Active config loaded by the app; overwritten by `just dev` |

`just load-config` copies `config.json` → `jbrowse2/config.json`.

## Static Server

The server is a minimal Bun HTTP server (`server/main.ts`) with automatic MIME type detection.

```sh
# Direct usage
bun serve <root> <assets> [-p <port>]

# Via just
just serve             # defaults: jbrowse2/, test_data/, port 3000
just serve 8080   # custom port
```

## Recreating the JBrowse Build

Downloads and unpacks a fresh JBrowse 2 web app into `jbrowse2/` and applies `config.json`.

```sh
just create
```
