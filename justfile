outdir := "jbrowse2"
remoteUrl := ""

add-assembly fasta_file:
  samtools faidx {{fasta_file}}
  bun run aa {{fasta_file}} --load copy

add-assembly-remote fasta_file:
  bun run aa {{fasta_file}}

create:
  bun run create {{outdir}}
  rm -rf {{outdir}}/test_data
  mv ./config.json {{outdir}}/config.json

# serve:
# fetch-config:
	
