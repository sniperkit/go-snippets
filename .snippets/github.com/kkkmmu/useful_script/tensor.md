1. The general outline of TensorFlow programs is as follows:
	a. Import and parse the data sets.
	b. Create feature columns to describe the data.
	c. Select the type of model.
	d. Train the model.
	e. Evaluate the model's effectiveness.
	f. Let the trained model make predictions.
2. Tensorflow seperates definition of computation from their execution.也就是说tensorflow将计算的定义和执行区分对待。
3. Tensorflow 的一般使用步骤：
	1. Assemble a graph: 定义图
	2. Use a session to execute operations in the graph.:利用Session来执行该图定义的操作。
4. What is a tensor ?
	An n-dimensional array.也就是说tensor是一个N维的数组。
	0-d tensor: Scalar(number)
    	1-d tensor: vector
	2-d tensor: matrix
	and so on
	也就是说tensors ar data: Data Flow -> Tensor Flow
5. 利用tensorflow计算两个数的和：
	a. import tensorflow as tf
	b. a = tf.add(3,5)
    	c. sess = tf.Session()
    	d. print(sess.run(a))
    	e. sess.close()

6. tf.Session()：
	A session object encapsulates the environment in which Opeation objects are executed, and Tensor object are evaluated.

7. 利用tensorflow进行更复杂的计算：
    1.
	x = 2
	y = 3
	op1 = tf.add(x,y)
    	op2 = tf.mul(x,y)
    	op3 = tf.pow(op2, op1)
    	with tf.Session() as sess:
		op3 = sess.run(op3)
    2.  x = 2
    	y = 3
	add_op = tf.add(x, y)
    	mul_op = tf.mul(x, y)
    	useless = tf.mul(x, add_op)
    	pow_op = tf.pow(add_op, mul_op)
    	with tf.Session() as sess:
		z = sess.run(pow_op)

8. Subgraphs:
	Possible to break graphs into several chunks and run them parallelly accross multiple CPUS, GPUs or devices.

9. tf.Graph()
    1. Create a graph:
    	g = tf.Graph()
    2. To add operators to a graph, set it as default.
    	g = tf.Graph()
        with g.as_default():
	    x = tf.add(3, 5)
	sess = tf.Session(graph=g)
	with tf.Session() as sess:
	    	sess.run(x)
    3. To handle the default graph:
      	g = tf.get_default_graph
10. Why graphs:
     1. Save computation(only run subgraphs that lead to the values you want to fetch).
     2. Break computation into small, different pieces to facilitates auto-differentiation.
     3. Facilitate distributed computation, spread the work accross multiple CPUs, GPUs or devies.
     4. Many common machine learning models are commonly taught and visualized as directed graphs already.
11. Using of tensorboard:
      1. 
      		import tensorflow as tf
      		a = tf.constant(2)
      		b = tf.constant(3)
      		x = tf.add(a, b)
      		with tf.Session() as sess:
		    writer = tf.summary.FileWriter('./graphs', sess.graph)
		    print(sess.run(x))
		writer.close()
     2. python xxx.py
     3. tensorboard --logdir=".graphs" --port 8080
     4. Visit http://localhost:8080

12. More tensorflow program examples:
	1. 
		import tensorflow as tf
		a = tf.constant(2)
    		b = tf.constant(3)
    		x = tf.add(a, b)
    		with tf.Session() as sess:
			print(sess.run(x))
    	2. 
		import tensorflow as tf
		a = tf.constant([2, 2], name='a')
		b = tf.constant([0, 1], [2, 4], name='b')
		x = tf.add(a, b, name='add')
		y = tf.mul(a, b, name='mul')
		with tf.Session() as sess:
			x, y = sess.run([x, y])
    			print(x,y)
13. Creates a tensor of shape and all elements be zero:
	tf.zeros(shape, dtype=tf.float32, name=None)
    	tf.zeros([2, 3], tf.int32) ==> [[0, 0, 0], [9, 0, 0]]
