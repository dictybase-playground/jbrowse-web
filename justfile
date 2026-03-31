outdir := "jbrowse2"
asset_url := "http://localhost:3000"

# Start local development: load local config and run the dev server
dev: load-config-local serve

# Copy local config into the jbrowse2 output directory
load-config-local:
  cp ./config.local.json {{outdir}}/config.json

# Run the dev server
serve:
  bun serve

# Index a local FASTA file and add its assembly to config.local.json
# Usage: just add-assembly <fasta_file>
add-assembly fasta_file:
  samtools faidx {{fasta_file}}
  bun run aa \
    --force \
    --out config.local.json \
    {{asset_url}}/{{fasta_file}}

# Add a remote assembly by URL to the JBrowse config
# Usage: just add-assembly-remote <fasta_url>
add-assembly-remote fasta_file:
  bun run aa {{fasta_file}}

add-track gff3_file:
  #!/bin/bash
  sorted=$(just sort-gff {{gff3_file}})
  tabix $sorted
  bun run at \
    --force \
    --out config.local.json \
    {{asset_url}}/$sorted

sort-gff gff3_file:
  #!/bin/bash
  bun sort-gff \
    {{gff3_file}} \
    | bgzip > {{gff3_file}}.sorted.gff.gz
  echo {{gff3_file}}.sorted.gff.gz

# Copy the production config into the jbrowse2 output directory
load-config:
  cp ./config.json {{outdir}}/config.json

# Scaffold a fresh JBrowse2 instance, remove test data, and apply the production config
create:
  bun run create {{outdir}} --force
  rm -rf {{outdir}}/test_data
  cp ./config.json {{outdir}}/config.json
