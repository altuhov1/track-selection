.PHONY: e2e frontend-install frontend-build

e2e:
	chmod +x tests/e2e_test_auth.sh
	./tests/e2e_test_auth.sh
e2es:
	chmod +x tests/e2e_student.sh
	./tests/e2e_student.sh

frontend-install:
	cd frontend && npm install

frontend-build:
	cd frontend && npm run build