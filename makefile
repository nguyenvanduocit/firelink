include .env.makefile
deploy:
	firebase deploy --token "${FIREBASE_TOKEN}"
