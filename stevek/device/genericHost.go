package device

func init() {
	Registry["Generic host"] = Device{
		Name:      "Generic host",
		filterFn:  AllowAll,
		Transform: Identity,
	}
}
