## Requirements

[bun]
[just 1.4+](https://github.com/casey/just)
[samtools](https://www.htslib.org/)
[tabix](https://www.htslib.org/doc/tabix.html)

## Development

### Local Development

If there is no remote storage for jbrowse assets available:

1. Create `./test-data` directory.
2. Place jbrowse data assets in `./test-data`
3. Generate index files of jbrowse data assets (.fai, gff3.tb.gz)
4. 
5. run `just dev`

