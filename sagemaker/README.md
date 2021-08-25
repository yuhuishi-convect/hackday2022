# Convect Forecat Sagemaker



## Context

### What
This project ports the Convect Forecast to AWS Sagemaker ecosystem, as a custom container that can support training and inference on user provided datasets.

Covnect Forecast - an automated time series model building tool.

AWS Sagemaker - a full lifecycle model management tool provided by Amazon.

### Why

* AWS sagemaker ecosystem - Sagemaker is a full model lifecycle management tool, including experiement tracking, deployment, parameter tuning... Porting convect forecast to Sagemaker we get those functions for free
* Quicker time to market - if a user is already on AWS, shipping forecasting products in the form of Sagemaker container expedites the time to integrate.
* Data control - if a user does not want her data to flow out of her system, then using Sagemaker for model building can make sure of that because resources are provisioned using the user's own account. 

## Setups 

### Permissions

Running this notebook requires permissions in addition to the normal `SageMakerFullAccess` permissions. This is because we'll creating new repositories in Amazon ECR. The easiest way to add these permissions is simply to add the managed policy `AmazonEC2ContainerRegistryFullAccess` to the role that you used to start your notebook instance. There's no need to restart your notebook instance when you do this, the new permissions will be available immediately.

### Build the container

```
chmod +x ./build_and_push.sh

./build_and_push.sh
```
This builds and push an image named `algoflow-sagemaker` to ECR


## Usage

### Testing

`local_test/train_local.sh` and `local_test/predict_local.sh`

### Training and inference

Training using convect forecast within Sagemaker - see `examples/remote-train.ipynb`

Batch tranformation using convect forecsat - see `examples/remote-train.ipynb`

Deploy a persistent inference service - WIP

## How does this work

### How Amazon SageMaker runs your Docker container 

Sagemaker runs the container with argument `train` or `serve` depending on the task environment. 

#### Training

```
/opt/ml
|-- input
|   |-- config
|   |   |-- hyperparameters.json
|   |   `-- resourceConfig.json
|   `-- data
|       `-- <channel_name>
|           `-- <input data>
|-- model
|   `-- <model files>
`-- output
    `-- failure
```

The input

* /`opt/ml/input/config` contains information to control how your program runs. `hyperparameters.json` is a JSON-formatted dictionary of hyperparameter names to values. These values will always be strings, so you may need to convert them. `resourceConfig.json` is a JSON-formatted file that describes the network layout used for distributed training. Since scikit-learn doesn't support distributed training, we'll ignore it here.
* `/opt/ml/input/data/<channel_name>/` (for File mode) contains the input data for that channel. The channels are created based on the call to CreateTrainingJob but it's generally important that channels match what the algorithm expects. The files for each channel will be copied from S3 to this directory, preserving the tree structure indicated by the S3 key structure.
* `opt/ml/input/data/<channel_name>_<epoch_number>` (for Pipe mode) is the pipe for a given epoch. Epochs start at zero and go up by one each time you read them. There is no limit to the number of epochs that you can run, but you must close each pipe before reading the next epoch.

The output


* `/opt/ml/model/` is the directory where you write the model that your algorithm generates. Your model can be in any format that you want. It can be a single file or a whole directory tree. SageMaker will package any files in this directory into a compressed tar archive file. This file will be available at the S3 location returned in the `DescribeTrainingJob` result.
* `/opt/ml/output` is a directory where the algorithm can write a file failure that describes why the job failed. The contents of this file will be returned in the `FailureReason` field of the `DescribeTrainingJob` result. For jobs that succeed, there is no reason to write this file as it will be ignored.

### Inference

Hosting has a very different model than training because hosting is reponding to inference requests that come in via HTTP. In this example, we use our recommended Python serving stack to provide robust and scalable serving of inference requests:

![](https://raw.githubusercontent.com/aws/amazon-sagemaker-examples/e2e246bbc77c39c7b727e5a92644be668bb2a34c/aws_marketplace/creating_marketplace_products/images/stack.png)

Amazon SageMaker uses two URLs in the container:


* `/ping` will receive GET requests from the infrastructure. Your program returns 200 if the container is up and accepting requests.
* `/invocations` is the endpoint that receives client inference POST requests. The format of the request and the response is up to the algorithm. If the client supplied ContentType and Accept headers, these will be passed in as well.

The container will have the model files in the same place they were written during training:

```
/opt/ml
`-- model
    `-- <model files>
```

### What we did

We add two scripts on top of the current convect forecasting docker container, as entrypoints during the training and inference process - `container/train` and `container/predictor.py`.

`container/train` parses the user provided parameters at `/opt/ml/input/config/hyperparamters.json` and data files at `/opt/ml/input/data/...` then start a normal training process by providing those file locations to it.

`container/predictor.py` implements a simple rest server which gets a json payload with `schema`, `config` and `freq` as keys. Then a batch transformation process is started by providing those arguments as inputs.
The result is returned as a plain text file as response.