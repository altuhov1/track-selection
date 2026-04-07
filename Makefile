.PHONY: test-e2e

test-e2e:
	chmod +x tests/e2e_test_auth_student.sh
	./tests/e2e_test_auth_student.sh