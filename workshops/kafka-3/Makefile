prepare-zsh-podman:
	echo "alias docker=podman">>~/.zshrc
	echo "alias docker-compose=podman-compose">>~/.zshrc

prepare-bash-podman:
	echo "alias docker=podman">>~/.zshrc
	echo "alias docker-compose=podman-compose">>~/.zshrc

#docker-compose=podman-compose
#docker=podman

dc-up:
	$(docker-compose) up -d

dc-down:
	$(docker-compose) down


pd-up:
	$(docker-compose) up -d

pd-down:
	$(docker-compose) down
