// Copyright (c) 2019 - goevent authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gpubsub

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gapsquare/goevent"
	"google.golang.org/api/option"
)

type errorTopicQueue struct {
	client *pubsub.Client
	topic  *pubsub.Topic
	appID  string
}

func newErrorTopicQueue(projectID, appID string, opts ...option.ClientOption) (*errorTopicQueue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}

	//Create or get a topic
	name := appID + "_error_queue"
	topic := client.Topic(name)

	if ok, err := topic.Exists(ctx); err != nil {
		return nil, err
	} else if !ok {
		if topic, err = client.CreateTopic(ctx, name); err != nil {
			return nil, err
		}
	}

	return &errorTopicQueue{
		client: client,
		topic:  topic,
		appID:  appID,
	}, nil
}

func (b *errorTopicQueue) publish(eventtopic goevent.EventTopic, msg *pubsub.Message, parentError error) error {
	publishCtx := context.Background()

	msg.Attributes = map[string]string{
		"event_topic": string(eventtopic),
		"error":       parentError.Error(),
	}

	res := b.topic.Publish(publishCtx, msg)
	if _, err := res.Get(publishCtx); err != nil {
		return errors.New("could not move pubsub message to error queue: " + err.Error())
	}

	return nil
}
