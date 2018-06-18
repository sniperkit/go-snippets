import tensorflow as tf

v = tf.Variable(5, dtype=tf.int64)
def fn(x):
    return x + v

dataset = (tf.contrib.data.Dataset.range(10).map(fn))

iterator = dataset.make_initializable_iterator()
next = iterator.get_next()

with tf.Session() as sess:
    sess.run(tf.global_variables_initializer())
    sess.run(iterator.initializer)

    for i in range(3):
        res = sess.run(next)
        print(res)
