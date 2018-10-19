#!/usr/bin/env bash
set -e

if ! [[ "$0" =~ scripts/kubekins-e2e.build.push.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

$(aws ecr get-login --no-include-email --region us-west-2)

if [[ -z "${GIT_COMMIT}" ]]; then
  GIT_COMMIT=$(git rev-parse --short=12 HEAD || echo "GitNotFound")
fi

_AWS_REGION=us-west-2
if [[ "${AWS_REGION}" ]]; then
  _AWS_REGION=${AWS_REGION}
fi

if [[ -z "${REGISTRY}" ]]; then
  REGISTRY=$(awstester ecr --region=${_AWS_REGION} get-registry)
fi

echo "Building:" ${REGISTRY}/kubekins-e2e:${GIT_COMMIT}

CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) \
  go build -v \
  -ldflags "-s -w \
  -X github.com/aws/awstester/version.GitCommit=${GIT_COMMIT} \
  -X github.com/aws/awstester/version.ReleaseVersion=${RELEASE_VERSION} \
  -X github.com/aws/awstester/version.BuildTime=${BUILD_TIME}" \
  -o ./bin/awstester \
  ./cmd/awstester

docker build \
  --tag ${REGISTRY}/kubekins-e2e:${GIT_COMMIT} \
  --file ./scripts/kubekins-e2e.Dockerfile .

docker push ${REGISTRY}/kubekins-e2e:${GIT_COMMIT}
