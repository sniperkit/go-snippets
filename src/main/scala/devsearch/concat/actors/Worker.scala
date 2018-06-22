package devsearch.concat.actors

import java.io._
import java.nio.file.{ Paths, Files, Path }

import akka.actor.{ Actor, ActorLogging, ActorRef, Props }
import devsearch.concat.Utils
import devsearch.concat.actors.Coordinator._
import devsearch.concat.actors.Worker._
import org.apache.commons.compress.archivers.{ ArchiveOutputStream, ArchiveStreamFactory }
import org.apache.commons.compress.archivers.tar.{ TarArchiveEntry, TarArchiveOutputStream }
import org.apache.commons.compress.compressors.{ CompressorOutputStream, CompressorStreamFactory }
import org.apache.commons.compress.utils.IOUtils
import org.apache.commons.io.FilenameUtils

/**
 * The worker actor is the one in charge of creating the large
 * files from smaller ones.
 *
 * @param master The actor to which it has to request new files
 */
class Worker(master: ActorRef) extends Actor with ActorLogging {

  case class Stats(totalBytesSeen: Long, totalBytesProcessed: Long) {
    def add(seen: Long, processed: Long) = Stats(totalBytesSeen + seen, totalBytesProcessed + processed)
  }

  object Stats {
    def empty: Stats = Stats(0, 0)
  }

  /**
   * Default worker behaviour at startup
   */
  override def receive: PartialFunction[Any, Unit] = {
    /* Start to work */
    case Begin =>
      master ! BlobRequest
      context.become(awaitBlob(Stats.empty))

  }

  def awaitBlob(stats: Stats): PartialFunction[Any, Unit] = {
    /* New blob to be created, it will contain the concatenation of many files */
    case BlobResponse(file) => {
      val out = new BufferedOutputStream(Files.newOutputStream(file))

      @SuppressWarnings(Array("org.brianmckenna.wartremover.warts.AsInstanceOf"))
      val tarOut: TarArchiveOutputStream = new ArchiveStreamFactory().createArchiveOutputStream(ArchiveStreamFactory.TAR, out).asInstanceOf[TarArchiveOutputStream]
      tarOut.setLongFileMode(TarArchiveOutputStream.LONGFILE_GNU)
      master ! RepoRequest
      context.become(fillBlob(file, tarOut, 0, stats))
    }
  }

  def fillBlob(blob: Path, stream: TarArchiveOutputStream, bytesWritten: Long, stats: Stats): PartialFunction[Any, Unit] = {
    /* One single file to append to the current blob */
    case RepoResponse(file, relativePath) =>
      val correctedPath = FilenameUtils.removeExtension(relativePath)

      val sizes = Utils.walkFiles(file) { fileEntry =>
        try {
          val size = fileEntry.size

          val isReasonableSize = size < Utils.MAX_FILE_SIZE
          /** This lazy is important because is TextFile might read the whole file in memory */
          lazy val isTextFile = Utils.isTextFile(fileEntry.inputStream)

          if (isReasonableSize && isTextFile) {
            val entry = new TarArchiveEntry(Paths.get(correctedPath, fileEntry.relativePath).toString)
            entry.setSize(size)
            stream.putArchiveEntry(entry)
            val processed = IOUtils.copy(fileEntry.inputStream, stream)
            stream.closeArchiveEntry()
            (processed, size)
          } else {
            (0L, size)
          }
        } catch {
          case e: IOException =>
            log.error(e, s"Encountered error when processing ${fileEntry.relativePath} for repo ${relativePath}")
            (0L, 0L)
        }
      }

      val (processed, seen) = sizes.foldLeft((0L, 0L)) { (b, a) => (b._1 + a._1, b._2 + a._2) }

      val updatedStats = stats.add(seen, processed)

      val totalWritten = bytesWritten + processed

      if (totalWritten >= Utils.BLOB_SIZE) {
        stream.close()
        log.info(s"Finished with blob : $blob")
        master ! BlobRequest
        context.become(awaitBlob(updatedStats))
      } else {
        master ! RepoRequest
        context.become(fillBlob(blob, stream, totalWritten, updatedStats))
      }

    /* No more files, end what you are doing and send finished message */
    case Shutdown =>
      stream.close()
      log.info(s"Closing last blob : $blob")
      sender ! Finished(stats.totalBytesSeen, stats.totalBytesProcessed)
  }
}

object Worker {
  def props(reader: ActorRef): Props = Props(new Worker(reader))

  case object RepoRequest

  case object BlobRequest

  case object Begin

  case class Finished(bytesSeen: Long, bytesProcessed: Long)

}
