 //
 //  Copyright (c) 2020-2021 Datastax, Inc.
 //  
 //  Licensed to the Apache Software Foundation (ASF) under one
 //  or more contributor license agreements.  See the NOTICE file
 //  distributed with this work for additional information
 //  regarding copyright ownership.  The ASF licenses this file
 //  to you under the Apache License, Version 2.0 (the
 //  "License"); you may not use this file except in compliance
 //  with the License.  You may obtain a copy of the License at
 //  
 //     http://www.apache.org/licenses/LICENSE-2.0
 //  
 //  Unless required by applicable law or agreed to in writing,
 //  software distributed under the License is distributed on an
 //  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 //  KIND, either express or implied.  See the License for the
 //  specific language governing permissions and limitations
 //  under the License.
 //

package util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apex/log"
)

// ConsumerResult is Pulsar Consumer result for channel communication
type ConsumerResult struct {
	Err             error
	InOrderDelivery bool
	Latency         time.Duration
	Timestamp       time.Time
}

// VerifyMessageByPulsarConsumer instantiates a Pulsar consumer and verifies an expected message
func VerifyMessageByPulsarConsumer(client pulsar.Client, topicName, expectedMessage string, completeChan chan *ConsumerResult) error {
	topicParts := strings.Split(topicName, "/")
	subscriptionName := "partition-sub" + topicParts[len(topicParts)-1]
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       topicName,
		SubscriptionName:            subscriptionName,
		Type:                        pulsar.Exclusive,
		ReceiverQueueSize:           1,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
	})
	if err != nil {
		log.Errorf("failed to created partition topic consumer, error: %v", err)
		return err
	}
	defer consumer.Close()

	receivedCount := 0
	receiveTimeout := 60 * time.Second
	start := time.Now()
	for time.Since(start) < 90*time.Second { // TODO: this should be passed in as timeout parameter
		cCtx, cancel := context.WithTimeout(context.Background(), receiveTimeout)
		defer cancel()

		log.Infof("%s wait to receive on message count %d", topicName, receivedCount)
		receivedCount++
		msg, err := consumer.Receive(cCtx)
		if err != nil {
			completeChan <- &ConsumerResult{
				Err: fmt.Errorf("consumer Receive() error: %v", err),
			}
			break
		}
		consumer.Ack(msg)
		if expectedMessage == string(msg.Payload()) {
			log.Infof("expected message received by %s", topicName)
			completeChan <- &ConsumerResult{
				InOrderDelivery: true,
				Timestamp:       time.Now(),
			}
			return nil
		}
	}
	return nil
}
