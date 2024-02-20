.PHONY: freeze
freeze:
	pip freeze > requirements.txt

.PHONY: protogen
protogen:
	python -m grpc_tools.protoc -I./api/oracle --python_out=./app/api/oracle --pyi_out=./app/api/oracle --grpc_python_out=./app/api/oracle ./api/oracle/oracle.proto
