package bthreads

import "errors"

func validateInstance(param *Config) (*Config, error) {
	// Verify Instance
	if param.FuncGoroutinesCount < 0 {
		return nil, errors.New("'Config.FuncGoroutinesCount' cannot < 0")
	}
	if param.GoroutinesDelay < 0 {
		return nil, errors.New("'Config.GoroutinesDelay' cannot < 0")
	}
	if param.LogDelay < 0 {
		return nil, errors.New("'Config.LogDelay' cannot < 0")
	}
	if param.StartDelay < 0 {
		return nil, errors.New("'Config.StartDelay' cannot < 0")
	}

	// Set Default
	if param.Name == "" {
		param.Name = dName
	}
	if param.FuncGoroutinesCount == 0 {
		param.FuncGoroutinesCount = dFuncGoroutinesCount
	}
	if param.GoroutinesDelay == 0 {
		param.GoroutinesDelay = dGoroutinesDelay
	}
	if param.LogDelay == 0 {
		param.LogDelay = dLogDelay
	}
	if param.StartDelay == 0 {
		param.StartDelay = dStartDelay
	}

	return param, nil
}