14. Create a tensor of shape and type(unless type is specified) as the input_tensor but all elements are zeros.
	tf.zeros_like(input_tensor, dtype=None, name=None, optimize=True)
    	# input_tensor: [[0, 1], [2, 3], [4, 5]]
    	tf.zero_like(input_tensor) ==> [[0, 0], [0, 0], [0, 0]]
15. tf.ones(shape, dtype=float32, name=None)
    tf.ones_like(input_tensor, dtype=None, name=None, optimize=True)
    tf.fill(dims, value, name=None)
    	tf.fill([2, 3], 8) ==> [[8, 8, 8], [8, 8, 8]]
    tf.linespace(start, stop, num, nmae=None)
    	tf.linespace(10.0, 13.0, 4) ==> [10.0, 11.0, 12,0, 13.0]
    tf.range(start, limit=None, delta=1, dtype=None, name='range')
    	tf.range(3, 18, 3) ==> [3, 6, 9, 12, 15]
    tf.random_normal(shape, mean=0.0, stddev=1.0 dtype=tf.float32, seed=None, name=None)
    tf.truncated_normal(shape, mean=0.0, stddev=1.0, dtype=tf.float32, seed=None, name=None)
    tf.random_uniform(shape, minval=0, maxval=None, dtype=tf.float32, seed=None, name=None)
    tf.random_shuffle(value, seed=None, name=None)
    tf.random_crop(value, size, seed=None, name=None)
    tf.multinomial(logits, num_samples, seed=None, name=None)
    tf.random_gamma(shape, alpha, beta=None, dtype=tf.float32, seed=None, name=None)
    tf.set_random_seed(seed)
16. Do not use Python native types for tensors because TensorFlow has to infer Python type.
17. What's wrong with constants ？
	1.Constants are stored in the graph definition.
	2.This makees loading graphs expensive when constants are big.
	3.So, only use constants for primitive types. Use variables or readers for more data that requires more memory.
18. Variables: 
	# create variable a with scalar value.
	a = tf.Variable(2, name="scalar")
	# create variable b as a vector
	b = tf.Variable([2, 3], name='vector')
	# create variable c as 2*2 matrix
	c = tf.Variable([[0,1],[2,3]], name='matrix')
	# create variable W as 784 * 10 tensor, filed with zeros
	W = tf.Variable(tf.zeros([784, 10]))
19. tf.Variable is a class, but tf.constant is an operator.
20. 
	W = tf.Variable(tf.truncated_normal([700, 10]))
    	with tf.Session() as sess:
		sess.run(W.initializer)
    		print(W.eval())
21. tf.Variable.assign: tf.Variable.assign does not assign the new value to the variable, it creates an assign operator, and that operator needs to be run to take effect.
	1. W = tf.Variable(10)
    	   W.assign(100)
    	   with tf.Session() as sess:
	   	sess.run(W.initializer)
    		print(W.eval) # >>>>> 10
	2. W = tf.Variable(10)
           assign_op = W.assign(100)
           with tf.Session() as sess:
		# sess.run(W.initializer) # This line is not need in this case.
    		sess.run(assign_op)
    	   print(W.eval()) # >>>>>> 100
	3. my_var = tf.Variable(2, name='my_var')
	   my_var_times_two = my_var.assign(2 * my_var)
    	   with tf.Session() as sess:
	   	sess.run(my_var.initializer)
    		sess.run(my_var_times_two) # >>>>>>>> 4
    		sess.run(my_var_times_two) # >>>>>>>> 8
    		sess.run(my_var_times_two) # >>>>>>>> 16
    		sess.run(my_var_times_two) # >>>>>>>> 32
	4. my_var = tf.Variable(10)
    	   with tf.Session() as sess:
	   	sess.run(my_var.initializer)
    		sess.run(my_var.assign_add(10)) # >>>>>>>> 20
    		sess.run(my_var.assign_sub(2)) # >>>>>>>> 18
