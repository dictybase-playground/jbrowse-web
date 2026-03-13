## Requirements

[go]
[bun]
[just 1.4+](https://github.com/casey/just)
[samtools](https://www.htslib.org/)
[tabix](https://www.htslib.org/doc/tabix.html)

## Development

### Local Development

#### Quick Start

1. Prepare local configuration file for the Jbrowse app.
```
just add-assembly canonical_core.fa
```
This generates an index file (`canonical_core.fa.fai`) and creates `config.local.json` at the root of the project.

2. Start the development server
```
just dev
```
This copies `config.local.json` into `./jbrowse2` as `config.json`. Then it serves the application and data files.


<!-- If there is no remote storage for jbrowse assets available: -->
<!---->
<!-- 1. Create `./test-data` directory. -->
<!-- 2. Place jbrowse data assets in `./test-data` -->
<!-- 3. Generate index files of jbrowse data assets (.fa.fai, gff3.tb.gz) -->
<!-- 4.  -->
<!-- 5. run `just dev` -->
<!---->
