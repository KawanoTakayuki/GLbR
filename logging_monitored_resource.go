package glbr

import (
	"os"
)

// AppEngineResource appengine log resource
// https://cloud.google.com/appengine/docs/standard/go111/runtime?authuser=0&hl=ja#environment_variables
func AppEngineResource() Option {
	return MonitoredResource("gae_app", map[string]string{
		"project_id": os.Getenv("GOOGLE_CLOUD_PROJECT"),
		"module_id":  os.Getenv("GAE_SERVICE"),
		"version_id": os.Getenv("GAE_VERSION"),
	})
}

// CloudFunctionsResource cloudfunctions log resource
// https://cloud.google.com/functions/docs/env-var?hl=ja#environment_variables_set_automatically
func CloudFunctionsResource() Option {
	return MonitoredResource("cloud_function", map[string]string{
		"project_id":    os.Getenv("GCP_PROJECT"),
		"function_name": os.Getenv("FUNCTION_NAME"),
		"region":        os.Getenv("FUNCTION_REGION"),
	})
}
