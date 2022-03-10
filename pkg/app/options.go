// Created by Hisen at 2022/3/1.
package app

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
)

type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets)
	Validate() []error
}

type ConfigurableOptions interface {
	ApplyFlags() []error
}

type CompleteableOptions interface {
	Complete() error
}

type PrintableOptions interface {
	String() string
}
