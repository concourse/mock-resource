package resource

type Source struct {
	// fetch the resource itself as an image
	MirrorSelf bool `json:"mirror_self"`

	// don't emit an initial version; useful for testing pipeline triggering
	NoInitialVersion bool `json:"no_initial_version"`
}

type Version struct {
	Version string `json:"version"`
}

type GetParams struct {
	// same as configuring MirrorSelf in source, but in params, so that we can
	// test params are respected in places
	MirrorSelfViaParams bool `json:"mirror_self_via_params"`
}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
