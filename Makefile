.PHONY: e2e frontend-install frontend-build

e2e-auth:
	chmod +x tests/e2e_test_auth.sh
	./tests/e2e_test_auth.sh
e2e-student:
	chmod +x tests/e2e_student.sh
	./tests/e2e_student.sh
e2e-tracks:
	chmod +x tests/e2e_tracks.sh
	./tests/e2e_tracks.sh

frontend-install:
	cd frontend && npm install

frontend-build:
	cd frontend && npm run build