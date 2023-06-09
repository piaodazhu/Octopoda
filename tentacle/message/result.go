package message

type Result struct {
	// Result code. Reserved
	Rcode int

	// Result message. OK, or Error Reason
	Rmsg string

	// For script or command execution. Script output.
	Output string

	// For version control. Version hash code.
	Version string

	// For version control. Modified flag.
	Modified bool
}
