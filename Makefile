.PHONY: frontend

frontend:
	сd frontend && npm install
	cd frontend && npm run build
