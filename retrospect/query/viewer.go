package query

type ViewerName struct {
	Viewer struct {
		Login string
	}
}

func (c *Client) ViewerName(userName string) (string, error) {
	if userName != "" {
		return userName, nil
	}

	variables := map[string]interface{}{}
	var query ViewerName
	if err := c.GQL.Query("ViewerName", &query, variables); err != nil {
		return "", err
	}
	return query.Viewer.Login, nil
}
