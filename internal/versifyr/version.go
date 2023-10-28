package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.11"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "sample value"

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_11"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2023-10-28 21:52:53"









