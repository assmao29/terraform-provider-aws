// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package cloudtrail

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// listTags lists cloudtrail service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn cloudtrailiface.CloudTrailAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &cloudtrail.ListTagsInput{
		ResourceIdList: aws.StringSlice([]string{identifier}),
	}

	output, err := conn.ListTagsWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.ResourceTagList[0].TagsList), nil
}

// ListTags lists cloudtrail service tags and set them in Context.
// It is called from outside this package.
func (p *servicePackage) ListTags(ctx context.Context, meta any, identifier string) error {
	tags, err := listTags(ctx, meta.(*conns.AWSClient).CloudTrailConn(ctx), identifier)

	if err != nil {
		return err
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(tags)
	}

	return nil
}

// []*SERVICE.Tag handling

// Tags returns cloudtrail service tags.
func Tags(tags tftags.KeyValueTags) []*cloudtrail.Tag {
	result := make([]*cloudtrail.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &cloudtrail.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from cloudtrail service tags.
func KeyValueTags(ctx context.Context, tags []*cloudtrail.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns cloudtrail service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*cloudtrail.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets cloudtrail service tags in Context.
func setTagsOut(ctx context.Context, tags []*cloudtrail.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates cloudtrail service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn cloudtrailiface.CloudTrailAPI, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.CloudTrail)
	if len(removedTags) > 0 {
		input := &cloudtrail.RemoveTagsInput{
			ResourceId: aws.String(identifier),
			TagsList:   Tags(removedTags),
		}

		_, err := conn.RemoveTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.CloudTrail)
	if len(updatedTags) > 0 {
		input := &cloudtrail.AddTagsInput{
			ResourceId: aws.String(identifier),
			TagsList:   Tags(updatedTags),
		}

		_, err := conn.AddTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates cloudtrail service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).CloudTrailConn(ctx), identifier, oldTags, newTags)
}