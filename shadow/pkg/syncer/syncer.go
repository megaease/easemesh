/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package syncer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
)

// Syncer syncs data from EaseMesh control plane, it uses a watcher to receive ShadowService object.
// It sends out the full data copy when ShadowService object is received from the watcher.
// The syncer also pulls full data from EaseMesh control plane at a configurable pull interval, this
// is to ensure data consistency, as EaseMesh control plane watcher may be cancelled if it cannot catch
// up with the ShadowService object.

type ShadowServiceSyncer struct {
	server       *Server
	pullInterval time.Duration
	done         chan struct{}
}

func NewSyncer(meshServer string, requestTimeout time.Duration, pullInterval time.Duration) (*ShadowServiceSyncer, error) {
	return &ShadowServiceSyncer{
		server: &Server{
			RequestTimeout: requestTimeout,
			MeshServer:     meshServer,
		},
		pullInterval: pullInterval,
		done:         make(chan struct{}),
	}, nil
}

func (s *ShadowServiceSyncer) pull(kind string) ([]object.ShadowService, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.server.RequestTimeout)
	defer cancelFunc()

	result, err := s.server.List(ctx, kind)
	if err != nil {
		log.Printf("failed to pull data for kind %s: %v", kind, err)
		return nil, err
	}
	log.Printf("Pull ShadowService Success.")
	return result, nil
}

func (s *ShadowServiceSyncer) watch(kind string, send func(data []object.ShadowService)) chan struct{} {
	watchChan := make(chan struct{})
	go func() {
		reader, err := s.server.Watch(kind)
		if err != nil {
			log.Printf("Watch response from MeshServer error: %s. Retry ...", err.Error())
			watchChan <- struct{}{}
			return
		}
		for {
			line, e := reader.ReadBytes('\n')
			if e != nil {
				log.Printf("Watch response from MeshServer error: %s. Retry ...", e.Error())
				watchChan <- struct{}{}
				return
			} else {
				if json.Valid(line) {
					var objects []object.ShadowService
					e = json.Unmarshal(line, &objects)
					if e != nil {
						log.Printf("MeshServer returns invalid json: %s, error: %s. Skipped.", line, e.Error())
						continue
					}
					send(objects)
				}
			}
		}
	}()
	return watchChan
}

func (s *ShadowServiceSyncer) run(kind string, send func(data []object.ShadowService)) {
	watchChan := s.watch(kind, send)
	defer close(watchChan)

	ticker := time.NewTicker(s.pullInterval)
	defer ticker.Stop()

	pullAndSend := func() {
		data, err := s.pull(kind)
		if err != nil {
			log.Printf("pull data for kind %s failed: %v", kind, err)
			return
		}
		send(data)
	}
	pullAndSend()

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			pullAndSend()
		case <-watchChan:
			close(watchChan)
			watchChan = s.watch(kind, send)
		}
	}
}

// Sync syncs a given EaseMesh kind's value through the returned channel.
func (s *ShadowServiceSyncer) Sync(kind string) (<-chan []object.ShadowService, error) {
	ch := make(chan []object.ShadowService)
	fn := func(data []object.ShadowService) {
		if data != nil {
			ch <- data
		}
	}
	go func() {
		defer close(ch)
		s.run(kind, fn)
	}()
	return ch, nil
}

// Close closes the syncer.
func (s *ShadowServiceSyncer) Close() {
	close(s.done)
}
