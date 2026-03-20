GEN_PATH ?= $(shell pwd)

.PHONY: gen

gen:
	jet -dsn=postgresql://$(TASK_MANAGER_DB_USER):$(TASK_MANAGER_DB_PASSWORD)@$(TASK_MANAGER_DB_HOST):$(TASK_MANAGER_DB_PORT)/$(TASK_MANAGER_DB_NAME)?sslmode=disable \
		-schema=$(TASK_MANAGER_DB_SCEHMA) \
		-path=$(GEN_PATH)/.gen
