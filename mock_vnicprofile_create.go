package ovirtclient

func (m *mockClient) CreateVNICProfile(
	name string,
	networkID string,
	params OptionalVNICProfileParameters,
	_ ...RetryStrategy,
) (VNICProfile, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if err := validateVNICProfileCreationParameters(name, networkID, params); err != nil {
		return nil, err
	}

	if _, ok := m.networks[networkID]; !ok {
		return nil, newError(ENotFound, "network not found")
	}

	id := m.GenerateUUID()
	m.vnicProfiles[id] = &vnicProfile{
		client: m,

		id:        id,
		networkID: networkID,
		name:      name,
	}

	return m.vnicProfiles[id], nil
}
