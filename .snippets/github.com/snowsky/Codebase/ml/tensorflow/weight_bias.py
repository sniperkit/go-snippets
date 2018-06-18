# coding: utf-8
import tensorflow as tf
import numpy as np
from tensorflow.python import debug as tf_debug

x = np.random.rand(100).astype(np.float32)
y=x*0.1+0.3
#x.reshape(10,10)
#W=tf.Variable(tf.random_uniform([1]))
W=tf.Variable(tf.random_uniform([1], -1.0, 1.0))
b=tf.Variable(tf.zeros([1]))
#y=W*x + b
y=x*0.1+0.3
y_pre=W*x + b
loss = tf.reduce_mean(tf.square(y_pre-y))
optimizer = tf.train.GradientDescentOptimizer(0.5)
train = optimizer.minimize(loss)
init = tf.global_variables_initializer()
sess = tf.Session()
sess.run(init)
        
for step in range(1000):
    sess.run(train)
    if step % 20 == 0:
        print step, sess.run(W), sess.run(b)
        
