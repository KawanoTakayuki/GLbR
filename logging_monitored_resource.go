package glbr

import (
	"context"

	"google.golang.org/genproto/googleapis/api/monitoredres"
)

// MonitoredResourceLabel 監視対象の情報を手動でセット（optionより優先）
// https://cloud.google.com/monitoring/api/resources
// Default: resourceType = project, resourceLabel = {"project_id": $PROJECT_ID}
func (s Service) MonitoredResourceLabel(resourceType string, resourceLabel map[string]string) Service {
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
