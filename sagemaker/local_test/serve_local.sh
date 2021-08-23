#!/bin/sh

image=algoflow-sagemaker

docker run -v $(pwd)/test_dir:/opt/ml -p 8080:8080 --rm ${image} serve
