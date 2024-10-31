package config

import "time"

type Config struct {
	Chain struct {
		ChainId string
		RpcAddr string
	} `toml:"chain"`

	Bench struct {
		// Number of concurrent workers
		Concurrency uint

		// Quantity of account per transaction
		Quantity uint

		// Duration of the benchmark
		Duration duration

		// Timeout of each transaction
		Timeout duration
	} `toml:"bench"`

	Option struct {
		DefaultPrivateKey string
	}
}

type duration time.Duration

func (d *duration) UnmarshalText(text []byte) error {
	if x, err := time.ParseDuration(string(text)); err != nil {
		return err
	} else {
		*d = duration(x)
	}
	return nil
}
