import org.apache.spark.SparkContext
import org.apache.spark.SparkContext._
import org.apache.spark.SparkConf
import org.apache.spark.graphx._


object Pagerank {
  def main(args: Array[String]) {
    // val sc = new SparkContext(args(0), "PageRank", System.getenv("SPARK_HOME"),
    //   SparkContext.jarOfClass(this.getClass).toSeq)

    val conf = new SparkConf()
                         .set("spark.driver.host", args(1))
                         .setJars(SparkContext.jarOfClass(this.getClass).toSeq)
                         .setSparkHome(System.getenv("SPARK_HOME"))
    val sc = new SparkContext(args(0), "PageRank", conf)

    // Load the edges as a graph
    val graph = GraphLoader.edgeListFile(sc, args(2))
    // Run PageRank
    val ranks = graph.pageRank(0.0001).vertices
    // Join the ranks with the usernames
    val users = sc.textFile(args(3)).map { line =>
      val fields = line.split(",")
                                          (fields(0).toLong, fields(1))
                                        }
    val ranksByUsername = users.join(ranks).map {
      case (id, (username, rank)) => (username, rank)
    }
    // Print the result
    ranksByUsername.saveAsTextFile(args(4))
  }
}
