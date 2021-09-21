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
	"time"

	"github.com/megaease/easegress/pkg/logger"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
)

// Syncer syncs data from Etcd, it uses an Etcd watcher to receive update.
// The syncer keeps a full copy of data, and keeps apply changes onto it when an
// update event is received from the watcher, and then send out the full data copy.
// The syncer also pulls full data from Etcd at a configurable pull interval, this
// is to ensure data consistency, as Etcd watcher may be cancelled if it cannot catch
// up with the key-value store.
type Syncer struct {
	server       *Server
	pullInterval time.Duration
	done         chan struct{}
}

func (server *Server) NewSyncer(pullInterval time.Duration) (*Syncer, error) {
	return &Syncer{
		server:       server,
		pullInterval: pullInterval,
		done:         make(chan struct{}),
	}, nil
}

func (s *Syncer) pull(kind string) ([]object.CustomObject, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), s.server.RequestTimeout)
	defer cancelFunc()

	result, err := s.server.List(ctx, kind)
	if err != nil {
		logger.Errorf("failed to pull data for kind %s: %v", kind, err)
		return nil, err
	}
	return result, nil
}

func (s *Syncer) watch(kind string, send func(data []object.CustomObject)) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), s.server.RequestTimeout)
	defer cancelFunc()

	reader, err := s.server.Watch(ctx, kind)
	if err != nil {
		logger.Errorf("Watch response from MeshServer error: %s. Stop watch.", err.Error())
		return
	}
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Errorf("Watch response from MeshServer error: %s. Stop watch.", err.Error())
			return
		} else {
			var objects []object.CustomObject
			err = json.Unmarshal(line, objects)
			if err != nil {
				logger.Errorf("MeshServer returns invalid json: %s, error: %s. Skipped.", line, err)
				continue
			}
			send(objects)
		}
	}
}

func (s *Syncer) run(kind string, send func(data []object.CustomObject)) {

	ticker := time.NewTicker(s.pullInterval)
	defer ticker.Stop()

	pullCompareSend := func() {
		data, err := s.pull(kind)
		if err != nil {
			logger.Errorf("pull data for kind %s failed: %v", kind, err)
			return
		}
		send(data)
	}

	pullCompareSend()

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			pullCompareSend()
		}
	}
	s.watch(kind, send)
}

// Sync syncs a given EaseMesh kind's value through the returned channel.
func (s *Syncer) Sync(kind string) (<-chan object.CustomObject, error) {
	ch := make(chan object.CustomObject, 10)
	fn := func(data []object.CustomObject) {
		for _, obj := range data {
			ch <- obj
		}
	}

	go func() {
		defer close(ch)
		s.run(kind, fn)
	}()
	return ch, nil
}

// Close closes the syncer.
func (s *Syncer) Close() {
	close(s.done)
}
