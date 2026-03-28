package kibela

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
)

// GetGroups retrieves groups from Kibela.
func (c *Client) GetGroups(ctx context.Context, first int) ([]Group, error) {
	var query GetGroupsQuery
	variables := map[string]interface{}{
		"first": graphql.Int(first),
	}

	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	groups := make([]Group, len(query.Groups.Edges))
	for i, edge := range query.Groups.Edges {
		groups[i] = Group{
			ID:         idToString(edge.Node.ID),
			Name:       string(edge.Node.Name),
			IsDefault:  bool(edge.Node.IsDefault),
			IsArchived: bool(edge.Node.IsArchived),
		}
	}

	return groups, nil
}
