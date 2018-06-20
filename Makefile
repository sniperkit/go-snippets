TARGET = greeting
BUILD = go build -o ${TARGET} -tags

.PHONY: run
run:
	@make first
	./${TARGET}
	@echo "-----"

	@make second
	./${TARGET}
	@echo "-----"

	@make third
	./${TARGET}
	@echo "-----"

	@make fourth
	./${TARGET}
	env name="lyc" ./${TARGET}
	@echo "-----"

	@make fifth
	./${TARGET}
	env name="lyc" ./${TARGET}
	@echo "-----"

	@make sixth
	./${TARGET}
	env name="lyc" ./${TARGET}
	env name="now" ./${TARGET}
	env name="now" ./${TARGET}
	@echo "-----"

.PHONY: first
first:
	${BUILD} first

.PHONY: second
second:
	${BUILD} second

.PHONY: third
third:
	${BUILD} third

.PHONY: fourth
fourth:
	${BUILD} fourth

.PHONY: fifth
fifth:
	${BUILD} fifth

.PHONY: sixth
sixth:
	${BUILD} sixth