22. Each session maintains its own copy of variable
	W = tf.Variable(10)
    	sess1 = tf.Session()
        sess2 = tf.Session()
    	
    	session1.run(W.initializer)
        session2.run(W.initializer)

    	print(sess1.run(W.assign_add(10)) # >>>>> 20
	print(sess2.run(W.assign_sub(2))  # >>>>> 8

	print(sess1.run(W.assign_add(100)) # >>>>> 120
	print(sess2.run(W.assign_sub(50)) # >>>>>> -42

	sess1.close()
	sess2.close()
23. Control Dependencies
	tr.Graph.control_dependencies(control_inputs)
	# define which ops should be run first
	# your graph g have 5 ops: a, b, c, d, e
	with g.control_dependencies([a, b, c]):
		d = ...
		e = ...
		# d and e will only run after 'a', 'b', 'c' have executed.
24. Placeholders
	A TF program often has 2 phases:
		1. Assemble a graph.
		2. Use a session to execute operations in the graph.
	=> Cann assemble the graph first without knowing the values needed for computation.
	=> Analogy: Can define the function f(x, y) = x*2 + y without knowing value of x or y. x, y are placeholders for the acttual values.

	tf.placeholder(dtype, shape=None, name=None)
	# create a placeholder of type float 32-bit, shape is a vector of 3 elements
    	a = tf.placeholder(tf.float32, shape=[3])
    	# create a constant of type float 32-bit, shape is a vector of 3 elements.
    	b = tf.constant([5, 5, 5], tf.float32)
    	# use the placeholder as you would a constant or a variable
    	c = a + b # short for t.f.add(a, b)
    	
    	with tf.Session() as sess:
		print(sess.run(c)) # Error because a doesn't have any value.
25. Feed the values to placeholders using a dictionary.
	tf.placeholder(dtype, shape=None, name=None)
    	# create a placeholder of type float 32-bit, shape is a vector of 3 elements.
    	a = tf.placeholder(tf.float32, shape=[3])
    	
    	# create a constant of type float 32-bit, shape is a vector of 3 elements.
    	b = tf.constant([5, 5, 5], tf.float32)

    	# use the placeholder as you would a constant or a variable
    	c = a + b # short for tf.add(a, b)

    	with tf.Session() as sess:
		# feed[1, 2, 3] to placeholder a via the dict{a: [1, 2, 3]}
		# fetch value of c
		print(sess.run(c, {a:[1, 2, 3]})) # the tensor a is the key, not the string 'a'.
	# >> [6, 7, 8]
26. Placeholders are valid ops. Placeholder is just a way to indicate that something must be fed.
27. Feeding values to TF ops:
	# create operations, tensors, etc(using the default graph)
	a = tf.add(2, 5)
    	b = tf.mul(a, 3)
    	with tf.Session() as sess:
		# define a dictionary that says to replace the value of 'a' with 15
		replace_dict = {a: 15}

		# Run the session, passing in 'replace_dict' as the value to 'feed_dict'
		sess.run(b, feed_dict=replace_dict) # return 45

28. Lazy loading.	
	Defer creating/initializing an object until it is needed.
29. List of optimizers in TF.
	tf.train.GradientDescentOptimizer
	tf.train.AdagradOptimizer
	tf.train.MomentumOptimizer
	tf.train.AdamOptimizer
	tf.train.ProximalGradientDescentOptimizer
	tf.train.ProximalAdagradOptimizer
	tf.train.RMSPropOptimizer
	and more
30. Overal all structure of a model in TensorFlow:
	Phase 1: Assemble graph.
		1. Define placeholders for input and output
		2. Define the weights
		3. Define the inference model
		4. Define the loss function.
		5. Define the optimizer
	Phase 2: Compute
		1. Initilize model parameters 
		2. Input trainning data
		3. Execute inference model on tranning data.
		4. Compute loss.
		5. Adjust mode parameters and go to 2.

