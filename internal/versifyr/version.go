package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.18"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "The light that burns twice as bright burns half as long."

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_18"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2026-05-01 19:02:43"
