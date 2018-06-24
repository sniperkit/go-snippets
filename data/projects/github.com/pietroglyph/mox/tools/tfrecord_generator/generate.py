r"""Convert raw Magic: The Gathering data into a set of ".tfrecord" files.

Example usage:
    python generate.py --logtostderr \
      --data_dir=${DATA_DIR} \
      --output_dir="${OUTPUT_DIR}"
"""
from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

import random
import hashlib
import io
import json
import os
import glob
import json
import math
import numpy as np

import tensorflow as tf

from PIL import Image, ImageEnhance, ImageDraw, ImageFilter
from object_detection.utils import dataset_util

flags = tf.app.flags
tf.flags.DEFINE_string('data_dir', './raw_data/',
                       'Training image directory.')
tf.flags.DEFINE_string('output_dir', '/tmp/', 'Output data directory.')
tf.flags.DEFINE_boolean('show_debug', False, 'Display a debug image with bounding boxes drawn.')
tf.flags.DEFINE_string('random_seed', '', 'Seed to supply to the random number generator.')

FLAGS = flags.FLAGS
MAX_UNDERFITTING_PIXELS = 200
MAX_COLOR_BALANCE_DELTA = 0.3
MAX_BRIGHTNESS_DELTA = 0.3
TRAINING_SEGMENTS = 2
BACKGROUND_NAMES = ["background1.png", "background2.png", "background3.png", "background4.png"]
NUM_EVAL_EXAMPLES = 10
MAX_ROTATION_DEGREES = 8
MAX_BLUR_RADIUS = 1.3

tf.logging.set_verbosity(tf.logging.INFO)

