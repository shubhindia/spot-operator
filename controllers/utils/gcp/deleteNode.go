package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/option"

	compute "cloud.google.com/go/compute/apiv1"
)

func DeleteNode(projectID string, zone string, instanceName string) error {

	ctx := context.Background()

	// TODO: Write an initialiser to get all the necessary values
	instancesClient, err := compute.NewInstancesRESTClient(ctx, option.WithCredentialsFile("to-be-picked-from-env"))
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}

	req := &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	_, err = instancesClient.Delete(ctx, req)
	if err != nil {

		return fmt.Errorf("unable to delete instance: %v", err)
	}

	return nil
}
