package pg

func getMockOptionsFalse() *Options {
	return &Options{
		ServerKey:   "abc",
		ClientId:    "abc",
		Logging:     &False,
		Environment: SandBox,
	}
}

func getMockOptionsTrue() *Options {
	return &Options{
		ServerKey:   "abc",
		ClientId:    "abc",
		Logging:     &True,
		Environment: SandBox,
	}
}
