package formatters

import (
	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/exporters"
	"github.com/justinbarrick/fluxcloud/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	fluxevent "github.com/weaveworks/flux/event"
	"testing"
)

func TestNewDefaultFormatter(t *testing.T) {
	config := config.NewFakeConfig()
	config.Set("github_url", "https://github.com/")

	formatter, err := NewDefaultFormatter(config)
	assert.Nil(t, err)
	assert.Equal(t, "https://github.com/", formatter.vcsLink)
	assert.Equal(t, config, formatter.config)
}

func TestNewDefaultFormatterNoGithubLink(t *testing.T) {
	config := config.NewFakeConfig()

	_, err := NewDefaultFormatter(config)
	assert.NotNil(t, err)
}

func TestDefaultFormatterImplementsFormatter(t *testing.T) {
	_ = Formatter(&DefaultFormatter{})
}

func TestDefaultFormatterFormatSyncEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}

	event := test_utils.NewFluxSyncEvent()

	msg := d.FormatEvent(event, &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com/commit/810c2e6f22ac5ab7c831fe0dd697fe32997b098f", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventSync, msg.Type)
	assert.Equal(t, `Event: Sync: 810c2e6, default:deployment/test
Commits: 

* <https://github.com/commit/810c2e6f22ac5ab7c831fe0dd697fe32997b098f|810c2e6>: change test image

Resources updated:

* default:deployment/test`, msg.Body)
	assert.Equal(t, event, msg.Event)
}

func TestDefaultFormatterFormatCommitEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}
	msg := d.FormatEvent(test_utils.NewFluxCommitEvent(), &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com/commit/d644e1a05db6881abf0cdb78299917b95f442036", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventCommit, msg.Type)
	assert.Equal(t, `Event: Commit: d644e1a, default:deployment/test

Resources updated:

* default:deployment/test`, msg.Body)
}

func TestDefaultFormatterFormatAutoReleaseEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}
	msg := d.FormatEvent(test_utils.NewFluxAutoReleaseEvent(), &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventAutoRelease, msg.Type)
	assert.Equal(t, `Event: Automated release of justinbarrick/nginx:test3

Resources updated:

* default:deployment/test`, msg.Body)
}

func TestDefaultFormatterFormatUpdatePolicyEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}
	msg := d.FormatEvent(test_utils.NewFluxUpdatePolicyEvent(), &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com/commit/d644e1a05db6881abf0cdb78299917b95f442036", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventSync, msg.Type)
	assert.Equal(t, `Event: Sync: d644e1a, default:deployment/test
Commits: 

* <https://github.com/commit/d644e1a05db6881abf0cdb78299917b95f442036|d644e1a>: Automated: default:deployment/test

Resources updated:

* default:deployment/test`, msg.Body)
}

func TestDefaultFormatterFormatSyncErrorEvent(t *testing.T) {
	d := DefaultFormatter{
		vcsLink: "https://github.com",
	}

	event := test_utils.NewFluxSyncErrorEvent()

	msg := d.FormatEvent(event, &exporters.FakeExporter{})
	assert.Equal(t, "https://github.com/commit/4997efcd4ac6255604d0d44eeb7085c5b0eb9d48", msg.TitleLink)
	assert.Equal(t, "Applied flux changes to cluster", msg.Title)
	assert.Equal(t, fluxevent.EventSync, msg.Type)
	assert.Equal(t, `Event: Sync: 4997efc, default:persistentvolumeclaim/test
Commits: 

* <https://github.com/commit/4997efcd4ac6255604d0d44eeb7085c5b0eb9d48|4997efc>: create invalid resource

Resources updated:

* default:persistentvolumeclaim/test

Errors:

Resource default:persistentvolumeclaim/test, file: manifests/test.yaml:

> running kubectl: The PersistentVolumeClaim "test" is invalid: spec: Forbidden: field is immutable after creation

Resource default:persistentvolumeclaim/lol, file: manifests/lol.yaml:

> running kubectl: The PersistentVolumeClaim "lol" is invalid: spec: Forbidden: field is immutable after creation
`, msg.Body)
	assert.Equal(t, event, msg.Event)
}
