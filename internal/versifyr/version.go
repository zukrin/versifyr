package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.1.2"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "<no value>"

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_1_2"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2026-03-31 15:12:50"
