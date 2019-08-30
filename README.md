# powerd

Utility for interacting with the ATEN SMASH-CLP System Management Shell to manage power on my server. It is hosted behind a standard OpenSSH server. The power can be controlled at `/system1/pwrmgtsvc1` with the following verbs: `start`, `stop`, `reset`.

The observed power states are `1=On`, `6=Off`. No intermediate states were observed when doing a quick test of turning the server on or off.

The purpose of this is primarily to allow for remote power management of the server via these on-board services.