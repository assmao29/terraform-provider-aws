// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ssm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// listTags lists ssm service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn ssmiface.SSMAPI, identifier, resourceType string) (tftags.KeyValueTags, error) {
	input := &ssm.ListTagsForResourceInput{
		ResourceId:   aws.String(identifier),
		ResourceType: aws.String(resourceType),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.TagList), nil
}

// ListTags lists ssm service tags and set them in Context.
// It is called from outside this package.
func (p *servicePackage) ListTags(ctx context.Context, meta any, identifier, resourceType string) error {
	tags, err := listTags(ctx, meta.(*conns.AWSClient).SSMConn(ctx), identifier, resourceType)

	if err != nil {
		return err
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(tags)
	}

	return nil
}

// []*SERVICE.Tag handling

// Tags returns ssm service tags.
func Tags(tags tftags.KeyValueTags) []*ssm.Tag {
	result := make([]*ssm.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &ssm.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from ssm service tags.
func KeyValueTags(ctx context.Context, tags []*ssm.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns ssm service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*ssm.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets ssm service tags in Context.
func setTagsOut(ctx context.Context, tags []*ssm.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(KeyValueTags(ctx, tags))
	}
}

// createTags creates ssm service tags for new resources.
func createTags(ctx context.Context, conn ssmiface.SSMAPI, identifier, resourceType string, tags []*ssm.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	return updateTags(ctx, conn, identifier, resourceType, nil, KeyValueTags(ctx, tags))
}

// updateTags updates ssm service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn ssmiface.SSMAPI, identifier, resourceType string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.SSM)
	if len(removedTags) > 0 {
		input := &ssm.RemoveTagsFromResourceInput{
			ResourceId:   aws.String(identifier),
			ResourceType: aws.String(resourceType),
			TagKeys:      aws.StringSlice(removedTags.Keys()),
		}

		_, err := conn.RemoveTagsFromResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.SSM)
	if len(updatedTags) > 0 {
		input := &ssm.AddTagsToResourceInput{
			ResourceId:   aws.String(identifier),
			ResourceType: aws.String(resourceType),
			Tags:         Tags(updatedTags),
		}

		_, err := conn.AddTagsToResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates ssm service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier, resourceType string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).SSMConn(ctx), identifier, resourceType, oldTags, newTags)
}