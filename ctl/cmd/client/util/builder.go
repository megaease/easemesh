package util

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var FileExtensions = []string{".json", ".yaml", ".yml"}

type VisitorBuilder struct {
	visitors          []Visitor
	decoder           Decoder
	httpGetAttempts   int
	errs              []error
	singleItemImplied bool
	filenameOptions   *FilenameOptions
	stdinInUse        bool
}

func NewVisitorBuilder() *VisitorBuilder {
	return &VisitorBuilder{httpGetAttempts: 3, decoder: newDefaultDecoder()}
}

type FilenameOptions struct {
	Filenames []string
	Recursive bool
}

func (b *VisitorBuilder) HttpAttemptCount(httpGetAttempts int) *VisitorBuilder {
	b.httpGetAttempts = httpGetAttempts
	return b
}

func (b *VisitorBuilder) FilenameParam(filenameOptions *FilenameOptions) *VisitorBuilder {
	b.filenameOptions = filenameOptions
	return b
}

func (b *VisitorBuilder) URL(httpAttemptCount int, urls ...*url.URL) *VisitorBuilder {
	for _, u := range urls {
		b.visitors = append(b.visitors, &URLVisitor{
			URL:              u,
			StreamVisitor:    NewStreamVisitor(nil, b.decoder, u.String()),
			HttpAttemptCount: httpAttemptCount,
		})
	}
	return b
}

func (b *VisitorBuilder) Do() ([]Visitor, error) {
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
	return b.visitors, nil
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

func (b *VisitorBuilder) Stream(r io.Reader, name string) *VisitorBuilder {
	b.visitors = append(b.visitors, NewStreamVisitor(r, b.decoder, name))
	return b
}
