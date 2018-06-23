import tensorflow as tf

sess = tf.Session()
p = tf.placeholder(tf.uint8, 1)
x = tf.one_hot(p, depth=10)
res = sess.run(x, feed_dict={p: [3]})
print(res)
