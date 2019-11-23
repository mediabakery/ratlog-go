package ratlog

import (
	"io"
	"sort"
	"strings"
)

// Ratlog implements Tag and Log function.
// It will be initialized with an io.Writer
type Ratlog struct {
	writer io.Writer
	tags   []string
}

// Fields holds a map of fields
// https://github.com/ratlog/ratlog-spec#fields
type Fields map[string]string

// Tag returns a new Ratlog instance with appended tags
// https://github.com/ratlog/ratlog-spec#tags
func (log Ratlog) Tag(tags []string) Ratlog {
	return Ratlog{writer: log.writer, tags: append(log.tags, tags...)}
}

func getEscapedString(value string, toEscape []string) string {

	for _, escape := range toEscape {
		value = strings.ReplaceAll(value, escape, "\\"+escape)
	}
	return value
}

func getTagsString(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	escapedTags := make([]string, len(tags))
	for i, v := range tags {
		escapedTags[i] = getEscapedString(v, []string{"]", "|"})
	}

	return "[" + strings.Join(escapedTags, "|") + "] "
}

// GetTags returns an array of tags
func GetTags(args ...string) []string {
	return args
}

// GetFields returns a map of fields.
// Every odd arg is used as a key and the even as corresponding value.
// The field is only used if there is a corresponding value
func GetFields(args ...string) map[string]string {
	var fields = make(map[string]string)
	for index := 0; index < len(args)/2; index++ {
		fields[args[index*2]] = args[index*2+1]
	}
	return fields
}

func getFieldsString(fields Fields) string {
	data := []string{}
	keys := []string{}
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if len(fields[k]) != 0 {
			data = append(data, getEscapedString(k, []string{"|", ":"})+": "+getEscapedString(fields[k], []string{"|", ":"}))
		} else {
			data = append(data, getEscapedString(k, []string{"|", ":"}))
		}
	}
	if len(data) == 0 {
		return ""
	}
	return " | " + strings.Join(data, " | ")
}

func getLogMessage(message string, fields Fields, tags []string) string {
	return strings.ReplaceAll(
		getTagsString(tags)+
			getEscapedString(message, []string{"[", "|"})+
			getFieldsString(fields),
		"\n", "\\n") + "\n"
}

// Props represents the properties of a Log call
type Props struct {
	message string
	fields  Fields
	tags    []string
}

// Log writes a log to the writer
func (log *Ratlog) Log(properties Props) error {
	_, err := log.writer.Write([]byte(getLogMessage(properties.message, properties.fields, properties.tags)))
	return err
}

// New returns a new instance of Ratlog without any tags
func New(writer io.Writer) Ratlog {
	return Ratlog{writer: writer, tags: []string{}}
}

// LogItem represents a log.
type LogItem struct {
	props Props
	log   *Ratlog
}

// Tag adds tags to the LogItem
func (logItem *LogItem) Tag(tags ...string) *LogItem {
	logItem.props.tags = tags
	return logItem
}

// Fields adds fields to the LogItem.
// Use pairs of key, values
func (logItem *LogItem) Fields(fields ...string) *LogItem {
	logItem.props.fields = GetFields(fields...)
	return logItem
}

// Log writes the log to the writer
func (logItem *LogItem) Log() error {
	return logItem.log.Log(logItem.props)
}

// Message Log a message and return a LogItem
// You can use this to chain a log.
// ratlog.Message("Foo").Tag("bla", "baz", "bar").Fields("a":"1")
func (log *Ratlog) Message(message string) *LogItem {
	return &LogItem{log: log, props: Props{message: message}}
}
