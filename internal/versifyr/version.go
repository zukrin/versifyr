package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.1.3"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "The sky above the port was the color of television, tuned to a dead channel."

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_1_3"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2026-05-01 17:05:40"
