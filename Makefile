.PHONY: test-e2e

test-e2e:
	chmod +x tests/e2e_test_auth.sh
	./tests/e2e_test_auth.sh