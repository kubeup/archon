package ucloud

type CommonRequest struct {
	Action    string
	PublicKey string
	ProjectId string
	Signature string
}

type CommonResponse struct {
	Action  string
	RetCode int
}

type Resource struct {
	ResourceType string
	ResourceName string
	ResourceId   string
	Zone         string
}
