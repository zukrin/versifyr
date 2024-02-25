package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.12"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "sample value"

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_12"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2024-02-25 12:02:59"












