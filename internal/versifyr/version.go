package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.11a"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "sample value"

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_11a"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2023-10-31 15:12:04"











