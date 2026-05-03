package versifyr

// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.19"

// $versifyr:template=const Sample = "{{ .sample }}"$
const Sample = "Fear is the mind-killer."

// $versifyr:template=const ActualTimestamp = "{{ .version | replace "." "_" }}"$
const ActualTimestamp = "v0_0_19"

// $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
const Compiled = "2026-05-03 16:00:59"