def randomise_image(foreground_bytes, background_path):
    im = Image.open(foreground_bytes)
    bg = Image.open(background_path, mode="r")

    origsize = im.size

    # Decide on a rotation
    rotation = random.randint(-MAX_ROTATION_DEGREES, MAX_ROTATION_DEGREES)

    # Downscale
    after_rotation_extra_size = (np.array(get_rotated_bounding_box_size(im.size, rotation))-im.size)
    im.thumbnail((bg.size[0]-random.randint(0, MAX_UNDERFITTING_PIXELS)-after_rotation_extra_size[0], bg.size[1]-random.randint(0, MAX_UNDERFITTING_PIXELS))-after_rotation_extra_size[1], Image.LANCZOS)

    after_thumbnail_size = np.array(im.size)

    # Rotate
    im = im.rotate(rotation, Image.BICUBIC, True)

    # Position the image randomly
    margin = (np.array(bg.size)//2)-(np.array(im.size)//2)
    transform = (random.randint(-margin[0],margin[0]), random.randint(-margin[1],margin[1]))
    bg.paste(im, tuple(margin+transform), im)

    # Apply some random blur
    bg = bg.filter(ImageFilter.GaussianBlur(radius=random.uniform(0, MAX_BLUR_RADIUS)))

    # Apply some random bgage enhancements (color balance, brightness, sharpness, contrast, etc.)
    bg = ImageEnhance.Color(bg).enhance(random.uniform(-MAX_COLOR_BALANCE_DELTA, MAX_COLOR_BALANCE_DELTA) + 1)
    bg = ImageEnhance.Brightness(bg).enhance(random.uniform(-MAX_BRIGHTNESS_DELTA, MAX_BRIGHTNESS_DELTA) + 1)

    bg.format = "PNG" # Hmm... Interesting.

    # bg.size should == (1024, 1024)
    return bg, tuple(margin+transform), after_thumbnail_size / origsize, origsize, rotation, (((np.array(bg.size)//2)-(np.array(after_thumbnail_size)//2))+transform)-(margin+transform)

def rotate_point(point, center, rotation_degrees):
    r = math.radians(rotation_degrees)
    return (math.cos(r) * (point[0] - center[0]) + math.sin(r) * (point[1] - center[1]) + center[0],
    math.sin(r) * -(point[0] - center[0]) + math.cos(r) * (point[1] - center[1]) + center[1])

def get_rotated_bounding_box_size(size, rotation_degrees):
    if math.fabs(rotation_degrees) != rotation_degrees:
        rotation_degrees += 90
        size = size[::-1]
    r = math.radians(rotation_degrees)
    return ((math.cos(r)*size[0])+(math.sin(r)*size[1]), (math.cos(r)*size[1])+(math.sin(r)*size[0]))

def create_tf_record_from_annotations(in_dir, in_ids, out_path):
    writer = tf.python_io.TFRecordWriter(out_path)

    for v in in_ids:
        background_png_path = os.path.join(in_dir, random.choice(BACKGROUND_NAMES))

        json_path = os.path.join(in_dir, v + ".json")
        with tf.gfile.GFile(json_path, "r") as f:
            json_raw = f.read()
        annotation = json.loads(json_raw)

        filename = v + ".png"
        card_png_path = os.path.join(in_dir, filename)
        print(card_png_path)
        with tf.gfile.GFile(card_png_path, 'rb') as fid:
            encoded_png = fid.read()
        randomised_image, card_transform, thumbnail_resize_factor, card_original_size, rotation, content_transform = randomise_image(io.BytesIO(encoded_png), background_png_path)
        if randomised_image.format != "PNG":
            raise ValueError("Image format not PNG")

        annotation["card"] = {
            "startX": 0,
            "startY": 0,
            "endX": card_original_size[0],
            "endY": card_original_size[1],
        }

        startx = []
        starty = []
        endx = []
        endy = []
        category_id = []
        category_name = []

        if FLAGS.show_debug:
            debug_draw = ImageDraw.Draw(randomised_image)

        card_rotated_center = None
        for k in sorted(annotation, key=category_name_to_id, reverse=True):
            if k != "uuid":
                start = (annotation[k]["startX"], annotation[k]["startY"])
                end = (annotation[k]["endX"], annotation[k]["endY"])

                if tuple(start) > tuple(end):
                    start, end = end, start # Swap the variables

                start *= np.array(thumbnail_resize_factor)

                if k != "card":
                    start += content_transform

                bb_size = np.array((math.fabs(annotation[k]["endX"] - annotation[k]["startX"]), math.fabs(annotation[k]["endY"] - annotation[k]["startY"])))*thumbnail_resize_factor
                rotated_size = get_rotated_bounding_box_size(bb_size, rotation)

                if k == "card":
                    rotated_center = start + (np.array(rotated_size) / 2)
                    card_rotated_center = rotated_center # We've sorted the dict, so this happens first
                else:
                    rotated_center = (bb_size/2) + start

                rotated_center = rotate_point(rotated_center, card_rotated_center, rotation)

                start = rotated_center-(np.array(rotated_size)/2)
                end = rotated_center+(np.array(rotated_size)/2)

                start += card_transform
                end += card_transform

                if FLAGS.show_debug:
                    debug_draw.rectangle([tuple(start), tuple(end)])
                    for i in xrange(360):
                        debug_draw.point(tuple(np.array(rotate_point(rotated_center, card_rotated_center, i)) + card_transform))

                # Normalize coordinates
                start /= randomised_image.size
                end /= randomised_image.size

                startx.append(start[0])
                starty.append(start[1])
                endx.append(end[0])
                endy.append(end[1])

                category_id.append(category_name_to_id(k))
                category_name.append(k.encode("utf8"))

        if FLAGS.show_debug:
            randomised_image.show()
            raw_input("Press any key to show next image.")
            continue

        # This is wildly inefficient, but I haven't found a better way to get this
        # into a format TensorFlow wants.
        randomised_image_bytes = io.BytesIO()
        randomised_image.save(randomised_image_bytes, format="PNG")

        example = tf.train.Example(features=tf.train.Features(feature={
            "image/width": dataset_util.int64_feature(randomised_image.size[0]),
            "image/height": dataset_util.int64_feature(randomised_image.size[1]),
            "image/filename": dataset_util.bytes_feature(filename.encode("utf8")),
            "image/source_id": dataset_util.bytes_feature(
                annotation["uuid"].encode("utf8")),
            "image/encoded": dataset_util.bytes_feature(randomised_image_bytes.getvalue()),
            "image/format": dataset_util.bytes_feature(b"png"),
            "image/object/bbox/xmin": dataset_util.float_list_feature(startx),
            "image/object/bbox/xmax": dataset_util.float_list_feature(endx),
            "image/object/bbox/ymin": dataset_util.float_list_feature(starty),
            "image/object/bbox/ymax": dataset_util.float_list_feature(endy),
            "image/object/class/text": dataset_util.bytes_list_feature(category_name),
            "image/object/class/label": dataset_util.int64_list_feature(category_id),
        }))

        writer.write(example.SerializeToString())
    writer.close()



def category_name_to_id(name):
    return {
        "cardName": 1,
        "setSymbol": 2,
        "collectorNumber": 3,
        "typeLine": 4,
        "card": 5,
        "uuid": 1000, # This has no equivalent in the label map; we make it big for sorting
    }[name]

def main(_):
    assert FLAGS.data_dir, '`data_dir` missing.'
    assert FLAGS.output_dir, '`output_dir` missing.'

    if not tf.gfile.IsDirectory(FLAGS.output_dir):
        tf.gfile.MakeDirs(FLAGS.output_dir)
    train_output_path = os.path.join(FLAGS.output_dir, 'mtg_train.record')
    val_output_path = os.path.join(FLAGS.output_dir, 'mtg_val.record')
    # testdev_output_path = os.path.join(FLAGS.output_dir, 'coco_testdev.record')

    if FLAGS.random_seed != "":
        random.seed(FLAGS.random_seed)

    files = glob.glob(FLAGS.data_dir+"*.json")
    ids = [[], [], []]
    files_added = 0
    segment = 0
    extra_data_added = False
    for i, v in enumerate(files):
        if files_added > len(files) / TRAINING_SEGMENTS:
            if segment == 0 and not extra_data_added:
                extra_data_added = True
                files_added += (NUM_EVAL_EXAMPLES + 1) - (len(files) / TRAINING_SEGMENTS)
            else:
                files_added = 0
                segment += 1
        ids[segment].append(files[i][len(FLAGS.data_dir):-5]) # strip the .json suffix to get the id
        files_added += 1

    print("Using", len(ids[0]), "images for training, and", len(ids[1]), "for evaluation.")

    create_tf_record_from_annotations(FLAGS.data_dir, ids[0], train_output_path)
    create_tf_record_from_annotations(FLAGS.data_dir, ids[1], val_output_path)
    # create_tf_record_from_annotations(FLAGS.data_dir, ids[2], testdev_output_path)

if __name__ == '__main__':
    tf.app.run()
