package glbr

import (
	"context"
	"os"

	"google.golang.org/genproto/googleapis/api/monitoredres"
)

// AppEngineResource appengine log resource
func (s Service) AppEngineResource() Service {
	return s.MonitoredResource("gae_app", map[string]string{
		"project_id": s.projectID,
		"module_id":  os.Getenv("GAE_SERVICE"),
		"version_id": os.Getenv("GAE_VERSION"),
	})
}

// CloudFunctionsResource cloudfunctions log resource
func (s Service) CloudFunctionsResource() Service {
	return s.MonitoredResource("cloud_function", map[string]string{
		"function_name": os.Getenv("FUNCTION_NAME"),
		"project_id":    s.projectID,
		"region":        os.Getenv("FUNCTION_REGION"),
	})
}

// MonitoredResource 監視対象の情報を手動でセット（optionより優先）
// https://cloud.google.com/monitoring/api/resources
// Default: resourceType = project, resourceLabel = {"project_id": $PROJECT_ID}
func (s Service) MonitoredResource(resourceType string, resourceLabel map[string]string) Service {
	// logging.CommonResource()はresourceBlockに反映されない？
	mr := monitoredres.MonitoredResource{
		Type:   resourceType,
		Labels: resourceLabel,
	}
	s.ctx = context.WithValue(s.ctx, &monitoredResourceKey, mr)
	return s
}

func getMonitoredResource(c context.Context) *monitoredres.MonitoredResource {
	if mr, ok := c.Value(&monitoredResourceKey).(monitoredres.MonitoredResource); ok {
		return &mr
	}
	return nil
}

// func newMonitordResource(c context.Context, projectID string, opts ...option.ClientOption) error {
// 	client, err := logadmin.NewClient(c, projectID, opts...)
// 	if err != nil {
// 		return err
// 	}
// 	descriptors := client.ResourceDescriptors(c)
// 	for {
// 		description, err := descriptors.Next()
// 		if err == iterator.Done {
// 			fmt.Printf("done.\n")
// 			break
// 		}
// 		fmt.Printf("Name:[%s]\n", description.GetName())
// 		fmt.Printf("Type:[%s]\n", description.GetType())
// 		fmt.Printf("Label:[%v]\n", description.GetLabels())
// 		fmt.Printf("Launch:[%s]\n", description.GetLaunchStage())
// 		fmt.Printf("Display:[%s]\n", description.GetDisplayName())
// 		fmt.Println()
// 	}
// 	return nil
// }
