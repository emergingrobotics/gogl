package reservations

// Modes are the modifiers governing a reservation operation.
//
// There are no mode selectors here. The single-binary era had Get/Set/Add/Del/Clear
// booleans plus a checkModes validator enforcing that exactly one was chosen; under a
// subcommand tree each is its own command, so the constraint is structural and both the
// booleans and the validator were deleted rather than carried forward as dead weight.
type Modes struct {
	// Domain sets the DNS domain, which reservation writes require.
	Domain string

	// Name, MAC and IP identify a target for Del. Exactly one must be set.
	Name string
	MAC  string
	IP   string

	// Force proceeds past conflicts and skips confirmation prompts.
	Force bool

	// Prune deletes reservations and names present on the router but absent from the
	// imported file.
	Prune bool

	// DryRun reports what would change without changing it.
	DryRun bool
}
