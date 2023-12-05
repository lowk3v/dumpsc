inc-version:
	@version=$$(cat VERSION); \
	if [ -z "$$version" ]; then \
	  	ver="0.1.0"; \
		echo $$ver > VERSION; \
		echo "Current version is $$ver"; \
		exit 1; \
	else \
		echo "Current version is $$version"; \
		ver=$$(echo $$version | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g'); \
		echo "New version is "$$ver; \
		echo $$ver > VERSION; \
	fi;
