package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.20"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "The stars are fire, and the void between them is ice"

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_20"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2026-05-12 15:27:36"
