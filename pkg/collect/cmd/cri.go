package cmd

const (
	ContainerdGetTagCommand    = "crictl images --digests | grep $(crictl ps | grep %s  | awk '{ print $2 }') | awk '{ print $2 }'"
	ContainerdGetSha256Command = "crictl inspect $(crictl ps | grep %s | awk '{ print $1 }')  | grep image | sed -n 2p | awk -F ':' '{ print $3 }' | awk -F '\"' '{ print $1 }'"
	DockerGetTagCommand        = "docker images --digests | grep `docker inspect $(docker ps | grep -v pause-amd | grep %s | awk '{ print $1 }') | grep Image | sed -n 2p | awk '{ print $2 }' | awk -F ',' '{ print $1 }' | awk -F ':' '{ print $2 }' | cut -c1-12` | awk '{ print $2 }'"
	DockerGetSha254Command     = "docker inspect $(docker ps | grep -v pause-amd | grep %s | awk '{ print $1 }') | grep Image | sed -n 1p | awk -F ':' '{ print $3 }' | awk -F '\"' '{ print $1 }'"
	Tag                        = "tag"
	Sha256                     = "sha256"
)
