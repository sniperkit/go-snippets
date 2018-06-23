# coding: utf-8
import tensorflow as tf
a=tf.Variable(0, name='counter')
b=tf.constant(1)
#new=tf.add(a,b)
#update = tf.assign(a, new)
update = tf.assign(a, tf.add(a,b))
init = tf.global_variables_initializer()
with tf.Session() as sess:
    sess.run(init)
    for _ in range(3):
#        sess.run(new)
        print sess.run(update)
        
with tf.Session() as sess:
    sess.run(init)
    for _ in range(3):
        sess.run(update)
        print sess.run(a)
