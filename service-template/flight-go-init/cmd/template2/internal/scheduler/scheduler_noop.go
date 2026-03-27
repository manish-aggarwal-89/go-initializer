package scheduler

// NoopSchedulerCredential is a no-op implementation for the starter pack.
// Replace with a real implementation that uses repository and predefine.
type NoopSchedulerCredential struct{}

func NewNoopSchedulerCredential() SchedulerCredential {
	return &NoopSchedulerCredential{}
}

func (NoopSchedulerCredential) CredentialPredefine() error {
	return nil
}

// NoopSchedulerSystemParameter is a no-op implementation for the starter pack.
// Replace with a real implementation that uses repository and predefine.
type NoopSchedulerSystemParameter struct{}

func NewNoopSchedulerSystemParameter() SchedulerSystemParameter {
	return &NoopSchedulerSystemParameter{}
}

func (NoopSchedulerSystemParameter) RefreshSystemParameterPredifine() error {
	return nil
}
