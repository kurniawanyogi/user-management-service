package logger

import (
	"context"
	"fmt"
)

type key string

const (
	// LoggingTagKey is reserved name in the context for auto logging
	LoggingTagKey key = "logging_tag"
)

// AddLoggingTag func to add context logging tag
func AddLoggingTag(ctx context.Context, tagsToAdd ...Tag) context.Context {
	allTags := ctx.Value(LoggingTagKey)
	if len(tagsToAdd) == 0 {
		return ctx
	}
	if contextTags, ok := allTags.(map[string]string); ok && contextTags != nil {
		return context.WithValue(ctx, LoggingTagKey, mergeTags(contextTags, tagsToAdd...))
	}
	return context.WithValue(ctx, LoggingTagKey, mergeTags(nil, tagsToAdd...))
}

func mergeTags(contextTags map[string]string, tagsToAdd ...Tag) map[string]string {
	if contextTags == nil {
		contextTags = make(map[string]string)
	}
	for _, tag := range tagsToAdd {
		contextTags[tag.Key] = fmt.Sprintf("%+v", tag.Value)
	}
	return contextTags
}

// AddRequestID to add X-REQUEST-ID to logging tag
func AddRequestID(ctx context.Context, msgID string) context.Context {
	return AddLoggingTag(ctx, Tag{
		Key:   RequestIDKey,
		Value: msgID,
	})
}

// GetAllLoggingTagInTagStr to get all tag str from logging tag
func GetAllLoggingTagInTagStr(ctx context.Context) []Tag {
	if ctx == nil {
		return nil
	}
	allTags := ctx.Value(LoggingTagKey)
	contextTags, ok := allTags.(map[string]string)
	if !ok || contextTags == nil {
		return nil
	}

	var tags []Tag
	for k, v := range contextTags {
		tags = append(tags, Tag{
			Key:   k,
			Value: v,
		})
	}
	return tags
}

// GetTagValue to get a value from specific key
func GetTagValue(ctx context.Context, tagKey string) string {
	allTags := ctx.Value(LoggingTagKey)
	if contextTags, ok := allTags.(map[string]string); ok && contextTags != nil {
		return contextTags[tagKey]
	}
	return ""
}
