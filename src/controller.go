package src

// 执行
func Run(instance string) error {
	config, err := getInstanceConfig(instance)
	if err != nil {
		return err
	}
	collector, err := NewCollector(instance, config.FromType, config.FromConfig)
	if err != nil {
		return err
	}
	dispenser, err := NewDispenser(instance, config.ToType, config.ToConfig)
	if err != nil {
		return err
	}
	return collector.Run(func(event *Event) error {
		return dispenser.Send(event)
	})
}
