outdir := "jbrowse2"
asset_url := "http://localhost:3000"
asset_dir := "data"
local_test_data_dir := "test_data"

dev: load-config-local serve

load-config-local:
  cp ./config.local.json {{outdir}}/config.json

[parallel]
serve: serve-app serve-assets

serve-app:
  caddy file-server --root {{outdir}} --listen :3000

serve-assets:
  caddy file-server --root {{local_test_data_dir}} --listen :8080

add-assembly fasta_file:
  samtools faidx {{local_test_data_dir}}/{{fasta_file}}
  bun run aa \
    --force \
    --out config.local.json \
    {{asset_url}}/{{local_test_data_dir}}/{{fasta_file}} 

add-assembly-remote fasta_file:
  bun run aa {{fasta_file}}

load-config: 
  cp ./config.json {{outdir}}/config.json

create:
  bun run create {{outdir}} --force
  rm -rf {{outdir}}/test_data
  cp ./config.json {{outdir}}/config.json
