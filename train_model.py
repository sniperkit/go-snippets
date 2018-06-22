import argparse
from collections import defaultdict
import gzip
import logging
import os
import pickle
import random
import re
import sys
from typing import List, Tuple, Dict, Optional

import humanize
import numpy

from chars import CHARS, WEIGHTS, OOV_WEIGHT


def setup():
    parser = argparse.ArgumentParser(formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument("-i", "--input", help="Path to the train dataset.")
    parser.add_argument("-l", "--layers", default="600",
                        help="Layers configuration: number of neurons on each layer separated by "
                             "comma.")
    parser.add_argument("-m", "--length", type=int, default=100, help="RNN sequence length.")
    parser.add_argument("-b", "--batch-size", type=int, default=128, help="Batch size.")
    parser.add_argument("-e", "--epochs", type=int, default=3, help="Number of epochs.")
    parser.add_argument("-t", "--type", default="LSTM",
                        choices=("GRU", "LSTM", "CuDNNLSTM", "CuDNNGRU"),
                        help="Recurrent layer type to use.")
    parser.add_argument("-v", "--validation", type=float, default=0.2,
                        help="Fraction of the dataset to use for validation.")
    parser.add_argument("--negative-code-samples", type=float, default=0.5,
                        help="Ratio of negative code boundary samples to the overall number.")
    parser.add_argument("-o", "--output", required=True,
                        help="Path to the resulting Tensorflow graph.")
    parser.add_argument("--snapshot", help="RNN snapshot to load.")
    parser.add_argument("--code-samples", default="code_samples.pickle",
                        help="Cached pickle with the dataset to train Code Neuron.")
    parser.add_argument("--optimizer", default="Adam", choices=("RMSprop", "Adam"),
                        help="Optimizer to apply.")
    parser.add_argument("--dropout", type=float, default=0, help="Dropout ratio.")
    parser.add_argument("--lr", default=0.001, type=float, help="Learning rate.")
    parser.add_argument("--decay", default=0.00002, type=float, help="Learning rate decay.")
    parser.add_argument("--enable-weights", action="store_true",
                        help="Weight character classes.")
    parser.add_argument("--seed", type=int, default=7, help="Random seed.")
    parser.add_argument("--devices", default="0,1", help="Devices to use. Empty means CPU.")
    parser.add_argument("--tensorboard", default="tb_logs",
                        help="TensorBoard output logs directory.")
    logging.basicConfig(level=logging.INFO)
    args = parser.parse_args()
    numpy.random.seed(args.seed)
    random.seed(args.seed)
    return args


def read_dataset(path: str, min_length: int, clean_code: bool, analyze_chars: bool) \
        -> Tuple[List[str], Dict[str, List[int]]]:
    log = logging.getLogger("reader")
    texts = []
    bufsize = 1 << 20
    buffer = bytearray(bufsize)
    bufpos = 0
    read_buffer = bytearray(bufsize)
    chars = defaultdict(list)

    def append_text(end_index: int):
        text = buffer[:end_index].decode("utf-8")
        if analyze_chars:
            index = len(texts)
            for c in sorted(set(text)):
                chars[c].append(index)
        if clean_code:
           text = text.replace("\x02", "").replace("\x03", "")
        if len(text) >= min_length:
            texts.append(text)

    with gzip.open(path) as gzf:
        size = gzf.readinto(read_buffer)
        while size > 0:
            sys.stderr.write("%d texts\r" % len(texts))
            rpos = 0
            while rpos < size:
                border = read_buffer.find(b"\x04", rpos)
                if border == -1:
                    border = size
                delta = border - rpos
                if bufpos + delta > bufsize:
                    raise OverflowError(
                        "%d %d %d %d" % (bufpos, delta, bufpos + delta, bufsize))
                buffer[bufpos:bufpos + delta] = read_buffer[rpos:border]
                if border < size:
                    append_text(bufpos + delta)
                    bufpos = 0
                else:
                    bufpos += delta
                rpos = border + 1
            size = gzf.readinto(read_buffer)
        if bufpos > 0:
            append_text(bufpos)
    log.info("%d texts, avg len %d, %d distinct chars, total chars %d",
             len(texts), numpy.mean([len(t) for t in texts]), len(chars),
             sum(len(t) for t in texts))
    return texts, chars


def config_keras():
    import tensorflow as tf
    from keras import backend
    config = tf.ConfigProto()
    config.gpu_options.allow_growth = True
    backend.tensorflow_backend.set_session(tf.Session(config=config))


def create_char_rnn_model(args: argparse.Namespace, classes: int,
                          weights: Optional[List[numpy.ndarray]] = None):
    # this late import prevents from loading Tensorflow too soon
    import tensorflow as tf
    tf.set_random_seed(args.seed)
    from keras import layers, models, initializers, optimizers, metrics
    log = logging.getLogger("model")
    if args.devices:
        dev1, dev2 = ("/gpu:" + dev for dev in args.devices.split(","))
    else:
        dev1 = dev2 = "/cpu:0"

    def add_rnn(device):
        with tf.device(device):
            input = layers.Input(batch_shape=(args.batch_size, args.length), dtype="uint8")
            log.info("Added %s", input)
            embedding = layers.Embedding(
                200, 200, embeddings_initializer=initializers.Identity(), trainable=False)(input)
            log.info("Added %s", embedding)
        layer = embedding
        layer_sizes = [int(n) for n in args.layers.split(",")]
        for i, nn in enumerate(layer_sizes):
            with tf.device(device):
                layer_type = getattr(layers, args.type)
                ret_seqs = (i < len(layer_sizes) - 1)
                try:
                    layer = layer_type(nn, return_sequences=ret_seqs, implementation=2)(layer)
                except TypeError:
                    # implementation kwarg is not present in CuDNN layers
                    layer = layer_type(nn, return_sequences=ret_seqs)(layer)
                log.info("Added %s", layer)
            if args.dropout > 0:
                layer = layers.Dropout(args.dropout)(layer)
                log.info("Added %s", layer)
        return input, layer

    forward_input, forward_output = add_rnn(dev1)
    reverse_input, reverse_output = add_rnn(dev2)
    with tf.device(dev1):
        merged = layers.Concatenate()([forward_output, reverse_output])
        log.info("Added %s", merged)
        dense = layers.Dense(classes, activation="softmax")
        decision = dense(merged)
        log.info("Added %s", decision)
    optimizer = getattr(optimizers, args.optimizer)(lr=args.lr, decay=args.decay)
    log.info("Added %s", optimizer)
    model = models.Model(inputs=[forward_input, reverse_input], outputs=[decision])
    log.info("Compiling...")
    model.compile(optimizer=optimizer, loss="categorical_crossentropy",
                  metrics=[metrics.categorical_accuracy, metrics.top_k_categorical_accuracy])
    if weights:
        log.info("Setting weights...")
        dense_weights = dense.get_weights()
        weights[-len(dense_weights):] = dense_weights
        model.set_weights(weights)
    log.info("Done")
    return model


def train_char_rnn_model(model, dataset: List[str], args: argparse.Namespace):
    from keras import callbacks, utils

    if args.length % 2 != 0:
        raise ValueError("--length must be even")
    log = logging.getLogger("train")
    numpy.random.seed(args.seed)
    log.info("Splitting into validation and train...")
    valid_doc_indices = set(numpy.random.choice(
        numpy.arange(len(dataset)), int(len(dataset) * args.validation), replace=False))
    train_doc_indices = set(range(len(dataset))) - valid_doc_indices
    train_docs = [dataset[i] for i in sorted(train_doc_indices)]
    valid_docs = [dataset[i] for i in sorted(valid_doc_indices)]

    # we cannot reach first l / 2 (forward) and last l / 2 (backward)
    valid_size = sum(len(text) - args.length for text in valid_docs)
    train_size = sum(len(text) - args.length for text in train_docs)
    log.info("train samples: %d\tvalidation samples: %d\t~%.3f",
             train_size, valid_size, valid_size / (valid_size + train_size))

    class Feeder(utils.Sequence):
        def __init__(self, texts: List[str]):
            self.texts = texts
            pos = 0
            self.index = index = numpy.zeros(len(texts) + 1, dtype=numpy.int32)
            for i, text in enumerate(texts):
                pos += len(text) - args.length
                index[i + 1] = pos
            # this takes much memory but is the best we can do.
            self.batches = numpy.arange(pos, dtype=numpy.uint32)
            log.info("Batches occupied %s", humanize.naturalsize(self.batches.size))
            self.on_epoch_end()

        def __len__(self):
            return self.index[-1] // args.batch_size

        def __getitem__(self, item):
            centers = self.batches[item * args.batch_size:(item + 1) * args.batch_size]
            batch = ([numpy.zeros((args.batch_size, args.length), dtype=numpy.uint8)
                      for _ in range(2)],
                     [numpy.zeros((args.batch_size, len(CHARS) + 1), dtype=numpy.float32)])
            text_indices = numpy.searchsorted(self.index, centers + 0.5) - 1
            for bi, (center, text_index) in enumerate(zip(centers, text_indices)):
                text = self.texts[text_index]
                x = center - self.index[text_index]
                assert 0 <= x < self.index[text_index + 1] - self.index[text_index]
                x += args.length // 2
                text_i = x - 1
                batch_i = args.length - 1
                while text_i >= 0 and batch_i >= 0:
                    batch[0][0][bi][batch_i] = CHARS.get(text[text_i], len(CHARS))
                    text_i -= 1
                    batch_i -= 1
                text_i = x + 1
                batch_i = args.length - 1
                while text_i < len(text) and batch_i >= 0:
                    batch[0][1][bi][batch_i] = CHARS.get(text[text_i], len(CHARS))
                    text_i += 1
                    batch_i -= 1
                batch[1][0][bi][CHARS.get(text[x], len(CHARS))] = 1
            return batch

        def on_epoch_end(self):
            log.info("Shuffling")
            numpy.random.shuffle(self.batches)

    log.info("Creating the training feeder")
    train_feeder = Feeder(train_docs)
    log.info("Creating the validation feeder")
    valid_feeder = Feeder(valid_docs)
    log.info("model.fit_generator")
    tensorboard = callbacks.TensorBoard(log_dir=args.tensorboard)
    checkpoint = callbacks.ModelCheckpoint(
        os.path.join(args.tensorboard, "checkpoint_{epoch:02d}_{val_loss:.3f}.hdf5"),
        save_best_only=True)
    if args.enable_weights:
        weights = [min(v, 100) for (c, v) in sorted(WEIGHTS.items())] + [OOV_WEIGHT]
    else:
        weights = None

    class Shuffler(callbacks.Callback):
        def on_epoch_end(self, epoch, logs=None):
            # Bug: it is not called automatically
            train_feeder.on_epoch_end()

    class Presenter(callbacks.Callback):
        def __init__(self):
            super().__init__()
            self.i = 0
            self.text = "I am working on a collection of classes used for video playback and " \
                        "recording. I have one main class which acts like the public interface, " \
                        "with methods like play(), stop(), pause(), record() etc...  Then I " \
                        "have workhorse classes which do the video decoding and video encoding."
            self.batches = []
            for x in range(args.length // 2, len(self.text) - args.length // 2 - 1):
                before = self.text[x - args.length//2:x]
                after = self.text[x + 1:x + 1 + args.length // 2]
                before_arr = numpy.zeros(args.length, dtype=numpy.uint8)
                after_arr = numpy.zeros_like(before_arr)
                for i, c in enumerate(reversed(before)):
                    before_arr[args.length - 1 - i] = CHARS.get(c, len(CHARS))
                for i, c in enumerate(reversed(after)):
                    after_arr[args.length - 1 - i] = CHARS.get(c, len(CHARS))
                self.batches.append((before_arr, after_arr))
            log.info("Testing on %d batches" % len(self.batches))
            self.vocab = [None] * (len(CHARS) + 1)
            for k, v in CHARS.items():
                self.vocab[v] = k
            self.vocab[len(CHARS)] = "?"

        def on_batch_end(self, batch, logs=None):
            if self.i % 10000 == 0:
                predicted = model.predict(self.batches, batch_size=len(self.batches))
                print("".join(self.vocab[numpy.argmax(p)] for p in predicted))

    class LRPrinter(callbacks.Callback):
        def on_epoch_end(self, epoch, logs=None):
            from keras import backend
            lr = self.model.optimizer.lr
            decay = self.model.optimizer.decay
            iterations = self.model.optimizer.iterations
            lr_with_decay = lr / (1. + decay * backend.cast(iterations, backend.dtype(decay)))
            print("Learning rate:", backend.eval(lr_with_decay))

    model.fit_generator(generator=train_feeder,
                        validation_data=valid_feeder,
                        validation_steps=len(valid_feeder),
                        steps_per_epoch=len(train_feeder),
                        epochs=args.epochs,
                        class_weight=weights,
                        callbacks=[tensorboard, checkpoint, Shuffler(), LRPrinter()],
                        use_multiprocessing=True)


def load_char_rnn_model(args: argparse.Namespace):
    from keras.models import load_model
    from keras.layers.recurrent import RNN

    log = logging.getLogger("load")
    log.info("Reading %s", args.snapshot)
    model = load_model(args.snapshot, compile=False)
    args.batch_size, args.length = model.layers[0].input_shape
    lengths = []
    for layer in model.layers:
        if isinstance(layer, RNN):
            args.type = type(layer).__name__
            length = layer.output_shape[-1]
            if not lengths or lengths[-1] != length:
                lengths.append(length)
    args.layers = ",".join(str(l) for l in lengths)
    log.info("Inferred parameters: --batch-size %d --length %d --type %s --layers %s",
             args.batch_size, args.length, args.type, args.layers)
    return model


def bake_code_neuron_dataset(texts: List[str], negative_ratio: float, length: int) \
        -> Tuple[List[Tuple[numpy.ndarray, numpy.ndarray]],
                 List[Tuple[numpy.ndarray, numpy.ndarray]],
                 List[Tuple[numpy.ndarray, numpy.ndarray]]]:
    assert 0 < negative_ratio < 1
    log = logging.getLogger("code_neuron_dataset")
    positive_beg = []
    positive_neg = []
    negative = []
    needle = re.compile("[\x02\x03]")

    def gen_sample(x: int, text: str):
        before = numpy.zeros(length, dtype=numpy.uint8)
        if text[x] in ("\x02", "\x03"):
            i = x - 1
        else:
            i = x
        j = length - 1
        while i >= 0 and j >= 0:
            if text[i] not in ("\x02", "\x03"):
                before[j] = CHARS.get(text[i], len(CHARS))
                j -= 1
            i -= 1
        after = numpy.zeros(length, dtype=numpy.uint8)
        i = x + 1
        j = length - 1
        while i < len(text) and j >= 0:
            if text[i] not in ("\x02", "\x03"):
                after[j] = CHARS.get(text[i], len(CHARS))
                j -= 1
            i += 1
        return before, after

    for text in texts:
        for match in needle.finditer(text):
            arr = positive_beg if match.group() == "\x02" else positive_neg
            arr.append(gen_sample(match.start(), text))
    positive_count = len(positive_beg) + len(positive_neg)
    negative_count = int(negative_ratio * positive_count / (1 - negative_ratio))
    log.info("Positive count: %d (%d, %d)", positive_count, len(positive_beg), len(positive_neg))
    log.info("Negative count: %d", negative_count)
    total_samples = sum(len(text) - length for text in texts)
    numpy.random.seed(7)
    choices = numpy.random.choice(numpy.arange(total_samples, dtype=numpy.int32),
                                  negative_count, replace=False)
    choices.sort()
    pos = 0
    ni = 0
    for ti, text in enumerate(texts):
        if ti % 100 == 0:
            sys.stderr.write("%d\r" % ti)
        delta = len(text) - length
        while pos + delta > choices[ni]:
            x = choices[ni] - pos + length // 2
            while x < len(text) - 1 and (
                    text[x] in ("\x02", "\x03") or text[x + 1] in ("\x02", "\x03")):
                x += 1
            if x == len(text) - 1:
                while x >= 0 and (text[x] in ("\x02", "\x03") or text[x + 1] in ("\x02", "\x03")):
                    x -= 1
                assert x >= 0
            # x and x+1 are not code boundaries => we look at the middle
            negative.append(gen_sample(x, text))
            ni += 1
            if ni == len(choices):
                break
        if ni == len(choices):
            break
        pos += delta
    sys.stderr.write("\n")
    return positive_beg, positive_neg, negative


def train_code_neuron_model(
        model_code,
        samples: Tuple[List[Tuple[numpy.ndarray, numpy.ndarray]],
                       List[Tuple[numpy.ndarray, numpy.ndarray]],
                       List[Tuple[numpy.ndarray, numpy.ndarray]]],
        args: argparse.Namespace):
    log = logging.getLogger("train_cn")
    size = sum(len(samples[i]) for i in range(3))
    val_size = int(size * args.validation)
    val_size -= val_size % args.batch_size
    size = int(val_size / args.validation)
    log.info("Final size: %d", size)
    train_x_before = numpy.zeros((size, args.length), dtype=numpy.uint8)
    train_x_after = numpy.zeros_like(train_x_before)
    train_y = numpy.zeros((size, 3), dtype=numpy.float32)

    def fill(offset: int, arr: List[Tuple[numpy.ndarray, numpy.ndarray]], y: numpy.ndarray):
        for before, after in arr:
            train_x_before[offset] = before
            train_x_after[offset] = after
            train_y[offset] = y
            offset += 1
            if offset == size:
                break
        return offset

    offset = fill(0, samples[0], numpy.array([1, 0, 0], dtype=numpy.float32))
    offset = fill(offset, samples[1], numpy.array([0, 1, 0], dtype=numpy.float32))
    fill(offset, samples[2], numpy.array([0, 0, 1], dtype=numpy.float32))

    from keras import callbacks

    tensorboard = callbacks.TensorBoard(log_dir=args.tensorboard)
    checkpoint = callbacks.ModelCheckpoint(
        os.path.join(args.tensorboard, "checkpoint_{epoch:02d}_{val_loss:.3f}.hdf5"),
        save_best_only=True)
    model_code.fit([train_x_before, train_x_after], train_y,
                   batch_size=args.batch_size, validation_split=args.validation,
                   epochs=args.epochs, callbacks=[tensorboard, checkpoint])


def export_model(model, path: str):
    from keras import backend
    import tensorflow as tf
    from tensorflow.python.framework import graph_util, graph_io

    log = logging.getLogger("export")
    log.info("Exporting %s to %s", model, path)
    session = backend.get_session()
    tf.identity(model.outputs[0], name="output")
    graph_def = session.graph.as_graph_def()
    # reset the devices
    for node in graph_def.node:
        node.device = ""
    constant_graph = graph_util.convert_variables_to_constants(session, graph_def, ["output"])
    graph_io.write_graph(constant_graph, *os.path.split(path), as_text=False)


def main():
    args = setup()
    try:
        if not args.snapshot:
            dataset, _ = read_dataset(args.input, args.length + 1, True, False)
            config_keras()
            model_char = create_char_rnn_model(args, len(CHARS) + 1)
            train_char_rnn_model(model_char, dataset, args)
            del dataset
        else:
            config_keras()
            model_char = load_char_rnn_model(args)
        if not os.path.exists(args.code_samples):
            if not args.input or not os.path.exists(args.input):
                raise FileNotFoundError("--input %s" % args.input)
            dataset, _ = read_dataset(args.input, args.length + 1, False, False)
            samples = bake_code_neuron_dataset(dataset, args.negative_code_samples, args.length)
            del dataset
            with open(args.code_samples, "wb") as fout:
                pickle.dump(samples, fout, protocol=-1)
        else:
            with open(args.code_samples, "rb") as fin:
                samples = pickle.load(fin)
        model_code = create_char_rnn_model(args, 3, model_char.get_weights())
        del model_char
        train_code_neuron_model(model_code, samples, args)
        export_model(model_code, args.output)
        del model_code
    finally:
        from keras import backend
        backend.clear_session()

if __name__ == "__main__":
    sys.exit(main())
