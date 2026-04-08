jbrowse_dir := "jbrowse2"
assets_dir := "test_data"
default_port := "3000"
root := justfile_directory()

# Copy shared config into the jbrowse2 output directory
load-config:
  cp ./config.json {{jbrowse_dir}}/config.json

# Load config and run the dev server
serve port=default_port: load-config
  bun run --filter server serve -- {{root}}/{{jbrowse_dir}} {{root}}/{{assets_dir}} -p {{port}}

# Index a local FASTA file and add its assembly to config.json
# Usage: just add-assembly <fasta_file>
add-assembly fasta_file:
  samtools faidx "{{fasta_file}}"
  bun run aa \
    --load inPlace \
    --force \
    --out config.json \
    "{{fasta_file}}"

# Sort a GFF3 file, compress with bgzip, index with tabix, and add it as a track in config.json
# Usage: just add-track <gff3_file>
add-track gff3_file:
  #!/bin/bash
  sorted=$(just sort-gff {{gff3_file}})
  tabix -f -p gff "$sorted"
  bun run at \
    --load inPlace \
    --force \
    --out config.json \
    "$sorted"

# Sort a GFF3 file by chromosome and position, then compress with bgzip
# Outputs the path to the resulting .sorted.gff3.gz file
# Usage: just sort-gff <gff3_file>
sort-gff gff3_file:
  #!/bin/bash
  bun sort-gff \
    {{gff3_file}} \
    | bgzip > {{gff3_file}}.sorted.gff3.gz
  echo {{gff3_file}}.sorted.gff3.gz

# Build all plugins
build-plugins:
  bun run --filter plugins build

# Scaffold a fresh JBrowse2 instance, remove test data, and apply the production config
create:
  bun run create {{jbrowse_dir}} --force
  rm -rf {{jbrowse_dir}}/test_data
  cp ./config.json {{jbrowse_dir}}/config.json
