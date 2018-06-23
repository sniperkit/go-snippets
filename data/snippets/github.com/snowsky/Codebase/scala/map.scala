object testMap {
  def p[A](args: List[A]): Unit = args.foreach(println)
  def pl[A](args: List[A]): Unit = println(args)
  def f(x: Int) = if (x > 2) Some(x) else None
  def g(v: Int) = List(v-1, v, v+1)
  def main(args: Array[String]) {
    val l = List(1,2,3,4,5)
    println(l)
    println(l.map(x => x*2))
    println(l.map(x => f(x)))
    println(l.map(x => g(x)))
    println(l.flatMap(x => g(x)))
    println(l.map(x => f(x)))
    println(l.flatMap(x => f(x)))

    val m = Map(1->2, 2->4, 3->6)
    m.toList
    println(m)
    println(m.mapValues(v => v*2))
    println(m.mapValues(v => f(v)))
    p(m.toList)
    println(m.flatMap(e => List(e._2)))
  }
}
testMap.main(args)
