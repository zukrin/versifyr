package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.10"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "<no value>"

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_10"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2023-10-11 17:20:54"

