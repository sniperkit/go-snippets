package devsearch.concat.actors

import java.nio.file.{ Files, InvalidPathException, Path }

import akka.actor.{ Actor, ActorLogging, Props }
import devsearch.concat.Utils
import devsearch.concat.actors.Coordinator._
import devsearch.concat.actors.Worker._
import org.apache.commons.io.FileUtils._

/**
 * The coordinator is in charge of listing the files in the input directory and
 * distribute them to the Workers.
 *
 * @param repoRoot The input folder where to find source files
 * @param outputFolder  The output folder where to output the big blobs
 * @param numWorkers The number of workers that it will have to coordinate
 */
case class Coordinator(repoRoot: Path, outputFolder: Path, numWorkers: Int) extends Actor with ActorLogging {

  import Utils._

  log.info(s"Starting up coordinator with $numWorkers worker on root folder : $repoRoot and ouput folder $outputFolder")

  /**
   * Default behaviour of coordiantor on startup
   */
  override def receive: PartialFunction[Any, Unit] = serverWorkers(0, getRepoPaths(repoRoot))

  /**
   * Serve the workers with repositories or blob numbers
   * @param currentBlobNum the number of the last blob that has been given out to a worker
   * @param repos all the repos that haven't been processed yet
   */
  def serverWorkers(currentBlobNum: Int, repos: Stream[Path]): PartialFunction[Any, Unit] = {
    /* Send next repo to worker */
    case RepoRequest => repos match {
      case repo #:: tail =>
        sender ! RepoResponse(repo, getRelativePath(repo, repoRoot))
        context.become(serverWorkers(currentBlobNum, tail))

      /* If there are no more files, shutdown workers */
      case _ =>
        sender ! Shutdown
        context.become(collectWorkers(numWorkers, 0, 0))
    }

    /* Send next available blob to worker */
    case BlobRequest =>
      val nextBlobNum = currentBlobNum + 1
      val blobName = "part-%05d.tar".format(nextBlobNum)
      sender ! BlobResponse(Files.createFile(outputFolder.resolve(blobName)))
      context.become(serverWorkers(nextBlobNum, repos))
  }

  /**
   * The coordinator will wait for all the workers to join before exiting
   */
  def collectWorkers(numLeft: Int, bytesSeen: Long, bytesProcessed: Long): PartialFunction[Any, Unit] = {
    /* Worker is done doing its work */
    case Finished(seen, processed) =>
      val totalSeen = bytesSeen + seen
      val totalProcessed = bytesProcessed + processed

      if (numLeft == 1) {
        log.info(s"Shutting down system! Processed ${byteCountToDisplaySize(totalProcessed)} out of ${byteCountToDisplaySize(totalSeen)}")
        context.system.shutdown()
      } else {
        context.become(collectWorkers(numLeft - 1, totalSeen, totalProcessed))
      }

    case _ => sender ! Shutdown
  }

}

object Coordinator {
  def props(langFolder: Path, outputFolder: Path, numWorkers: Int): Props =
    Props(new Coordinator(langFolder, outputFolder, numWorkers))

  case class RepoResponse(file: Path, relativeParentPath: String)

  case class BlobResponse(file: Path)

  case object Shutdown

  def getRelativePath(file: Path, repoRoot: Path): String = {
    repoRoot.relativize(file).toString
  }
}
