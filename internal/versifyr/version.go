package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.14"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "<no value>"

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_14"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2025-02-20 15:46:49"
