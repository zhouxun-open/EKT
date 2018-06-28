is_testing=${1-"x"}
NAME=ekt8
REGISTRY=registry.cloudhua.com/ekt8

version=$(python version.py)
image=${REGISTRY}/${NAME}:${version}
latest=${REGISTRY}/${NAME}:latest

if [ ${is_testing} = "--testing" ];then
    image=${image}-testing
fi

docker build -f docker/Dockerfile -t ${image} .

docker push ${image}

if [ ${is_testing} != "--testing" ];then
    docker tag ${image} ${latest}
    docker push ${latest}
fi
