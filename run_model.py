import argparse
import logging
import os
import sys

import numpy

from chars import CHARS


def setup():
    parser = argparse.ArgumentParser(formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument("-m", "--model", required=True,
                        help="Path to the trained model in Tensorflow GraphDef format.")
    logging.basicConfig(level=logging.ERROR)
    os.environ["TF_CPP_MIN_LOG_LEVEL"] = "3"
    return parser.parse_args()


def main():
    args = setup()
    import tensorflow as tf
    tf.logging.set_verbosity(tf.logging.ERROR)
    graph_def = tf.GraphDef()
    with open(args.model, "rb") as fin:
        graph_def.ParseFromString(fin.read())
    tf.import_graph_def(graph_def, name="")
    del graph_def
    graph = tf.get_default_graph()
    input1 = graph.get_operation_by_name("input_1_1").outputs[0]
    input2 = graph.get_operation_by_name("input_2_1").outputs[0]
    output = graph.get_operation_by_name("output").inputs[0]
    batch_size, length = (d.value for d in input1.shape)
    text = sys.stdin.read()
    batches = [[], []]

    def gen_sample(x: int, text: str):
        before = numpy.zeros(length, dtype=numpy.uint8)
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
        buf = [""] * length
        while i < len(text) and j >= 0:
            if text[i] not in ("\x02", "\x03"):
                buf[j] = text[i]
                after[j] = CHARS.get(text[i], len(CHARS))
                j -= 1
            i += 1
        return before, after

    for i in range(length // 2, len(text) - length // 2):
        before, after = gen_sample(i, text)
        batches[0].append(before)
        batches[1].append(after)

    inputs_size = len(batches[0])
    if inputs_size == 0:
        return
    target_size = int(numpy.ceil(inputs_size / batch_size)) * batch_size
    while len(batches[0]) < target_size:
        batches[0].append(numpy.zeros(length, dtype=numpy.uint8))
        batches[1].append(numpy.zeros(length, dtype=numpy.uint8))
    batches[0] = numpy.array(batches[0])
    batches[1] = numpy.array(batches[1])
    result = []
    for bi in range(target_size // batch_size):
        batch1 = batches[0][bi * batch_size:(bi + 1) * batch_size]
        batch2 = batches[1][bi * batch_size:(bi + 1) * batch_size]
        with tf.Session() as session:
            probs = session.run(output, {input1: batch1, input2: batch2})
            result.extend(numpy.argmax(probs, axis=-1))

    result = result[:inputs_size]
    if result[-1] == 0:
        result[-1] = 2  # never ends with code
    print()
    sys.stdout.write(text[:length // 2])
    for i, (x, r) in enumerate(zip(text[length // 2:-length // 2], result)):
        if r == 2:
            sys.stdout.write(x)
            continue
        if r == 0:
            sys.stdout.write(x + "<code>")
        else:
            sys.stdout.write(x + "</code>")
    sys.stdout.write(text[-length // 2:])

if __name__ == "__main__":
    sys.exit(main())
