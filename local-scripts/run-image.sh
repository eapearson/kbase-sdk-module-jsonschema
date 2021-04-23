echo "Starting dev image..."
root=$(git rev-parse --show-toplevel)
container_root=/kb/module
docker run -i -t \
  --name=JSONSchema  \
  --dns=8.8.8.8 \
  -p 5000:5000 \
  --mount type=bind,src=${root}/schemas,dst=${container_root}/schemas \
  --rm test/jsonschema:dev