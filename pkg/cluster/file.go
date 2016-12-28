package cluster

type FileSpec struct {
	Name               string `json:"name,omitempty" yaml:"name,omitempty"`
	Encoding           string `json:"encoding,omitempty" yaml:"encoding,omitempty" valid:"^(base64|b64|gz|gzip|gz\\+base64|gzip\\+base64|gz\\+b64|gzip\\+b64)$"`
	Content            string `json:"content,omitempty" yaml:"content,omitempty"`
	Template           string `json:"template,omitempty" yaml:"template,omitempty"`
	Owner              string `json:"owner,omitempty" yaml:"owner,omitempty"`
	Path               string `json:"path,omitempty" yaml:"path,omitempty"`
	RawFilePermissions string `json:"permissions,omitempty" yaml:"permissions,omitempty" valid:"^0?[0-7]{3,4}$"`
}
