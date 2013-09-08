Diffusion
=========

Diffusion is a simple websocket broadcaster written in Go.

Everything sent on /b/{{channel}} is broadcast to every connection on /{{channel}}.

Broadcaster and client accesses can be restricted by a key for each and they will have to connect to /b/{{channel}}?{{broadcaster key}} or /{{channel}}?{{client key}}.
These keys are definable using command line arguments when starting the program.