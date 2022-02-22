package api

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/googleapis/gax-go/v2/apierror"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	pubsubApiLogger = ctrl.Log.WithName("google_pubsub_api.go")
)

type GCPUBSUBClientCalls interface {
	GetTopic(topicName string) (*pubsub.Topic, error)
	CreateSubscription(topic *pubsub.Topic, filter string, subscriptionName string) error
	GetSubscription(subscriptionName string) *pubsub.Subscription
	UpdateSubscription(topic *pubsub.Topic, filter string, subscriptionName string) error
	DeleteSubscription(subscription *pubsub.Subscription) error
}

type GCPPUBClient struct {
	PubsubClient *pubsub.Client
}

func (c *GCPPUBClient) DeleteSubscription(subscription *pubsub.Subscription) error {
	ctx := context.Background()
	return subscription.Delete(ctx)
}

// Get a topic from google pubsub
//
// topicName the name of the topic to reference
//
// returns a topic or an error if failure to get topic or it does not exist
func (c *GCPPUBClient) GetTopic(topicName string) (*pubsub.Topic, error) {
	ctx := context.Background()
	topic := c.PubsubClient.Topic(topicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return topic, nil
	} else {
		return nil, fmt.Errorf("topic %s does not exist", topicName)
	}
}

// Create a pub sub subscription to a topic if it does not already exist
//
// topic the name of the topic to subscribe to
//
// filter  string to filter by
//
// subscriptionName name of the subscription to create
//
// returns an error if there was a failure creating the subscription
func (c *GCPPUBClient) CreateSubscription(topic *pubsub.Topic, filter string, subscriptionName string) error {

	//filter := fmt.Sprintf("hasPrefix(attributes.objectId, \"%s\")", filterByPrefix)
	pubsubConfig := pubsub.SubscriptionConfig{
		Topic:               topic,
		AckDeadline:         30 * time.Second,
		RetainAckedMessages: false,
	}
	if len(filter) > 0 {
		pubsubConfig.Filter = filter
	}
	ctx := context.Background()
	_, err := c.PubsubClient.CreateSubscription(ctx, subscriptionName, pubsubConfig)
	if err != nil {
		gcpErr, ok := err.(*apierror.APIError)
		if ok {
			// Subscription already exists
			if gcpErr.GRPCStatus().Code() == 6 {
				pubsubApiLogger.Info(fmt.Sprintf("Subscription %s for topic %s already exists.", subscriptionName, topic.ID()))
				return nil
			} else {
				return err
			}
		}
		return err
	}
	pubsubApiLogger.Info(fmt.Sprintf("Created subscription %s for topic %s.", subscriptionName, topic.ID()))
	return nil
}

// Get a pubsub subscription
//
// subscriptionName name of the subscription to receive
func (c *GCPPUBClient) GetSubscription(subscriptionName string) *pubsub.Subscription {
	subscription := c.PubsubClient.Subscription(subscriptionName)
	return subscription
}

// Get a pubsub subscription
//
// subscriptionName name of the subscription to receive
func (c *GCPPUBClient) UpdateSubscription(topic *pubsub.Topic, filter string, subscription *pubsub.Subscription) error {
	// TODO: implement logic for update not needed atm

	// ctx := context.Background()
	// pubsubConfig := pubsub.SubscriptionConfigToUpdate{
	// 	Topic:               topic,
	// 	AckDeadline:         30 * time.Second,
	// 	RetainAckedMessages: false,
	// }
	// if len(filter) > 0 {
	// 	pubsubConfig.Filter = filter
	// }
	// _, err := subscription.Update(ctx, pubsubConfig)
	// return subscription
	return nil
}
