package resource

const DefaultInitialVersion = "mirror"

type Source struct {
	// fetch the resource itself as an image
	MirrorSelf bool `json:"mirror_self"`

	// initial version that the mirrored resource image should emit from /check
	// (default: 'mirror')
	RawInitialVersion string `json:"mirrored_version"`

	// don't emit an initial version; useful for testing pipeline triggering
	NoInitialVersion bool `json:"no_initial_version"`
}

func (s Source) InitialVersion() string {
	if s.RawInitialVersion == "" {
		return DefaultInitialVersion
	}

	return s.RawInitialVersion
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
