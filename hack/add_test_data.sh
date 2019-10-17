#!/bin/bash

main() {
  # zip-proto
  post-test-data
  # clean
}

random-version() {
  let "major = $RANDOM % 10"
  let "minor = $RANDOM % 10"
  let "patch = $RANDOM % 10"
  echo ${major}.${minor}.${patch}
}

zip-proto() {
  cd test_protobuf && zip -r proto.zip * && cd ..
}

clean() {
  rm ./test_protobuf/proto.zip
}

post-test-data() {
  if [[ -z "${NUM}" ]] ; then
    NUM=10
  fi
  if [[ -z "${REGISTRY_HOST}" ]] ; then
    REGISTRY_HOST="localhost:8080"
  fi
  # b64data=$(cat ./test_protobuf/proto.zip | base64 --wrap=0)
  messagenames=($(shuf -n ${NUM}  /usr/share/dict/words | sed "s/'//g" | tr '[:upper:]' '[:lower:]'))
  for msg in "${messagenames[@]}" ; do
    # curl \
    #   -X POST ${REGISTRY_HOST}/api/proto \
    #   --data "
    #     {
    #       \"version\": \"$(random-version)\",
    #       \"name\": \"$msg-proto\",
    #       \"body\": \"${b64data}\"
    #     }"
    go run util.go upload "${msg}-proto" "$(random-version)" ./test_protobuf/
  done
}

if [[ "${0}" == "add_test_data.sh" ]] ; then
  main
fi
