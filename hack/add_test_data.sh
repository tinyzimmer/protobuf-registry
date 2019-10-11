#!/bin/bash

main() {
  zip-proto
  post-test-data
  clean
}

random-version() {
  let "major = $RANDOM % 10"
  let "minor = $RANDOM % 10"
  let "patch = $RANDOM % 10"
  echo ${major}.${minor}.${patch}
}

zip-proto() {
  cd hack/test_protobuf && zip -r proto.zip * && cd ../..
}

clean() {
  rm hack/test_protobuf/proto.zip
}

post-test-data() {
  if [[ -z "${NUM}" ]] ; then
    NUM=10
  fi
  b64data=$(cat hack/test_protobuf/proto.zip | base64 --wrap=0)
  messagenames=($(shuf -n ${NUM}  /usr/share/dict/words | sed "s/'//g" | tr '[:upper:]' '[:lower:]'))
  for msg in "${messagenames[@]}" ; do
    curl \
      -X POST localhost:8080/api/proto \
      --data "
        {
          \"version\": \"$(random-version)\",
          \"name\": \"$msg-proto\",
          \"body\": \"${b64data}\",
          \"remoteDeps\": [
            {
              \"url\": \"github.com/googleapis/api-common-protos\"
            }
          ]
        }"
  done
}

main
