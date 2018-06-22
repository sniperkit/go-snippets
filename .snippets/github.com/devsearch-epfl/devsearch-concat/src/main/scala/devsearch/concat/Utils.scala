package devsearch.concat

import java.io._
import java.nio.file.{ Paths, Files, Path }

import org.apache.commons.compress.archivers.ArchiveStreamFactory
import org.apache.commons.compress.archivers.tar.{ TarArchiveInputStream, TarArchiveEntry }
import org.apache.tika.Tika
import org.apache.tika.detect.TextDetector
import scala.collection.JavaConverters._

import scala.util.{ Failure, Try, Success }

/**
 * Contains useful functions used by the workers as well as default values
 *
 * Created by dengels on 04/05/15.
 */
object Utils {

  import scala.Console._

  /** Default values **/

  /** Size of one blob we build */
  val BLOB_SIZE: Long = 640L << 20 // 640 Mb

  /** Text files bigger than 2Mb are probably not source code */
  val MAX_FILE_SIZE: Long = 2L << 20 // 2Mb

  /** Default number of worker nodes */
  val DEFAULT_PARALLELISM: Int = 4

  /**
   * Checks whether a stream contains text data
   *
   * @param is The input stream to read from
   * @return true if is contains text data
   */
  def isTextFile(is: InputStream): Boolean = try {
    val contentType = Option(new Tika(new TextDetector()).detect(is))
    contentType.exists(_.contains("text"))
  } catch {
    case o: IOException =>
      err.println("Input Stream could not be read to check if file is text")
      false
  }

  /**
   * Checks whether a file contains text data
   *
   * @return true if file is a text file
   */
  def isTextFile(file: Path): Boolean = try {
    isTextFile(new FileInputStream(file.toFile))
  } catch {
    case e: FileNotFoundException =>
      err.println(s"Could not find file $file")
      false
  }

  /** List all the files in a folder recursively */
  def listFilesRec(folder: Path): Stream[Path] = {

    def recScan(toScan: Stream[Path]): Stream[Path] = toScan match {
      case f #:: fs =>
        val children = f.toFile.listFiles.toSeq.map(_.toPath)
        val goodFiles = children.filterNot {
          p =>
            try {
              Files.isHidden(p) || Files.isSymbolicLink(p)
            } catch {
              case e: IOException => true
            }
        }
        val (folders, files) = goodFiles.partition(Files.isDirectory(_))
        files.toStream #::: recScan(folders.toStream #::: fs)
      case _ => Stream.empty[Path]
    }

    recScan(Stream(folder))
  }

  /**
   * Returns all the repos contained in the repository root
   *
   * repos are organized by language, then owner, then reponame
   *
   * repos can be either folder, or tarballs
   */
  def getRepoPaths(repoRoot: Path): Stream[Path] = {
    val languages = repoRoot.toFile.listFiles.toStream

    languages.filterNot(_.isDirectory).foreach[Unit] { file =>
      err.println(s"Found regular file $file when expecting language directory")
    }

    val owners = languages.filter(_.isDirectory).flatMap(_.listFiles)

    owners.filterNot(_.isDirectory).foreach { file =>
      err.println(s"Found regular file $file when expecting owner directory")
    }

    owners.filter(_.isDirectory).flatMap(_.listFiles.toStream.map(_.toPath))
  }

  /**
   * Generic way of representing a file whether it comes from a tar ball or the regular file system
   */
  trait FileEntry {
    def relativePath: String

    @throws[IOException]
    def inputStream: InputStream

    def size: Long
  }

  /** Whether the repo is a directory or has been put in a tar ball */
  def walkFiles[T](repo: Path)(processEntry: FileEntry => T): TraversableOnce[T] = {
    if (Files.isDirectory(repo)) {
      listFilesRec(repo).map { p =>

        val tryProcess = Try(Files.size(p)).map { fileSize =>

          lazy val is = new BufferedInputStream(Files.newInputStream(p))
          val entry = new FileEntry {

            override def size: Long = fileSize

            override def relativePath: String = repo.relativize(p).toString

            override def inputStream: InputStream = is
          }

          val res = processEntry(entry)
          try {
            is.close()
          } catch {
            case e: IOException => err.println(s"Could not close stream for entry $p")
          }
          res
        }

        tryProcess match {
          case Failure(e) => err.println(s"Could not process $p : ${e.getMessage}")
          case _ =>
        }

        tryProcess
      }.collect { case Success(value) => value }
    } else if (repo.toString.endsWith(".tar")) {
      val is = new BufferedInputStream(Files.newInputStream(repo))

      @SuppressWarnings(Array("org.brianmckenna.wartremover.warts.AsInstanceOf"))
      val tarInput = new ArchiveStreamFactory().createArchiveInputStream(ArchiveStreamFactory.TAR, is).asInstanceOf[TarArchiveInputStream]

      @SuppressWarnings(Array("org.brianmckenna.wartremover.warts.Var"))
      var entry: TarArchiveEntry = tarInput.getNextTarEntry

      @SuppressWarnings(Array("org.brianmckenna.wartremover.warts.Var"))
      var results = Seq.empty[T]

      while (Option(entry).isDefined) {
        /* Since the tarInput does not support marking, we have to cheat */
        lazy val is = new BufferedInputStream(tarInput)
        val fileEntry = new FileEntry {
          override def size: Long = entry.getSize

          override def relativePath: String = entry.getName.dropWhile(_ != '/').tail

          override def inputStream: InputStream = is
        }
        /** Filter symbolic files and directories and hidden files */
        val isHidden = Paths.get(fileEntry.relativePath).iterator.asScala.exists(_.toString.startsWith("."))
        if (!entry.isSymbolicLink && !entry.isDirectory && !isHidden) {
          val res = processEntry(fileEntry)
          results = results :+ res
        }

        entry = tarInput.getNextTarEntry
      }
      tarInput.close()
      results
    } else {
      err.println(s"$repo does not appear to be a folder or tarball repository")
      Iterator.empty
    }

  }
}
