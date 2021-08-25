# This is the file that implements a flask server to do inferences. It's the file that you will modify to
# implement the scoring for your own algorithm.

import io
import os
import pickle
import tempfile
import pathlib
from algoflow.app.forecast.cli.run_forecast import prepare_convect_data, run_convect_transform

import flask
import pandas as pd

prefix = "/opt/ml/"
model_path = os.path.join(prefix, "model")
run_id = "fit-run"

# The flask app for serving predictions
app = flask.Flask(__name__)


@app.route("/ping", methods=["GET"])
def ping():
    """Determine if the container is working and healthy. In this sample container, we declare
    it healthy if we can load the model successfully."""

    health = pathlib.Path(model_path).exists()

    status = 200 if health else 404
    return flask.Response(response="\n", status=status, mimetype="application/json")


@app.route("/invocations", methods=["POST"])
def transformation():
    data = None

    # Convert from CSV to pandas
    if flask.request.content_type == "application/json":
        # handle the zip file
        data = flask.request.json

        schema = data.get("schema")
        config = data.get("config")
        freq = data.get("freq")

        if not schema:
            return flask.Response(
                status=415, 
                response="The payload shall contain a schema key"
            )
        if not config:
            return flask.Response(
                status=415, 
                response="The payload shall contain a config key"
            )
        if not freq:
            return flask.Response(
                status=415, 
                response="The payload shall contain a freq key"
            )

        # write to temp file
        temp_input_path = tempfile.mkdtemp()
        with open(os.path.join(temp_input_path, "config.yaml"), 'w') as f:
            f.write(config)
        with open(os.path.join(temp_input_path, "schema.yaml"), 'w') as f:
            f.write(schema)
    else:
        return flask.Response(
            response="This predictor only supports json payload", status=415, mimetype="text/plain"
        )

    # invoke the prepare data
    print("Start preparing data")
    converted_path = tempfile.mkdtemp()

    prepare_convect_data.callback(
            schema=os.path.join(temp_input_path, "schema.yaml"),
            output=converted_path,
            freq=freq,
            scheduler=None,  # using local cluster
            run_id=None,
            workspace=None,
    )

    print("Start transforming data")

    output_path = tempfile.mktemp()
    run_convect_transform.callback(
        data=converted_path,
        output=output_path,
        config=os.path.join(temp_input_path, "config.yaml"),
        scheduler=None,
        workspace=model_path,
        run_id=run_id
    )

    print("Reading prediction result")

    # return the result
    out = io.StringIO()
    pd.read_csv(
        output_path
    ).to_csv(out, index=False)
    result = out.getvalue()

    return flask.Response(response=result, status=200, mimetype="text/csv")
