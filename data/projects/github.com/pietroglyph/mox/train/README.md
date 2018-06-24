# Training
Although a pretrained inference graph is provided, you can train your own with a bit of effort.

*There is a known bug that breaks the training script in Python 3. You can manually fix it according to [this](https://github.com/tensorflow/models/issues/3705#issuecomment-375563179), or you can set everything up for Python 2 and run that way.*

1. Install [Tensorflow](https://www.tensorflow.org/install/)
2. Clone `https://github.com/tensorflow/models.git` and follow [this](https://github.com/tensorflow/models/blob/master/research/object_detection/g3doc/installation.md) guide to install dependencies, compile protocol buffers, and add libraries to `$PYTHONPATH`.
3. ` $ cd models/research/; python setup.py build; python setup.py install` (`models` is the repo that you should've just cloned).
4. Test your installation according to the guide in step 2.
5. ` $ cd mox/; python ./train/tools/tfrecord_generator/generate.py; mv /tmp/*.record ./train/data"`
6. Run the training script (you must replace `{PATH TO MODELS REPO}` with the path to where you cloned the repo in step 2):

  (Current directory must be `mox/train/`.)

  ` $ python {PATH TO TENSORFLOW MODELS REPO}/research/object_detection/train.py     --logtostderr     --pipeline_config_path=./models/text_localize_faster_rcnn/faster_rcnn_resnet101_mtg.config     --train_dir=./models/text_localize_faster_rcnn/train/`

  This should run until you kill it. See [this](https://github.com/tensorflow/models/blob/master/research/object_detection/g3doc/running_locally.md) for more info, especially in regards to running Tensorboard and an evaluation job.
7. Export the inference graph:
  (Current directory must be `mox/train/`. You must replace {LATEST CHECKPOINT NUMBER} with the number of the latest checkpoint.)

  ` $ python {PATH TO TENSORFLOW MODELS REPO}/research/object_detection/export_inference_graph.py     --input_type image_tensor     --pipeline_config_path ./models/text_localize_faster_rcnn/faster_rcnn_resnet101_mtg.config     --trained_checkpoint_prefix ./models/text_localize_faster_rcnn/train/model.ckpt-{LATEST CHECKPOINT NUMBER}     --output_directory inference_graph`

8. Replace the old inference graph:

  (Current directory must be `mox/train/data`.)

  ` $ rm -rf ../inference_graph; mv ./inference_graph/ ../`
# Running the evaluation job
Although the best way to evaluate the performance of the system is often to run it, you can visually verify the performance of the localizer using the following command:

  (Current directory must be `mox/`)

  ` $ python {PATH TO TENSORFLOW MODELS REPO}/research/object_detection/train.py     --logtostderr     --pipeline_config_path=./models/text_localize_faster_rcnn/faster_rcnn_resnet101_mtg.config     --checkpoint_dir=./models/text_localize_faster_rcnn/train/     --eval_dir=./models/text_localize_faster_rcnn/eval/`
# Annotating more data
There are about 300 human-annotated examples provided in `raw_data`, but you may annotate more using the `label_helper` tool.

This assumes a working installation of [Go](https://golang.org/).
1. ` $ cd mox/tools/label_helper`
2. ` $ go get github.com/BlueMonday/go-scryfall; go get github.com/ogier/pflag`
3. ` $ go install`
4. ` $ label_helper -o "../../raw_data/"`
