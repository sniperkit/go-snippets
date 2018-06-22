package devmine.srccat

import java.nio.file.{ Path, Files, Paths }
import java.util.concurrent.Executors

import akka.actor.ActorSystem
import devmine.srccat.actors.Worker.Begin
import devmine.srccat.actors.{ Coordinator, Worker }
import scopt.OptionParser

import scala.concurrent.ExecutionContext

/**
 * Main object of concat project
 */
object Main {

  @SuppressWarnings(Array("org.brianmckenna.wartremover.warts.DefaultArguments"))
  case class Config(repoRoot: String = "", outputFolder: String = "", parallelism: Int = Utils.DEFAULT_PARALLELISM)

  def concat(repoRoot: Path, outputFolder: Path, parallelism: Int): Unit = {
    val numWorkers = parallelism
    val threadPool = Executors.newFixedThreadPool(numWorkers)
    val executionContext = ExecutionContext.fromExecutor(threadPool)

    /* Create new actor system */
    val system = ActorSystem("srccat", defaultExecutionContext = Some(executionContext))

    /* Initiate devmine.srccat.actors */
    val master = system.actorOf(Coordinator.props(repoRoot, outputFolder, numWorkers))
    val workers = Vector.fill(numWorkers)(system.actorOf(Worker.props(master)))

    /* Start working */
    workers.foreach {
      _ ! Begin
    }

    /* Wait for termination */
    system.awaitTermination()

    threadPool.shutdown()
  }

  def main(args: Array[String]) {
    import Files._

    @SuppressWarnings(Array("org.brianmckenna.wartremover.warts.NonUnitStatements"))
    val parser: OptionParser[Config] = new OptionParser[Config]("srccat") {
      opt[Int]('j', "jobs").text("Maximum number of jobs to run").action((j, c) => c.copy(parallelism = j))
      arg[String]("<REPO_ROOT>").text("Repository root").action((repo, c) => c.copy(repoRoot = repo))
      arg[String]("<OUTPUT_FOLDER>").text("Output folder for big files").action((out, c) => c.copy(outputFolder = out))
    }

    def fail(msg: String): Nothing = {
      Console.err.println(s"[srccat] ERROR : $msg")
      sys.exit(1)
    }

    val conf = parser.parse(args, Config()) getOrElse sys.exit(1)

    val repoRoot = Paths.get(conf.repoRoot)
    val outputFolder = Paths.get(conf.outputFolder)

    if (!isDirectory(repoRoot)) fail("Repository root is not a directory")
    if (!isDirectory(outputFolder)) fail("Output folder is not a directory")

    if (!outputFolder.toFile.list.isEmpty) fail("Output folder is not empty")

    repoRoot.toFile.listFiles.filterNot(_.isDirectory).foreach[Unit] { file =>
      fail(s"Found $file in the repository root which is not a directory!")
    }

    concat(repoRoot, outputFolder, conf.parallelism)
  }
}

