package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/jellydator/ttlcache/v3"
	"google.golang.org/api/option"

	compute "cloud.google.com/go/compute/apiv1"
	nodeCache "github.com/shubhindia/spot-operator/controllers/utils/cache"
)

func DeleteNode(projectID string, zone string, instanceName string) error {

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx, option.WithCredentialsFile("/Users/shubhamgopale/Desktop/Work/Infra/go_scripts/spot-operator.json"))
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
	// add the nodeName and update its deletion status in cache
	hit := nodeCache.Cache().Get(instanceName)
	if hit == nil {
		// add the node into cache
		nodeCache.Cache().Set(instanceName, true, ttlcache.DefaultTTL)
	}

	return nil
}
