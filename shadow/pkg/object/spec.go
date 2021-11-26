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

package object

type (
	// MySQL configures the MySQL configuration of the shadow service.
	MySQL struct {
		Uris     string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string `yaml:"password" json:"password" jsonschema:"required"`
	}

	// Kafka configures the Kafka configuration of the shadow service.
	Kafka struct {
		Uris string `yaml:"uris" json:"uris" jsonschema:"required"`
	}

	// Redis configures the Redis configuration of the shadow service.
	Redis struct {
		Uris     string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string `yaml:"password" json:"password" jsonschema:"required"`
	}

	// RabbitMQ configures the RabbitMQ configuration of the shadow service.
	RabbitMQ struct {
		Uris     string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string `yaml:"password" json:"password" jsonschema:"required"`
	}

	// ElasticSearch configures the ElasticSearch configuration of the shadow service.
	ElasticSearch struct {
		Uris     string `yaml:"uris" json:"uris" jsonschema:"required"`
		UserName string `yaml:"userName" json:"userName" jsonschema:"required"`
		Password string `yaml:"password" json:"password" jsonschema:"required"`
	}

	// ShadowService is used to create a shadow service of an existing service.
	ShadowService struct {
		Name          string         `yaml:"name" jsonschema:"required"`
		ServiceName   string         `yaml:"serviceName" jsonschema:"required"`
		Namespace     string         `yaml:"namespace" jsonschema:"required"`
		MySQL         *MySQL         `yaml:"mysql" jsonschema:"omitempty"`
		Kafka         *Kafka         `yaml:"kafka" jsonschema:"omitempty"`
		Redis         *Redis         `yaml:"redis" jsonschema:"omitempty"`
		RabbitMQ      *RabbitMQ      `yaml:"rabbitMq" jsonschema:"omitempty"`
		ElasticSearch *ElasticSearch `yaml:"elasticSearch" jsonschema:"omitempty"`
	}
)
