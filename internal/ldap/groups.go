package ldap

import (
	"fmt"

	goldap "github.com/go-ldap/ldap/v3"
)

// GetAllGroups queries groups from AD
func (c *Client) GetAllGroups(filter string) ([]Group, error) {
	if filter == "" {
		filter = "(objectClass=group)"
	}

	searchRequest := goldap.NewSearchRequest(
		c.config.BaseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{
			AttrSAMAccountName, AttrDescription, AttrGroupType, AttrMember,
		},
		nil,
	)

	sr, err := c.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching for groups: %w", err)
	}

	var groups []Group
	for _, entry := range sr.Entries {
		groupType := 0
		fmt.Sscanf(entry.GetAttributeValue(AttrGroupType), "%d", &groupType)

		groups = append(groups, Group{
			DN:             entry.DN,
			SAMAccountName: entry.GetAttributeValue(AttrSAMAccountName),
			Description:    entry.GetAttributeValue(AttrDescription),
			GroupType:      groupType,
			Member:         entry.GetAttributeValues(AttrMember),
		})
	}

	return groups, nil
}

// GetGroupBySAM fetches a single group by sAMAccountName.
func (c *Client) GetGroupBySAM(sam string) (*Group, error) {
	safeSAM := goldap.EscapeFilter(sam)
	searchRequest := goldap.NewSearchRequest(
		c.config.BaseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=group)(sAMAccountName=%s))", safeSAM),
		[]string{
			AttrSAMAccountName, AttrDescription, AttrGroupType, AttrMember,
		},
		nil,
	)

	sr, err := c.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching for group %q: %w", sam, err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("group %q not found", sam)
	}

	entry := sr.Entries[0]
	groupType := 0
	fmt.Sscanf(entry.GetAttributeValue(AttrGroupType), "%d", &groupType)

	return &Group{
		DN:             entry.DN,
		SAMAccountName: entry.GetAttributeValue(AttrSAMAccountName),
		Description:    entry.GetAttributeValue(AttrDescription),
		GroupType:      groupType,
		Member:         entry.GetAttributeValues(AttrMember),
	}, nil
}

// CreateGroup creates a new security group in AD under the default Users container.
func (c *Client) CreateGroup(name, description string) error {
	// Derive the default Users OU from BaseDN (e.g. DC=in,DC=ibict,DC=br → CN=Users,DC=in,DC=ibict,DC=br)
	usersOU := "CN=Users," + c.config.BaseDN
	dn := fmt.Sprintf("CN=%s,%s", goldap.EscapeFilter(name), usersOU)

	addRequest := goldap.NewAddRequest(dn, nil)
	addRequest.Attribute("objectClass", []string{"top", "group"})
	addRequest.Attribute(AttrSAMAccountName, []string{name})
	// Group type: -2147483646 = security, global (0x80000002)
	addRequest.Attribute(AttrGroupType, []string{"-2147483646"})

	if description != "" {
		addRequest.Attribute(AttrDescription, []string{description})
	}

	if err := c.conn.Add(addRequest); err != nil {
		return fmt.Errorf("failed to create group %q: %w", name, err)
	}

	return nil
}

// UpdateGroup modifies an existing group's attributes by DN.
func (c *Client) UpdateGroup(dn, description string) error {
	modRequest := goldap.NewModifyRequest(dn, nil)

	if description != "" {
		modRequest.Replace(AttrDescription, []string{description})
	}

	if err := c.conn.Modify(modRequest); err != nil {
		return fmt.Errorf("failed to update group %q: %w", dn, err)
	}

	return nil
}

// DeleteGroup removes a group by its DN.
func (c *Client) DeleteGroup(dn string) error {
	delRequest := goldap.NewDelRequest(dn, nil)
	if err := c.conn.Del(delRequest); err != nil {
		return fmt.Errorf("failed to delete group %q: %w", dn, err)
	}
	return nil
}
