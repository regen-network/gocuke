package gocuke

import messages "github.com/cucumber/messages/go/v21"

// Scenario is a special step argument type which describes the running scenario
// and that can be used in a step definition or hook method.
type Scenario interface {
	Name() string
	Tags() []string
	URI() string
	AddAttachment(attachment *Attachment)
	private()
}

type scenario struct {
	runner *scenarioRunner
}

// AddAttachment adds an attachment to test reporting, if reporting is enabled
// (check ReportingEnabled()).
func (s scenario) AddAttachment(attachment *Attachment) {
	if s.runner.reporter != nil {
		s.runner.reporter.Report(&messages.Envelope{Attachment: &messages.Attachment{
			Body:              attachment.Body,
			ContentEncoding:   messages.AttachmentContentEncoding(attachment.ContentEncoding),
			FileName:          attachment.FileName,
			MediaType:         attachment.MediaType,
			Source:            nil,
			TestCaseStartedId: s.runner.testCaseStartedId,
			TestStepId:        s.runner.testStepId,
			Url:               attachment.Url,
		}})
	}
}

func (s scenario) ReportingEnabled() bool {
	return s.runner.reporter != nil
}

// Name returns the scenario name.
func (s scenario) Name() string {
	return s.runner.pickle.Name
}

// Tags returns the scenario tags.
func (s scenario) Tags() []string {
	tags := make([]string, len(s.runner.pickle.Tags))
	for i, tag := range s.runner.pickle.Tags {
		tags[i] = tag.Name
	}
	return tags
}

func (s scenario) URI() string {
	return s.runner.pickle.Uri
}

func (s scenario) private() {}

var _ Scenario = scenario{}
