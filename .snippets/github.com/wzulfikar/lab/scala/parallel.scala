/**
 * Using parallel function and thread in Scala.
 * 
 * Run in scastie: https://scastie.scala-lang.org
 * 
 */

// Higher-order function that displays 
// time elapsed to execute any given function
def timer (taskId : String, f: () => Unit) = {
	// use `val` to declare immutable object
	val t1 = System.currentTimeMillis()

	// run the function
	f()
  
	// use `var` to declare mutable object
	var msg = "→ Time elapsed: " + ((System.currentTimeMillis()-t1) / 1e3) + " µs"
	if (taskId != Nil) {
		msg += " (" + taskId + ")"
	}
	println(msg)
}

// Higher-order function to run any function as thread
def threader (f: () => Unit) = {
	(new Thread() {
		override def run() {
			f()
		}
	}).start()
}

// Function that we'll use as filter
def isPrime(x: Int): Boolean = {
	for (i <- 2 until x) {
		if ((x % i) == 0) {
			return true
		}
	}
	return false
}

// declare an immutable collection
val numbers = 1 to 1e5.toInt

threader(() => {
	println("Task 1: Running filter on 1e5 numbers in sequence..")
	timer("Task 1", () => numbers.filter(isPrime))
})

threader(() => {
	println("Task 2: Running filter on 1e5 numbers in parallel..")
	timer("Task 2", () => numbers.par.filter(isPrime))
})

println("End of code")
