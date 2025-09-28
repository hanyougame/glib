package firebasex

type Config struct {
	Type           string `json:"type"`
	ProjectID      string `json:"project_id"`
	PrivateKeyID   string `json:"private_key_id"`
	PrivateKey     string `json:"private_key"`
	ClientEmail    string `json:"client_email"`
	ClientID       string `json:"client_id"`
	AuthURI        string `json:"auth_uri"`
	TokenURI       string `json:"token_uri"`
	AuthProvider   string `json:"auth_provider_x509_cert_url"`
	ClientCert     string `json:"client_x509_cert_url"`
	UniverseDomain string `json:"universe_domain"`
}

type WebConfig struct {
	VAPIDPublicKey string       `json:"vapid_public_key"`
	SdkConfig      WebSdkConfig `json:"sdk_config"`
}

type WebSdkConfig struct {
	ApiKey            string `json:"api_key"`
	AuthDomain        string `json:"auth_domain"`
	ProjectID         string `json:"project_id"`
	StorageBucket     string `json:"storage_bucket"`
	MessagingSenderId string `json:"messaging_sender_id"`
	AppID             string `json:"app_id"`
	MeasurementId     string `json:"measurement_id"`
}
