.PHONY: build check lint test

override empty :=
override space := $(empty) $(empty)
override makefile_space := __GO_RATELIMITER_MAKEFILE_SPACE__
override encoded_makefile_list := $(patsubst $(makefile_space)%,%,$(subst $(space),$(makefile_space),$(MAKEFILE_LIST)))
override ROOT := $(subst $(makefile_space),$(space),$(abspath $(dir $(lastword $(encoded_makefile_list)))))

lint test build: check

check:
	@"$(ROOT)/scripts/check-baseline.sh"
