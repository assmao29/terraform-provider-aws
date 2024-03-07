// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package cloudformation

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
)

// []*SERVICE.Tag handling

// Tags returns cloudformation service tags.
func Tags(tags tftags.KeyValueTags) []*cloudformation.Tag {
	result := make([]*cloudformation.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &cloudformation.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from cloudformation service tags.
func KeyValueTags(ctx context.Context, tags []*cloudformation.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns cloudformation service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*cloudformation.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets cloudformation service tags in Context.
func setTagsOut(ctx context.Context, tags []*cloudformation.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(KeyValueTags(ctx, tags))
	}
}