package ovirtclient

func (o *oVirtClient) GetBlankTemplate(retries ...RetryStrategy) (result Template, err error) {
	templateList, err := o.ListTemplates(retries...)
	if err != nil {
		return nil, err
	}
	for _, tpl := range templateList {
		if tpl.ID() == blankTemplateID {
			return tpl, nil
		}
	}
	for _, tpl := range templateList {
		if tpl.IsBlank() {
			return tpl, nil
		}
	}

	return nil, newError(ENotFound, "No blank template found.")
}
