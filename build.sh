BUILD_TIME=$(date +%s) && GIT_COMMIT=$(git log -1 --format=%h) && docker-compose build --build-arg GIT_COMMIT=$GIT_COMMIT --build-arg BUILD_TIME=$BUILD_TIME app && docker-compose up -d
