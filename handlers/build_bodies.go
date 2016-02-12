package handlers

const (
	defaultBuildImg = "quay.io/arschles/gbs-env:0.0.1"
)

type buildReq struct {
	BuildImage   string `json:"build_image"`
	CGOEnabled   bool   `json:"cgo_enabled"`
	CrossCompile bool   `json:"cross_compile"`
}

func newBuildReq() *buildReq {
	return &buildReq{}
}

func (b *buildReq) buildImage() string {
	ret := defaultBuildImg
	if b.BuildImage != "" {
		ret = b.BuildImage
	}
	return ret
}

func (b *buildReq) envs() []string {
	var env []string
	if !b.CGOEnabled {
		env = append(env, "CGO_ENABLED=0")
	} else {
		env = append(env, "CGO_ENABLED=1")
	}
	if b.CrossCompile {
		env = append(env, "CROSS_COMPILE=1")
	} else {
		env = append(env, "CROSS_COMPILE=0")
	}
	return env
}

type buildResp struct {
	StatusURL string `json:"status_url"`
}
