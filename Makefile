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
e2e-recomm:
	chmod +x tests/e2e_recommendations.sh
	./tests/e2e_recommendations.sh
e2e-recomm:
	chmod +x tests/e2e_track_selection.sh
	./tests/e2e_track_selection.sh

frontend:
	сd frontend && npm install
	cd frontend && npm run build
