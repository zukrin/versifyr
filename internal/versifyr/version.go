package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.17"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "I've seen things you people wouldn't believe."

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_17"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2026-05-01 17:18:36"
