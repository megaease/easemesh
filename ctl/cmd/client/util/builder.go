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

package util

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var FileExtensions = []string{".json", ".yaml", ".yml"}

type (
	VisitorBuilder struct {
		visitors          []Visitor
		decoder           Decoder
		httpGetAttempts   int
		errs              []error
		singleItemImplied bool
		commandOptions    *CommandOptions
		filenameOptions   *FilenameOptions
		stdinInUse        bool
	}

	CommandOptions struct {
		// Kind is required.
		Kind string
		// Name is allowed to be empty.
		Name string
	}

	FilenameOptions struct {
		Filenames []string
		Recursive bool
	}
)

func NewVisitorBuilder() *VisitorBuilder {
	return &VisitorBuilder{httpGetAttempts: 3, decoder: newDefaultDecoder()}
}

func (b *VisitorBuilder) HTTPAttemptCount(httpGetAttempts int) *VisitorBuilder {
	b.httpGetAttempts = httpGetAttempts
	return b
}

func (b *VisitorBuilder) FilenameParam(filenameOptions *FilenameOptions) *VisitorBuilder {
	b.filenameOptions = filenameOptions
	return b
}

func (b *VisitorBuilder) CommandParam(commandOptions *CommandOptions) *VisitorBuilder {
	b.commandOptions = commandOptions
	return b
}

func (b *VisitorBuilder) Command() *VisitorBuilder {
	if b.commandOptions == nil {
		return b
	}

	b.visitors = append(b.visitors, NewCommandVisitor(
		b.commandOptions.Kind,
		b.commandOptions.Name,
	))

	return b
}

func (b *VisitorBuilder) Do() ([]Visitor, error) {
	b.Command()
	b.File()

	if len(b.errs) != 0 {
		return nil, fmt.Errorf("%+v", b.errs)
	}

	return b.visitors, nil
}

func (b *VisitorBuilder) File() *VisitorBuilder {
	if b.filenameOptions == nil {
		return b
	}

	recursive := b.filenameOptions.Recursive
	paths := b.filenameOptions.Filenames
	for _, s := range paths {
		switch {
		case s == "-":
			b.Stdin()
		case strings.Index(s, "http://") == 0 || strings.Index(s, "https://") == 0:
			url, err := url.Parse(s)
			if err != nil {
				b.errs = append(b.errs, fmt.Errorf("the URL passed to filename %q is not valid: %v", s, err))
				continue
			}
			b.URL(b.httpGetAttempts, url)
		default:
			if !recursive {
				b.singleItemImplied = true
			}
			b.Path(recursive, s)
		}
	}

	return b
}

func (b *VisitorBuilder) URL(httpAttemptCount int, urls ...*url.URL) *VisitorBuilder {
	for _, u := range urls {
		b.visitors = append(b.visitors, &URLVisitor{
			URL:              u,
			StreamVisitor:    NewStreamVisitor(nil, b.decoder, u.String()),
			HTTPAttemptCount: httpAttemptCount,
		})
	}
	return b
}

func (b *VisitorBuilder) Stdin() *VisitorBuilder {
	if b.stdinInUse {
		b.errs = append(b.errs, errors.Errorf("Stdin already in used"))
	}
	b.stdinInUse = true
	b.visitors = append(b.visitors, FileVisitorForSTDIN(b.decoder))
	return b
}

func (b *VisitorBuilder) Path(recursive bool, paths ...string) *VisitorBuilder {
	for _, p := range paths {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			b.errs = append(b.errs, fmt.Errorf("the path %q does not exist", p))
			continue
		}
		if err != nil {
			b.errs = append(b.errs, fmt.Errorf("the path %q cannot be accessed: %v", p, err))
			continue
		}

		visitors, err := ExpandPathsToFileVisitors(b.decoder, p, recursive, FileExtensions)
		if err != nil {
			b.errs = append(b.errs, fmt.Errorf("error reading %q: %v", p, err))
		}

		b.visitors = append(b.visitors, visitors...)
	}
	if len(b.visitors) == 0 {
		b.errs = append(b.errs, fmt.Errorf("error reading %v: recognized file extensions are %v", paths, FileExtensions))
	}
	return b
}
