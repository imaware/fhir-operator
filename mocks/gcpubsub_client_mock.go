package mocks

import "cloud.google.com/go/pubsub"

type MockGCPPUBSUBClientGood struct {
	PubSubClient *pubsub.Client
}

func (m *MockGCPPUBSUBClientGood) GetTopic(topicName string) (*pubsub.Topic, error) {
	mockedTopic := &pubsub.Topic{}
	return mockedTopic, nil
}

func (m *MockGCPPUBSUBClientGood) CreateSubscription(topic *pubsub.Topic, filterByPrefix string, subscriptionName string) error {
	return nil
}

func (m *MockGCPPUBSUBClientGood) GetSubscription(subscriptionName string) *pubsub.Subscription {
	mockedSubscription := &pubsub.Subscription{}
	return mockedSubscription
}
